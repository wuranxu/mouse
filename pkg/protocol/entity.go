package protocol

import (
	json "github.com/json-iterator/go"
	"log"
)

type HTTPMethod string

const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
)

type Response struct {
	success bool
	Data    string `json:"data"`
	raw     []byte
}

func NewResponse(success bool, raw []byte) *Response {
	return &Response{
		success: success,
		Data:    string(raw),
		raw:     raw,
	}
}

func (r *Response) JSON(out interface{}) error {
	return json.Unmarshal(r.raw, out)
}

func (r *Response) Text() string {
	return r.Data
}

type RequestOption func(h *HTTPRequest)

type HTTPRequest struct {
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	// allow redirect to another url
	AllowRedirect bool `json:"allowRedirect"`
	// request body
	Body   any        `json:"body"`
	Method HTTPMethod `json:"method"`
}

type HTTPResponse struct {
	*Response
	Error      error
	StatusCode int `json:"statusCode"`
	//ResponseHeaders map[string]string `json:"responseHeaders"`
	Headers map[string]string `json:"headers"`
	Elapsed int64             `json:"elapsed"`
	Request *HTTPRequest      `json:"request"`
}

type GRPCResponse struct {
	Response
}

func WithHeader(key, value string) RequestOption {
	return func(h *HTTPRequest) {
		if h.Headers == nil {
			h.Headers = map[string]string{key: value}
			return
		}
		h.Headers[key] = value
	}
}

func WithHeaders(kv map[string]string) RequestOption {
	return func(h *HTTPRequest) {
		h.Headers = kv
	}
}

func WithHeaderString(data string) RequestOption {
	return func(h *HTTPRequest) {
		var header map[string]string
		if err := json.Unmarshal([]byte(data), &header); err != nil {
			log.Println("http header unmarshal failed, ", err)
			return
		}
		h.Headers = header
	}
}

func WithBody(data any) RequestOption {
	return func(h *HTTPRequest) {
		h.Body = data
	}
}

func WithRedirect(redirect bool) RequestOption {
	return func(h *HTTPRequest) {
		h.AllowRedirect = redirect
	}
}
