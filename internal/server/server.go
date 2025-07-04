package server

import (
	"cloud_market/internal/cache"
	"cloud_market/internal/storage"
	"context"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

func Run() {
	ctx := context.Background()

	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Ошибка загрузки .env файла: %v\n", err)
		return
	}

	databaseURI := viper.GetString("DATABASE_URI")
	addr := viper.GetString("SERVER_ADDRESS")
	kafkaBroker := viper.GetString("KAFKA_BROKER")
	kafkaTopic := viper.GetString("KAFKA_TOPIC")
	kafkaGroup := viper.GetString("KAFKA_GROUP")

	// 2. Создается экземпляр структуры storage.Storage для дальнейшей работы с БД (хранилищем)
	strg, err := storage.New(databaseURI)
	if err != nil {
		fmt.Println("postgres DB initialization error:", err)
		return
	}
	defer strg.DB.Close()

	//...Создаем экземпляр кэш
	c := cache.NewCache()
	//...Заполняем кэш из БД
	err = strg.FillingCache(ctx, c)
	if err != nil {
		fmt.Printf("ошибка заполнения кэш: %v", err)
	}

	kc := kafkaConfig{
		Brokers:     kafkaBroker,
		Topic:       kafkaTopic,
		Group:       kafkaGroup,
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
	router := NewRouter(strg, c)

	// 6. Запуск локального сервера
	fmt.Printf("Server starting on %s\n", addr)
	err = http.ListenAndServe(addr, router.Routers())
	if err != nil {
		fmt.Println("error of start server:", err)
		return
	}

}
