package utils

import "testing"

func TestGetExternalIP(t *testing.T) {
	ip, err := GetExternalIP()
	if err != nil {
		t.Error(err)
	}
	t.Log(ip)
}
