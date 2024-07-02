package http

import (
	stdHttp "net/http"
	"strings"
	"testing"
	"time"
)

type EmptyRequest struct {
}

type Body struct {
	Name string `json:"name"`
}

type TestBodyRequest struct {
	Body Body
}

type TestHeaderRequest struct {
	Header string `header:"Test-Header"`
}

const GET_ROOT = "get root"
const POST_ROOT = "post root"
const GET_CHILD = "get child"

func TestRouting(t *testing.T) {
	var server HttpServer
	server.Port = 8080
	AddRequestMapping(&server, GET, "/", func(request *EmptyRequest) Result {
		return NewOk(GET_ROOT)
	})
	AddRequestMapping(&server, POST, "/", func(request *EmptyRequest) Result {
		return NewOk(POST_ROOT)
	})
	AddRequestMapping(&server, GET, "/child", func(request *EmptyRequest) Result {
		return NewOk(GET_CHILD)
	})
	go server.StartHttpServer()
	time.Sleep(1 * time.Millisecond)

	buffer := make([]byte, 1024)

	// Test GET
	response, err := stdHttp.Get("http://localhost:8080/")
	if err != nil {
		t.Error(err)
		return
	}
	if response.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", response.StatusCode)
		return
	}

	n, err := response.Body.Read(buffer)
	if err != nil {
		t.Error(err)
		return
	}

	if string(buffer[:n]) != GET_ROOT {
		t.Errorf("Expected '%s', got '%s'", GET_ROOT, string(buffer))
		return
	}
	response.Body.Close()

	// Test POST
	buffer = make([]byte, 1024)
	response, err = stdHttp.Post("http://localhost:8080/", "text/plain", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if response.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", response.StatusCode)
		return
	}

	n, err = response.Body.Read(buffer)
	if err != nil {
		t.Error(err)
		return
	}

	if string(buffer[:n]) != POST_ROOT {
		t.Errorf("Expected '%s', got '%s'", POST_ROOT, string(buffer))
		return
	}
	response.Body.Close()

	// Test GET child
	buffer = make([]byte, 1024)
	response, err = stdHttp.Get("http://localhost:8080/child")
	if err != nil {
		t.Error(err)
		return
	}
	if response.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", response.StatusCode)
		return
	}

	n, err = response.Body.Read(buffer)
	if err != nil {
		t.Error(err)
		return
	}

	response.Body.Close()
}

func TestNotFound(t *testing.T) {
	var server HttpServer
	server.Port = 8081
	go server.StartHttpServer()
	time.Sleep(1 * time.Millisecond)

	response, err := stdHttp.Get("http://localhost:8080/notfound")
	if err != nil {
		t.Error(err)
		return
	}
	if response.StatusCode != 404 {
		t.Errorf("Expected 404, got %d", response.StatusCode)
		return
	}
	response.Body.Close()
}

func TestBody(t *testing.T) {
	body := `{"name": "test"}`
	var server HttpServer
	server.Port = 8082
	AddRequestMapping(&server, POST, "/", func(request *TestBodyRequest) Result {
		return NewOk(request.Body.Name)
	})
	go server.StartHttpServer()
	time.Sleep(1 * time.Millisecond)

	client := &stdHttp.Client{}
	req, err := stdHttp.NewRequest("POST", "http://localhost:8082/", strings.NewReader(body))
	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.ContentLength = int64(len(body))

	response, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	if response.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", response.StatusCode)
		return
	}

	buffer := make([]byte, 1024)
	n, err := response.Body.Read(buffer)

	if err != nil {
		t.Error(err)
		return
	}

	if string(buffer[:n]) != "test" {
		t.Errorf("Expected 'test', got '%s'", string(buffer))
		return
	}
}

func TestHeaders(t *testing.T) {
	var server HttpServer
	server.Port = 8083
	AddRequestMapping(&server, GET, "/", func(request *TestHeaderRequest) Result {
		return NewOk(request.Header)
	})
	go server.StartHttpServer()
	time.Sleep(1 * time.Millisecond)

	client := &stdHttp.Client{}
	req, err := stdHttp.NewRequest("GET", "http://localhost:8083/", nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Add("Test-Header", "test")

	response, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	if response.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", response.StatusCode)
		return
	}

	buffer := make([]byte, 1024)
	n, err := response.Body.Read(buffer)

	if err != nil {
		t.Error(err)
		return
	}

	if string(buffer[:n]) != "test" {
		t.Errorf("Expected 'test', got '%s'", string(buffer))
		return
	}
}
