package storage

import (
	"cloud_market/internal/model"
	"context"
	"fmt"

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

<<<<<<< HEAD
func (s *Storage) AddOrder(ctx context.Context, order model.Order) error {
	// начало транзакции
	fmt.Println("Начало транзакции")
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	// откат транзакции в случае ошибки
	defer func() {
		if err != nil {
			fmt.Println("Произошёл Rollback")
			tx.Rollback()
		}
	}()
=======
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
>>>>>>> ef9e0d1f54fc0fa8a9e3b1a59864572f69c3bcc6

	fmt.Println("Начало order")
	// запись данных в таблицу orders и получение из неё order_id для последующей записи order_id  в таблицы delivery, payment, items
	row := tx.QueryRowContext(ctx, `
	INSERT INTO orders
		(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_chard)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id;
	`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)

	var orderID int
	if err := row.Scan(&orderID); err != nil {
		return fmt.Errorf("ошибка добавления order: %w", err)
	}
	fmt.Println("Конец order")

	fmt.Println("Начало delivery")
	// запись данных в таблицу delivery
	_, err = tx.ExecContext(ctx, `
	INSERT INTO delivery
		(name, phone, zip, city, address, region, email, order_id)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8);
	`, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email, orderID)
	if err != nil {
		return fmt.Errorf("ошибка добавления delivery: %w", err)
	}
	fmt.Println("Конец delivery")

	fmt.Println("Начало payment")
	// запись данных в таблицу payment
	_, err = tx.ExecContext(ctx, `
	INSERT INTO payment
		(transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, order_id)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee, orderID)
	if err != nil {
		return fmt.Errorf("ошибка добавления payment: %w", err)
	}
	fmt.Println("Конец delivery")

	fmt.Println("Начало items")
	// запись данных из слайса в таблицу items
	if len(order.Items) > 0 {
		for i := range order.Items {
			order.Items[i].OrderID = orderID
		}
		_, err = tx.NamedExecContext(ctx, `
		INSERT INTO items
			(chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_id)
		VALUES
			(:chrt_id, :track_number, :price, :rid, :name, :sale, :size, :total_price, :nm_id, :brand, :status, :order_id);
		`, order.Items)
		if err != nil {
			return fmt.Errorf("ошибка добавления item: %w", err)
		}
	}
	fmt.Println("Конец items")

	// выполнение commit транзакции
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка выполнения commit транзакции: %w", err)
	}

<<<<<<< HEAD
=======
	`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee, append(order.Items, order.Items...))
>>>>>>> ef9e0d1f54fc0fa8a9e3b1a59864572f69c3bcc6
	return nil
}
