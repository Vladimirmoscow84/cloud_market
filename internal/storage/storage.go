package storage

import (
	"cloud_market/internal/model"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// структура для работы с БД (хранилищем)
type Storage struct {
	DB *sqlx.DB
}

// 1.Создается функция-конструктор для создания экземпляра структруры Storage
func New(databaseURI string) (*Storage, error) {
	db, err := sqlx.Connect("pgx", databaseURI)
	if err != nil {
		return nil, err
	}

	return &Storage{
		DB: db,
	}, nil
}

func (s *Storage) AddOrder(order model.Order) error {
	s.DB.Exec(`
		BEGIN;
	
		INSERT INTO orders
			(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_chard)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	
		INSERT INTO delivery
			(name, phone, zip, city, address, region, email)
		VALUES
			($12, $13, $14, $15, $16, $17, $18);
		
		INSERT INTO payment
			(transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES
			($19, $20, $21, $22, $23, $24, $25, $26, $27, $28);

		INSERT INTO items
			(chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
		VALUES
			($29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39)

		COMMIT;

	`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee, append(order.Items, order.Items...))
	return nil
}
