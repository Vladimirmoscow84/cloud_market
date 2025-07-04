package main

import (
	"cloud_market/internal/model"
	"cloud_market/internal/storage"
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()

	// databaseURI := "host=localhost port=7701 user=postgres password=password dbname=cloud_market sslmode=disable"
	databaseURI := "host=localhost port=5432 user=postgres password=password dbname=cloud_market sslmode=disable"
	// 2. Создается экземпляр структуры storage.Storage для дальнейшей работы с БД (хранилищем)
	strg, err := storage.New(databaseURI)
	if err != nil {
		fmt.Printf("postgres DB initialization error: %v\n", err)
		return
	}
	defer strg.DB.Close()

	t, err := time.Parse(time.RFC3339Nano, "2021-11-26T06:22:19Z")
	if err != nil {
		fmt.Printf("error time parse: %v\n", err)
	}
	for i := range 5 {
		order := model.Order{
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

		err := strg.AddOrder(ctx, order)
		if err != nil {
			fmt.Printf("error DB: %v\n", err)
			return
		}
	}

}
