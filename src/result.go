package http

import (
	"fmt"
)

type Result struct {
	statusCode int
	body       string
	headers    map[string]string
}

func (self *Result) AddHeader(key, value string) {
	self.headers[key] = value
}

func (self *Result) GetHeader(key string) string {
	return self.headers[key]
}

func (self *Result) RemoveHeader(key string) {
	delete(self.headers, key)
}

func (self *Result) String() string {
	result := fmt.Sprintf("HTTP/1.1 %d\r\n", self.statusCode)
	for key, value := range self.headers {
		result += key + ": " + value + "\r\n"
	}
	result += "\r\n" + self.body
	return result
}

func NewResult(statusCode int, body string) Result {
	return Result{statusCode, body, make(map[string]string)}
}

func NewOk(body string) Result {
	return NewResult(200, body)
}

func NewBadRequest(body string) Result {
	return NewResult(400, body)
}

func NewNotFound(body string) Result {
	return NewResult(404, body)
}

func NewInternalServerError(body string) Result {
	return NewResult(500, body)
}
