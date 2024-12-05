package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/codingpierogi/streaming-demo/consumers/notify/config"
	"github.com/codingpierogi/streaming-demo/consumers/notify/util"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Starting Notify Consumer...")

	config, err := config.LoadConfig("consumers/notify/")

	if err != nil {
		log.Fatalf("Fatal error loading config: %s\n", err)
	}

	log.Printf("Creating Consumer [bootstrap.servers=%s] [group.id=%s] [auto.offset.reset=%s]\n", config.BootstrapServers, config.GroupId, config.AutoOffsetReset)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"group.id":          config.GroupId,
		"auto.offset.reset": config.AutoOffsetReset,
	})

	if err != nil {
		log.Fatalf("Fatal error creating consumer: %s\n", err)
	}

	defer consumer.Close()

	log.Printf("Subscribing to topic [topic=%s]", config.Topic)

	err = consumer.SubscribeTopics([]string{config.Topic}, nil)

	if err != nil {
		log.Fatalf("Fatal error subscribing to topic: %s\n", err)
	}

	run := true
	for run {
		select {
		case <-sigs:
			log.Println("Received signal, exiting...")
			run = false
		default:
			log.Printf("Polling for messages [timeoutMs=%d]", config.TimeoutMs)

			ev := consumer.Poll(config.TimeoutMs)

			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				consumeNotifyMessage(e)
			case kafka.Error:
				log.Printf("Kafka error: %v\n", e)
			}
		}
	}

}

func consumeNotifyMessage(e *kafka.Message) {
	nm, err := util.DeserializeMessage(e.Value)

	if err != nil {
		log.Fatalf("Fatal message deserialization error: %v\n", err)
	}

	log.Printf("Consumed message [topic=%s] [message=%+v]\n", *e.TopicPartition.Topic, nm)
}
