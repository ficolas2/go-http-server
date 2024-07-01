package http

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type HttpServer struct {
	mappings map[Method]map[string]func(string) string
}

func (server *HttpServer) AddRequestMapping(method Method, path string, f func(string) string) {
	if server.mappings == nil {
		server.mappings = make(map[Method]map[string]func(string) string)
	}
	if _, methodExists := server.mappings[method]; !methodExists {
		server.mappings[method] = make(map[string]func(string) string)
	}
	if _, exists := server.mappings[method][path]; exists {
		log.Fatalf("Mapping already exists for %s %s", method, path)
		return
	}
	server.mappings[method][path] = f
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

	if _, exists := server.mappings[method][path]; exists {
		body := server.mappings[method][path](body)
		fmt.Println("Response body: " + body)
		connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n" + body))
	} else {
		fmt.Println("No mapping found")
		result = NewNotFound("No mapping found")
	}

	connection.Write([]byte(result.String()))
}
