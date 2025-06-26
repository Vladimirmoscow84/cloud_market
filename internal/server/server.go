package server

import (
	"cloud_market/internal/storage"
	"fmt"
	"net/http"
)

func Run() {

	databaseURI := "host=localhost port=5432 user=postgres password=password dbname=cloud_market sslmode=disable"
	// 2. Создается экземпляр структуры storage.Storage для дальнейшей работы с базой данных
	strg, err := storage.New(databaseURI)
	if err != nil {
		fmt.Println("postgres DB initialization error:", err)
		return
	}
	// 4. Создается экземпляр структуры Router для
	router := NewRouter(strg)

	// 6. Запуск локального сервера
	err = http.ListenAndServe(":7540", router.Routers())
	if err != nil {
		fmt.Println("error of start server:", err)
		return
	}

	router.strg.DB.Close()

}
