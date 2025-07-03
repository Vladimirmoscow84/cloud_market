package server

import (
	"cloud_market/internal/cache"
	"cloud_market/internal/storage"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	err = strg.FillingCache(ctx, c)
	if err != nil {
		fmt.Printf("ошибка заполнения кэш: %v", err)
	}

	//Проверка работы fillinCache
	// fmt.Println("Внимание!!! КЭШ!!!")
	// fmt.Println("Спасибо за внимание")
	// c.Out()
	// time.Sleep(3 * time.Second)

	// 7. Создание экземпляра структуры консьюмер
	// consumer := newConsumer(strg)

	//8 запускаем в отдельной горутине  консьюмер
	// go consumer.readMessage()
	uid := "b563feb7b2b84b6test_2"
	fmt.Printf("получение данных по order_uid: %s", uid)
	answer, err := strg.GetOrderById(ctx, uid)
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := json.MarshalIndent(answer, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(data))

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
