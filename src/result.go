package http

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Result struct {
	statusCode int
	body       string
	headers    map[string]string
}

func (self *Result) SetBodyFromStruct(body interface{}) error {
	return nil
}

func (self *Result) SetBody(body string) *Result {
	self.body = body
	return self
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

func NewResult(statusCode int, body string) *Result {
	result := Result{statusCode, body, make(map[string]string)}
	return &result
}

func NewOk(body interface{}) *Result {
	result := NewResult(200, "")

	bodyType := reflect.TypeOf(body)

	switch bodyType.Kind() {
	case reflect.Struct:
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		result.body = string(bodyBytes)
	case reflect.String:
		result.body = body.(string)
	default:
		panic("Unsupported body type " + bodyType.String())
	}
	return result
}

func NewBadRequest(body string) *Result {
	return NewResult(400, body)
}

func NewNotFound(body string) *Result {
	return NewResult(404, body)
}

func NewInternalServerError(body string) *Result {
	return NewResult(500, body)
}
