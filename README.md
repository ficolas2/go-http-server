A simple http server library, build from scratch, using a TCP connection.

Check out **/examples** for example ussage.
Check out **/src**, for the library code.

This project is not supossed to be efficient, secure, or complete. It's just a simple project to 
learn go and networking.

## Features
- HTTP 1.1 server
- Routing with simple path matching
- Demarshalling of request JSON body

## Library overview
To create a server, you need to create a new instance of the HttpServer struct, and call
StartHttpServer() method.

To add request mappings to the server, you can use the http.AddRequestMapping() method, which 
receives the http server struct, a method, and a function.

```go
func main() {
	var server http.HttpServer
	http.AddRequestMapping(&server, http.GET, "/", getRoot)
	server.StartHttpServer()
}
```

The function signature must be:
```go
func functionName(r *CustomRequest) http.Response { 
    // code
}
```

The CustomStruct is a struct that contains the request data.

The field named Body will be the body of the request, unmarsheled to a struct, if the request is 
body a json.

You can get header info from the request using tags.

```go
type CustomBody struct {
    Field string `json:"field"`
}

type CustomRequest struct {
    Body   CustomBody
    Header string `header:"header-name"`
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

