package http

type Method int

const (
	GET Method = iota
	POST
	PUT
	DELETE
	HEAD
	OPTIONS
	PATCH
	TRACE
	CONNECT
)

func MethodFromString(method string) Method {
	switch method {
	case "GET":
		return GET
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "DELETE":
		return DELETE
	case "HEAD":
		return HEAD
	case "OPTIONS":
		return OPTIONS
	case "PATCH":
		return PATCH
	case "TRACE":
		return TRACE
	case "CONNECT":
		return CONNECT
	default:
		return GET
	}
}

func (method Method) String() string {
	switch method {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case HEAD:
		return "HEAD"
	case OPTIONS:
		return "OPTIONS"
	case PATCH:
		return "PATCH"
	case TRACE:
		return "TRACE"
	case CONNECT:
		return "CONNECT"
	default:
		return "GET"
	}
}
