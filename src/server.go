package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"
)

type HttpServer struct {
	mappings map[Method]map[string]reflect.Value
	Port     int
}

func AddRequestMapping[F any](server *HttpServer, method Method, path string, f func(*F) Result) {
	if server.mappings == nil {
		server.mappings = make(map[Method]map[string]reflect.Value)
	}
	if _, methodExists := server.mappings[method]; !methodExists {
		server.mappings[method] = make(map[string]reflect.Value)
	}
	if _, exists := server.mappings[method][path]; exists {
		log.Fatalf("Mapping already exists for %s %s", method, path)
		return
	}
	server.mappings[method][path] = reflect.ValueOf(f)
}

func (server *HttpServer) StartHttpServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", server.Port))
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		go server.handleConnection(connection)
	}
}

func (server *HttpServer) handleConnection(connection net.Conn) {
	defer connection.Close()

	buf := make([]byte, 8192)
	n, err := connection.Read(buf)
	if err != nil {
		log.Fatal(err)
		return
	}
	str := string(buf[:n])

	// Read request line
	start := 0
	length := strings.Index(str, "\r\n")

	firstLine := strings.Split(str[start:start+length], " ")

	method := MethodFromString(firstLine[0])
	path := firstLine[1]
	protocol := firstLine[2]

	// Read headers
	start = start + length + 2
	headers := make(map[string]string)
	for {
		length := strings.Index(str[start:], "\r\n")
		line := str[start : start+length]
		split := strings.Split(line, ": ")

		if len(split) == 1 {
			break
		}

		name := split[0]
		value := split[1]

		headers[name] = value
		start = start + length + 2
	}

	// Read body
	start = start + 2
	body := str[start:]

	// Process path params
	paramStart := strings.Index(path, "?")
	paramList := make(map[string]string)
	if paramStart != -1 {
		stringParamList := strings.Split(path[paramStart+1:], "&")
		for _, param := range stringParamList {
			split := strings.Split(param, "=")
			paramList[split[0]] = split[1]
		}
		path = path[:paramStart]
	}

	var result Result

	if protocol != "HTTP/1.1" {
		fmt.Println("Unsupported protocol")
		result = NewBadRequest("Unsupported protocol")
	}

	if f, exists := server.mappings[method][path]; exists {
		result = callMapping(f, headers, body, paramList)
	} else {
		result = NewNotFound("No mapping found " + path)
	}

	connection.Write([]byte(result.String()))
}

func callMapping(fnValue reflect.Value, headers map[string]string, bodyStr string, paramList map[string]string) Result {
	argType := fnValue.Type().In(0).Elem()

	if argType.Kind() == reflect.Struct {
		instance := reflect.New(argType)

		for i := 0; i < argType.NumField(); i++ {
			fieldName := argType.Field(i).Name
			fieldTag := argType.Field(i).Tag
			fieldType := argType.Field(i).Type

			required := fieldTag.Get("Required") == "true"

			if fieldName == "Body" {
				body := reflect.New(fieldType)
				err := json.Unmarshal([]byte(bodyStr), body.Interface())
				if err != nil {
					return NewBadRequest("Error parsing body " + err.Error())
				}

				instance.Elem().Field(i).Set(reflect.ValueOf(body.Elem().Interface()))
				continue
			}

			header := fieldTag.Get("header")
			if header != "" {
				if value, exists := headers[header]; exists {
					instance.Elem().Field(i).SetString(value)
				} else if required {
					return NewBadRequest("Header " + header + " is required")
				}
				continue
			}

			param := fieldTag.Get("param")
			if param != "" {
				if value, exists := paramList[param]; exists {
					instance.Elem().Field(i).SetString(value)
				} else if required {
					return NewBadRequest("Param " + param + " is required")
				}
				continue
			}

		}

		return fnValue.Call([]reflect.Value{instance})[0].Interface().(Result)
	}

	return NewInternalServerError("Mapping first argument is not a struct.")
}
