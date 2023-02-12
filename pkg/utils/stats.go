package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"io"
	"log"
	"math"
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

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// MD5 returns the md5 hash of strings.
func MD5(slice ...string) string {
	h := md5.New()
	for _, v := range slice {
		io.WriteString(h, v)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
