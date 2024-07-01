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
	listener, err := net.Listen("tcp", ":8080")
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
	
	firstLine := strings.Split(str[start : start+length], " ")

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

	var result Result

	if protocol != "HTTP/1.1" {
		fmt.Println("Unsupported protocol")
		result = NewBadRequest("Unsupported protocol")
	}

	if f, exists := server.mappings[method][path]; exists {
		result = callMapping(f, headers, body)
	} else {
		fmt.Println("No mapping found")
		result = NewNotFound("No mapping found")
	}

	connection.Write([]byte(result.String()))
}

func callMapping(fnValue reflect.Value, headers map[string]string, bodyStr string) Result {
	argType := fnValue.Type().In(0).Elem()

	if argType.Kind() == reflect.Struct{
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
					fmt.Printf("Error parsing body: %s\n", err)
					return NewBadRequest("Error parsing body")
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

		}
		
		return fnValue.Call([]reflect.Value{instance})[0].Interface().(Result)
	}

	return NewInternalServerError("Mapping first argument is not a struct.")
}
