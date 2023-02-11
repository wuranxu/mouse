package core

import (
	"context"
	"log"
	"testing"
)

func TestNewRateLimiter(t *testing.T) {
	data := NewRateLimiter(10, 10)
	ctx := context.Background()
	for i := 0; i < 10000; i++ {
		err := data.Wait(ctx)
		if err == nil {
			log.Println("ok")
		}
	}
}
