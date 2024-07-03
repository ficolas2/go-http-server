package main

import (
	"github.com/ficolas2/http-server"
)

type CustomRequest struct {
	Body   CustomRequestBody
	Param  string `param:"param"`
	Header string `header:"header"`
}

type CustomRequestBody struct {
	Message string `json:"message"`
}

type CustomResponseBody struct {
	Message string `json:"message"`
}

func validate(body CustomRequestBody) *http.Result {
	if body.Message == "" {
		err := http.NewBadRequest("Message is required")
		return err
	}
	return nil
}

func main() {
	var server http.HttpServer
	server.Port = 8080
	http.AddRequestMapping(&server, http.POST, "/", func(request *CustomRequest) *http.Result {
		err := validate(request.Body)
		if err != nil {
			return err
		}
		return http.NewOk(CustomResponseBody{"Hello world"})
	})
	server.StartHttpServer()
}
