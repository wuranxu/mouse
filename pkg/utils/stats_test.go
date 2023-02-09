package utils

import "testing"

func TestGetCpuAndMemoryUsage(t *testing.T) {
	mem, cpu := GetCpuAndMemoryUsage()
	t.Log("memory: ", mem, "cpu: ", cpu)
	if mem == 0.0 && cpu == 0.0 {
		t.Error("get current pid cpu/mem usage failed.")
	}
}
