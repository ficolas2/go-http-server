package http

import (
	"github.com/ficolas2/http-server"
)

func getRoot(request string) string {
	return "Hello World"
}

func main() {
	var server http.HttpServer
	server.AddRequestMapping(http.GET, "/", getRoot)
	server.StartHttpServer()
}
