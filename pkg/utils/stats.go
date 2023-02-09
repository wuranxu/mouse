package utils

import (
	"github.com/shirou/gopsutil/process"
	"log"
	"os"
	"runtime"
)

// GetCpuAndMemoryUsage get current cpu/memory usage percent
func GetCpuAndMemoryUsage() (float64, float64) {
	pid := os.Getpid()
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		log.Println("get cpu usage failed: ", err)
		return 0.0, 0.0
	}
	percent, err := proc.CPUPercent()
	if err != nil {
		log.Println("get cpu percent failed: ", err)
		return 0.0, 0.0
	}
	info, err := proc.MemoryInfo()
	if err != nil {
		log.Println("get memory usage failed: ", err)
		return 0.0, 0.0
	}
	// covert bytes to mb
	return float64(info.RSS >> 20), percent / float64(runtime.NumCPU())
}
