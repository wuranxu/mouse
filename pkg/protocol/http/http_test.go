package http

import (
	"github.com/wuranxu/mouse/pkg/protocol"
	"log"
	"testing"
)

func TestHTTPClient_Get(t *testing.T) {
	client := NewHTTPClient()
	request := NewRequest("https://www.baidu.com", protocol.GET)
	resp := client.Get(request)
	if resp.Error != nil {
		t.Error("request failed", resp.Error)
		return
	}
	log.Println("request success: ", resp.Response.Data)
}

func TestClient_Get(t *testing.T) {
	resp := Post("https://api.pity.fun/auth/login", map[string]string{"Content-Type": "application/json"}, `
{"username": "wpop", "password": "23"}
`)
	if resp.Error != nil {
		t.Error("request failed", resp.Error)
		return
	}
	result := struct {
		Code int    `json:"code"`
		Data any    `json:"data,omitempty"`
		Msg  string `json:"msg"`
	}{}
	err := resp.JSON(&result)
	if err != nil {
		t.Error("parse result error: ", err)
		return
	}
	log.Println("request success: ", result)
}
