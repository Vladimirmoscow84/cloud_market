package server

import (
	"cloud_market/internal/storage"
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	strg   *storage.Storage
	reader *kafka.Reader
}

// создаем конструктор консьюмера
func newConsumer(strg *storage.Storage) *Consumer {

	return &Consumer{
		strg: strg,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "order",
		}),
	}

}

// создаем отдельный фметод для запуска консьюмера
func (c *Consumer) readMessage() {
	for {
		message, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Printf("получено сообщение: %s\n", string(message.Value))

		// Коммитим оффсет вручную после обработки
		err = c.reader.CommitMessages(context.Background(), message)
		if err != nil {
			panic(err)
		}
	}

}
