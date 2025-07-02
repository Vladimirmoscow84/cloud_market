package server

import (
	"cloud_market/internal/cache"
	"cloud_market/internal/storage"
	"context"
	"fmt"
	"net/http"
	"time"
)

func Run() {
	ctx := context.Background()

	databaseURI := "host=localhost port=7701 user=postgres password=password dbname=cloud_market sslmode=disable"
	// 2. Создается экземпляр структуры storage.Storage для дальнейшей работы с БД (хранилищем)
	strg, err := storage.New(databaseURI)
	if err != nil {
		fmt.Println("postgres DB initialization error:", err)
		return
	}
	defer strg.DB.Close()

	//...Создаем экземпляр кэш
	c := cache.NewCashe()
	//...Заполняем кэш из БД
	strg.FillingCache(ctx, c)

	//Проверка работы fillinCache
	fmt.Println("Внимание!!! КЭШ!!!")
	fmt.Println("Спасибо за внимание")
	c.Out()
	time.Sleep(3 * time.Second)

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

	ks := NewKafkaService(ctx, kc)
	defer ks.Stop()
	go ks.Start(ctx)
	krs := kafkaReadService{
		messagesChan: kc.MessageChan,
		strg:         strg,
	}
	go krs.Process(ctx)

	// 4. Создается экземпляр структуры Router для
	router := NewRouter(strg)

	// 6. Запуск локального сервера
	err = http.ListenAndServe(":7540", router.Routers())
	if err != nil {
		fmt.Println("error of start server:", err)
		return
	}

}
