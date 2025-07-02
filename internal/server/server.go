package server

import (
	"cloud_market/internal/storage"
	"context"
	"fmt"
	"net/http"
)

func Run() {

	databaseURI := "host=localhost port=7701 user=postgres password=password dbname=cloud_market sslmode=disable"
	// 2. Создается экземпляр структуры storage.Storage для дальнейшей работы с БД (хранилищем)
	strg, err := storage.New(databaseURI)
	if err != nil {
		fmt.Println("postgres DB initialization error:", err)
		return
	}
	defer strg.DB.Close()

	// 7. Создание экземпляра структуры консьюмер
	// consumer := newConsumer(strg)

	//8 запускаем в отдельной горутине  консьюмер
	// go consumer.readMessage()

	kc := kafkaConfig{
		Brokers:     ":9092",
		Topic:       "cloud_market",
		Group:       "cloud_market_group",
		MessageChan: make(chan []byte),
	}

	ks := NewKafkaService(context.Background(), kc)
	defer ks.Stop()
	go ks.Start(context.Background())
	krs := kafkaReadService{
		messagesChan: kc.MessageChan,
		strg:         strg,
	}
	go krs.Process(context.Background())

	// 4. Создается экземпляр структуры Router для
	router := NewRouter(strg)

	// 6. Запуск локального сервера
	err = http.ListenAndServe(":7540", router.Routers())
	if err != nil {
		fmt.Println("error of start server:", err)
		return
	}

}
