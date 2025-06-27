package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBM/sarama"
)

type kafkaService struct {
	topic    string
	handler  sarama.ConsumerGroupHandler
	consumer sarama.ConsumerGroup
}

type consumerHandler struct {
	topic       string
	messageChan chan []byte
}

type kafkaConfig struct {
	Brokers     string
	Topic       string
	Group       string
	MessageChan chan []byte
}

func NewKafkaService(ctx context.Context, cfg kafkaConfig) *kafkaService {
	brokers := strings.Split(cfg.Brokers, ",")
	consumerCfg := sarama.NewConfig()
	consumerCfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	consumerCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	client, err := sarama.NewConsumerGroup(brokers, cfg.Group, consumerCfg)
	if err != nil {
		fmt.Printf("failed to create consumer: %v\n", err)
	}
	handler := &consumerHandler{
		topic:       cfg.Topic,
		messageChan: cfg.MessageChan,
	}
	return &kafkaService{
		topic:    cfg.Topic,
		handler:  handler,
		consumer: client,
	}
}

func (k *kafkaService) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := k.consumer.Consume(ctx, []string{k.topic}, k.handler)
			if err != nil {
				fmt.Printf("consume message error: %v\n", err)
			}
		}
	}
}

func (k *kafkaService) Stop() {
	k.consumer.Close()
}

func (h consumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.messageChan <- message.Value
		session.MarkMessage(message, "")
	}
	return nil
}

// type Consumer struct {
// 	strg   *storage.Storage
// 	reader *kafka.Reader
// }

// https://habr.com/ru/articles/894056/

// // создаем конструктор консьюмера
// func newConsumer(strg *storage.Storage) *Consumer {

// 	return &Consumer{
// 		strg: strg,
// 		reader: kafka.NewReader(kafka.ReaderConfig{
// 			Brokers: []string{"localhost:9092"},
// 			Topic:   "order",
// 		}),
// 	}
// }

// // создаем отдельный фметод для запуска консьюмера
// func (c *Consumer) readMessage() {
// 	for {
// 		message, err := c.reader.ReadMessage(context.Background())
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Printf("получено сообщение: %s\n", string(message.Value))

// 		// Коммитим оффсет вручную после обработки
// 		err = c.reader.CommitMessages(context.Background(), message)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// }
