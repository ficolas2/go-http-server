A simple http server library, build from scratch, using a TCP connection.

Check out **/examples** for example ussage. [Here](examples/example.go)

This project is not supossed to be efficient, secure, or complete. It's just a simple project to 
learn go and networking.

- [Features](#Features)
- [Library overview](#Library-overview)
    - [Start and configure the server](#Start-and-configure-the-server)
    - [Add request mappings](#Add-request-mappings)
        - [Request fields](#Request-fields)
        - [Response](#Response)
- [Testing](#Testing)

## Features
- HTTP/1.1 server
- Routing with simple path matching
- Headers and path parameters parsing
- Demarshalling of request JSON body
- Marshalling of response JSON body

## Library overview
### Start and configure the server
To create a server, you need to create a new instance of the HttpServer struct, and call
StartHttpServer() method.

```go
func main() {
    var server http.HttpServer
    server.Port = 8080
    server.StartHttpServer()
}
```

### Add request mappings
To add request mappings to the server, you can use the http.AddRequestMapping() method, which 
receives the http server struct, a method, and a function.

```go
http.AddRequestMapping(&server, http.GET, "/", func (request *CustomRequest) http.Result { 
    // code
})
```

#### Request fields
The first struct argument is a custom struct that contains the request data.

The field named Body will be the body of the request, unmarsheled to a struct, if the request is 
body a json.

You can get header info and path params from the request using tags.

```go
type CustomRequestBody struct {
    Message string `json:"message"`
}

type CustomRequest struct {
    Body   CustomRequestBody
    Header string `header:"header-name"`
    Param  string `param:"param-name"`
}
```

#### Response
The response returned by the function can correspond to an error, or an OK response.

The OK response can contain a deserialized JSON, if created with a struct, or a string.

Example:
```go
type CustomResponseBody struct {
    Field string `json:"field"`
}

func getRoot(r *CustomRequest) http.Result {
    err := validate(r.Body)
    if err != nil {
        return http.NewBadRequest("Validation failed")
    }

    return http.NewOK(CustomResponseBody{"Hello world"})
}
```

## Testing
To run the tests, run:

```bash
cd src
go test
```

## TODO
- TLS
- TCP server from scratch
- Error handling and better concurrency
- Websockets
- Content negotiation (Accepts)
- Static file serving (images, css, html)

