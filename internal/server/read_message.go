package server

import (
	"cloud_market/internal/model"
	"cloud_market/internal/storage"
	"context"
	"encoding/json"
	"fmt"
)

type kafkaReadService struct {
	messagesChan <-chan []byte
	strg         *storage.Storage
}

func (krs *kafkaReadService) Process(ctx context.Context) {
	for msg := range krs.messagesChan {
		m := model.Order{}

		err := json.Unmarshal(msg, &m)
		if err != nil {
			fmt.Printf("err unmarshal: %v\n", err)
			continue
		}
		fmt.Println("message:", m)
		err = krs.strg.AddOrder(ctx, m)
		if err != nil {
			fmt.Printf("err addOrder: %v\n", err)
		}
	}
}
