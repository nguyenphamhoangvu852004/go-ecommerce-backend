package initialize

import (
	"go-ecommerce-backend-api/global"
	"log"

	"github.com/segmentio/kafka-go"
)

// Init Kafka Producer

var KafkaProducer *kafka.Writer

func InitKafka() {
	global.KafkaProducer = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "opt-auth-topic",
		Balancer: &kafka.LeastBytes{},
	}
}

func CloseKafka() {
	if err := global.KafkaProducer.Close(); err != nil {
		log.Fatalf("failed to close kafka producer: %s", err)
	}
}
