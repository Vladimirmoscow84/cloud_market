package server

import (
	"context"
	"encoding/json"
	"fmt"
)

type kafkaReadService struct {
	messagesChan <-chan []byte
}

type message struct {
	Field string `json:"field"`
}

func (krs *kafkaReadService) Process(ctx context.Context) {
	for msg := range krs.messagesChan {
		m := message{}

		err := json.Unmarshal(msg, &m)
		if err != nil {
			fmt.Printf("err unmarshal: %v\n", err)
			continue
		}
		fmt.Println("message:", m)
	}
}
