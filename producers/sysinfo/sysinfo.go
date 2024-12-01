package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/codingpierogi/streaming-demo/producers/sysinfo/config"
	"github.com/codingpierogi/streaming-demo/producers/sysinfo/util"
	"github.com/codingpierogi/streaming-demo/types"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/elastic/go-sysinfo"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Starting SysInfo Producer...")
	log.Printf("Detected [OS=%s] [ARCH=%s]\n", runtime.GOOS, runtime.GOARCH)

	config, err := config.LoadConfig("producers/sysinfo/")

	if err != nil {
		log.Fatalf("Fatal error loading config: %s\n", err)
	}

	log.Printf("Creating Producer [bootstrap.servers=%s]\n", config.BootstrapServers)

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
	})

	if err != nil {
		log.Fatalf("Fatal error creating producer: %s\n", err)
	}

	defer producer.Close()

	ticker := time.NewTicker(time.Duration(config.IntervalMs) * time.Millisecond)
	defer ticker.Stop()

	log.Printf("Started Producer [topic=%s]\n", config.Topic)

	for {
		select {
		case <-ticker.C:
			err = produceSysInfoMessage(producer, config.Topic)
			if err != nil {
				log.Fatalf("Received producer message error %s\n", err)
			}
		case <-sigs:
			log.Println("Received signal, exiting...")
			return
		}
	}
}

func produceSysInfoMessage(p *kafka.Producer, topic string) error {
	host, err := sysinfo.Host()

	if err != nil {
		return fmt.Errorf("host Error occurred: %s", err)
	}

	memory, _ := host.Memory()

	hmi := types.HostMemoryInfo{
		Total:     memory.Total,
		Used:      memory.Used,
		Available: memory.Available,
		Free:      memory.Free,
	}

	sim := types.SysInfoMessage{
		Key:   "host_memory_info",
		Value: hmi,
	}

	sm, err := util.SerializeMessage(sim)

	if err != nil {
		return fmt.Errorf("failed to serialize message: %s", err)
	}

	km := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: sm,
	}

	dc := make(chan kafka.Event)
	err = p.Produce(km, dc)

	if err != nil {
		return fmt.Errorf("failed to produce message %w", err)
	}

	e := <-dc
	ekm := e.(*kafka.Message)

	if ekm.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %s", ekm.TopicPartition.Error)
	}

	close(dc)
	return nil
}
