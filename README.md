A simple http server library, build from scratch, using a TCP connection.

Check out examples for example ussage.
Check out src, for the library code.

This project is not supossed to be efficient, or complete. It's just a simple project to learn go.

In the future, I may implement the TCP server from scratch too, or add TLS, but for now, I'm using
the net package, and https is not supported.

## Library overview
To create a server, you need to create a new instance of the HttpServer struct, and call
ListenAndServe() method.

To add request mappings to the server, you can use the AddRequestMapping() method, which receives
a method, and a function.

```go
func main() {
    var server http.HttpServer
    server.AddRequestMapping(http.GET, "/", getRoot)
    server.StartHttpServer()
}
```

