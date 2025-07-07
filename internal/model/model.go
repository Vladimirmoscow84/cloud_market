package model

import "time"

type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required" db:"order_uid"`
	TrackNumber       string    `json:"track_number" validate:"required" db:"track_number"`
	Entry             string    `json:"entry" validate:"required" db:"entry"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required"`
	Locale            string    `json:"locale" validate:"required" Db:"locale"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required" db:"customer_id"`
	DeliveryService   string    `json:"delivery_service" validate:"required" db:"delivery_service"`
	Shardkey          string    `json:"shardkey" validate:"required" db:"shardkey"`
	SmID              int       `json:"sm_id" validate:"required" db:"sm_id"`
	DateCreated       time.Time `json:"date_created" validate:"required" db:"date_created"`
	OofShard          string    `json:"oof_shard" validate:"required" db:"oof_chard"`
	ID                int       `json:"-" db:"id"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required" db:"name"`
	Phone   string `json:"phone" validate:"required" db:"phone"`
	Zip     string `json:"zip" validate:"required" db:"zip"`
	City    string `json:"city" validate:"required" db:"city"`
	Address string `json:"address" validate:"required" db:"address"`
	Region  string `json:"region" validate:"required" db:"region"`
	Email   string `json:"email" validate:"required" db:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required" db:"transaction"`
	RequestID    string `json:"request_id" db:"request_id"`
	Currency     string `json:"currency" validate:"required" db:"currency"`
	Provider     string `json:"provider" validate:"required" db:"provider"`
	Amount       int    `json:"amount" validate:"required" db:"amount"`
	PaymentDt    int    `json:"payment_dt" validate:"required" db:"payment_dt"`
	Bank         string `json:"bank" validate:"required" db:"bank"`
	DeliveryCost int    `json:"delivery_cost" validate:"required" db:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total" validate:"required" db:"goods_total"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id" validate:"required" db:"chrt_id"`
	TrackNumber string `json:"track_number" validate:"required" db:"track_number"`
	Price       int    `json:"price" validate:"required" db:"price"`
	Rid         string `json:"rid" validate:"required" db:"rid"`
	Name        string `json:"name" validate:"required" db:"name"`
	Sale        int    `json:"sale" db:"sale"`
	Size        string `json:"size" validate:"required" db:"size"`
	TotalPrice  int    `json:"total_price" db:"total_price"`
	NmID        int    `json:"nm_id" validate:"required" db:"nm_id"`
	Brand       string `json:"brand" validate:"required" db:"brand"`
	Status      int    `json:"status" validate:"required" db:"status"`
	OrderID     int    `json:"-" db:"order_id"`
}
