package main

import (
	"cloud_market/internal/model"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

func main() {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// Создание нового продюсера

	producer, err := sarama.NewSyncProducer([]string{"172.31.141.173:9092"}, config)
	if err != nil {
		fmt.Printf("failed to start producer: %v\n", err)
	}
	defer producer.Close()

	t, err := time.Parse(time.RFC3339Nano, "2021-11-26T06:22:19Z")
	if err != nil {
		fmt.Printf("error time parse: %v\n", err)
	}
	for i := range 50 {
		m := model.Order{
			OrderUID:    fmt.Sprintf("b563feb7b2b84b6test_%d", i),
			TrackNumber: "sgg4",
			Entry:       "sdw",
			Delivery: model.Delivery{
				Name: "sdell",
			},
			Payment: model.Payment{
				Transaction: "sdgsdg",
			},
			Items: []model.Item{
				{ChrtID: 234234},
			},
			DateCreated: t,
		}

		mJSON, err := json.Marshal(m)
		if err != nil {
			fmt.Printf("error marshal: %v\n", err)
			return
		}

		//Отправка сообщение в топик kafka
		message := &sarama.ProducerMessage{
			Topic: "cloud_market",
			Value: sarama.ByteEncoder(mJSON),
		}

		partititon, offset, err := producer.SendMessage(message)
		if err != nil {
			fmt.Printf("failed to sendt message: %v\n", err)
		} else {
			fmt.Printf("Message sent to partition %d at offset %d\n", partititon, offset)
		}
		time.Sleep(time.Second)
	}

}
