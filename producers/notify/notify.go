package main

import (
	"fmt"
	"log"

	"github.com/codingpierogi/streaming-demo/producers/notify/config"
	"github.com/codingpierogi/streaming-demo/producers/notify/util"
	"github.com/codingpierogi/streaming-demo/types"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting Notify Producer...")

	config, err := config.LoadConfig("producers/notify/")

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

	log.Printf("Creating Router [mode=%s] [addr=%s]\n", config.Mode, config.Addr)

	gin.SetMode(config.Mode)
	router := gin.New()
	router.POST("/notify", func(c *gin.Context) {
		message := c.PostForm("message")
		log.Printf("POST /notify [message=%s]\n", message)
		err := produceNotifyMessage(producer, config.Topic, message)

		if err != nil {
			log.Fatalf("Received producer message error %s\n", err)
		}
	})

	log.Printf("🚀 listening and serving on 0.0.0.0%s\n", config.Addr)

	err = router.Run(config.Addr)

	if err != nil {
		log.Fatalf("Fatal router error: %s\n", err)
	}
}

func produceNotifyMessage(p *kafka.Producer, topic string, message string) error {
	n := types.Notification{
		Message: message,
	}

	nm := types.NotifyMessage{
		Key:   "notify",
		Value: n,
	}

	snm, err := util.SerializeMessage(nm)

	if err != nil {
		return fmt.Errorf("failed to serialize message: %s", err)
	}

	km := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: snm,
	}

	dc := make(chan kafka.Event)
	err = p.Produce(km, dc)

	if err != nil {
		return fmt.Errorf("failed to produce message %w", err)
	}

	log.Printf("Produced message [topic=%s] [message=%+v]\n", topic, message)

	e := <-dc
	ekm := e.(*kafka.Message)

	if ekm.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %s", ekm.TopicPartition.Error)
	}

	close(dc)
	return nil
}
