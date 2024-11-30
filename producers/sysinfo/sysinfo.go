package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/codingpierogi/streaming-demo/producers/sysinfo/config"
	"github.com/codingpierogi/streaming-demo/types"
	"github.com/elastic/go-sysinfo"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Starting SysInfo Producer...")
	log.Printf("Detected [OS=%v] [ARCH=%v]", runtime.GOOS, runtime.GOARCH)

	config, err := config.LoadConfig("producers/sysinfo/")

	if err != nil {
		log.Fatalf("Fatal error loading config: %v", err)
	}

	ticker := time.NewTicker(time.Duration(config.IntervalMs) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			hmi := produceSysInfo()
			log.Printf("%+v", hmi)
		case <-sigs:
			log.Println("Received signal, exiting...")
			return
		}
	}
}

func produceSysInfo() types.HostMemoryInfo {
	host, err := sysinfo.Host()

	if err != nil {
		log.Fatalf("Host Error occurred: %v", err)
	}

	memory, _ := host.Memory()

	return types.HostMemoryInfo{
		Total:     memory.Total,
		Used:      memory.Used,
		Available: memory.Available,
		Free:      memory.Free,
	}
}
