package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/wuranxu/mouse/pkg/protocol"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrorParsedUrl    = errors.New("parse url error")
	ErrorReadResponse = errors.New("read response error")
)

type Invoker interface {
	Get(*protocol.HTTPRequest) *protocol.HTTPResponse
	Post(*protocol.HTTPRequest) *protocol.HTTPResponse
	Put(*protocol.HTTPRequest) *protocol.HTTPResponse
	Delete(*protocol.HTTPRequest) *protocol.HTTPResponse
	Do(*protocol.HTTPRequest) *protocol.HTTPResponse
}

type Option func(*Client)

type Client struct {
	Timeout int64 `json:"timeout"`
	client  *http.Client
}

var httpClient = &Client{client: new(http.Client)}

func NewRequest(url string, method protocol.HTTPMethod, options ...protocol.RequestOption) *protocol.HTTPRequest {
	request := &protocol.HTTPRequest{Url: url, AllowRedirect: true, Method: method}
	for _, opt := range options {
		opt(request)
	}
	return request
}

func Get(url string, headers map[string]string, options ...protocol.RequestOption) *protocol.HTTPResponse {
	opts := append(options, protocol.WithHeaders(headers))
	request := NewRequest(url, protocol.GET, opts...)
	return httpClient.Get(request)
}

func Post(url string, headers map[string]string, data any, options ...protocol.RequestOption) *protocol.HTTPResponse {
	opts := append(options, protocol.WithHeaders(headers), protocol.WithBody(data))
	request := NewRequest(url, protocol.POST, opts...)
	return httpClient.Do(request)
}

func NewHTTPClient(options ...Option) *Client {
	client := &Client{client: new(http.Client)}
	for _, opt := range options {
		opt(client)
	}
	return client
}

func WithTimeout(ms int64) Option {
	return func(c *Client) {
		if ms < 500 {
			return
		}
		c.Timeout = ms
	}
}

func (h *Client) Get(request *protocol.HTTPRequest) *protocol.HTTPResponse {
	request.Method = protocol.GET
	return h.Do(request)
}

func (h *Client) Post(request *protocol.HTTPRequest) *protocol.HTTPResponse {
	request.Method = protocol.POST
	return h.Do(request)
}

func (h *Client) Put(request *protocol.HTTPRequest) *protocol.HTTPResponse {
	request.Method = protocol.PUT
	return h.Do(request)
}

func (h *Client) Delete(request *protocol.HTTPRequest) *protocol.HTTPResponse {
	request.Method = protocol.DELETE
	return h.Do(request)
}

func (h *Client) makeRequestBody(request *protocol.HTTPRequest, req *http.Request) {
	if request.Body != nil {
		var data []byte
		switch body := request.Body.(type) {
		case []byte:
			data = body
		case string:
			data = []byte(body)
		default:
			data, _ = json.Marshal(body)
		}
		req.Body = io.NopCloser(bytes.NewReader(data))
	}
}

func (h *Client) makeHeaders(resp *http.Response) map[string]string {
	headers := make(map[string]string, len(resp.Header))
	for k, v := range resp.Header {
		headers[k] = v[0]
	}
	return headers
}

func (h *Client) Do(request *protocol.HTTPRequest) *protocol.HTTPResponse {
	uri, err := url.Parse(request.Url)
	if err != nil {
		return &protocol.HTTPResponse{Error: ErrorParsedUrl, Response: protocol.NewResponse(false, []byte{})}
	}
	header := make(map[string][]string)
	if request.Headers != nil {
		for k, v := range request.Headers {
			header[k] = []string{v}
		}
	}
	req := &http.Request{
		Method: string(request.Method),
		URL:    uri,
		Header: header,
	}
	h.makeRequestBody(request, req)
	if h.Timeout > 0 {
		h.client.Timeout = time.Duration(h.Timeout) * time.Microsecond
	}
	start := time.Now()
	resp, err := h.client.Do(req)
	end := time.Since(start).Microseconds() / 1000
	if err != nil {
		return &protocol.HTTPResponse{Error: err, Elapsed: end, Response: protocol.NewResponse(false, nil)}
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return &protocol.HTTPResponse{StatusCode: resp.StatusCode, Error: ErrorReadResponse, Elapsed: end, Response: protocol.NewResponse(false, result)}
	}
	if err != nil {
		return &protocol.HTTPResponse{StatusCode: resp.StatusCode, Error: err, Elapsed: end, Response: protocol.NewResponse(false, result)}
	}
	return &protocol.HTTPResponse{
		Error:      nil,
		StatusCode: resp.StatusCode,
		Headers:    h.makeHeaders(resp),
		Elapsed:    end,
		Request:    request,
		Response:   protocol.NewResponse(true, result),
	}

}
