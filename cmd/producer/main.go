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
	for i := range 5 {
		m := model.Order{
			OrderUID:    fmt.Sprintf("b563feb7b2b84b6test_%d", i),
			TrackNumber: "WBILMTESTTRACK",
			Entry:       "WBIL",
			Delivery: model.Delivery{
				Name:    "Test Testov",
				Phone:   "+9720000000",
				Zip:     "2639809",
				City:    "Kiryat Mozkin",
				Address: "Ploshad Mira 15",
				Region:  "Kraiot",
				Email:   "test@gmail.com",
			},
			Payment: model.Payment{
				Transaction:  "b563feb7b2b84b6test",
				RequestID:    "",
				Currency:     "USD",
				Provider:     "wbpay",
				Amount:       1817,
				PaymentDt:    1637907727,
				Bank:         "alpha",
				DeliveryCost: 1500,
				GoodsTotal:   317,
				CustomFee:    0,
			},
			Items: []model.Item{
				{
					ChrtID:      9934930,
					TrackNumber: "WBILMTESTTRACK",
					Price:       453,
					Rid:         "ab4219087a764ae0btest",
					Name:        "Mascaras",
					Sale:        30,
					Size:        "0",
					TotalPrice:  317,
					NmID:        2389212,
					Brand:       "Vivienne Sabo",
					Status:      202,
				},
				{
					ChrtID:      9934922,
					TrackNumber: "WBILMTESTTRACK",
					Price:       453,
					Rid:         "ab4219087a764ae0btest",
					Name:        "Mascaras",
					Sale:        30,
					Size:        "0",
					TotalPrice:  317,
					NmID:        2389212,
					Brand:       "Vivienne Sabo",
					Status:      203,
				},
				{
					ChrtID:      9934910,
					TrackNumber: "WBILMTESTTRACK",
					Price:       453,
					Rid:         "ab4219087a764ae0btest",
					Name:        "Mascaras",
					Sale:        30,
					Size:        "0",
					TotalPrice:  317,
					NmID:        2389212,
					Brand:       "Vivienne Sabo",
					Status:      207,
				},
			},
			Locale:            "en",
			InternalSignature: "",
			CustomerID:        "test",
			DeliveryService:   "meest",
			Shardkey:          "9",
			SmID:              99,
			DateCreated:       t,
			OofShard:          "1",
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
		time.Sleep(time.Second * 15)
	}

}
