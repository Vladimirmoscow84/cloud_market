package storage

import (
	"cloud_market/internal/cache"
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

	return nil
}

// FillingCache заполняет кэш из БД

func (s *Storage) FillingCache(ctx context.Context, cache *cache.Cache) error {
	//временная структура для результатов JOIN
	type orderRow struct {
		model.Order
		Delivery model.Delivery `db:"delivery"`
		Payment  model.Payment  `db:"payment"`
	}
	var rows []orderRow

	// Основной запрос с JOIN для delivery и payment
	err := s.DB.SelectContext(ctx, &rows, `
	SELECT 
		o.*,
		d.name as "delivery.name",
		d.phone as "delivery.phone",
		d.zip as "delivery.zip",
		d.city as "delivery.city",
		d.address as "delivery.address",
		d.region as "delivery.region",
		d.email as "delivery.email",
		p.transaction as "payment.transaction",
		p.request_id as "payment.request_id",
		p.currency as "payment.currency",
		p.provider as "payment.provider",
		p.amount as "payment.amount",
		p.payment_dt as "payment.payment_dt",
		p.bank as "payment.bank",
		p.delivery_cost as "payment.delivery_cost",
		p.goods_total as "payment.goods_total",
		p.custom_fee as "payment.custom_fee"
		FROM orders o
		LEFT JOIN delivery d ON d.order_id = o.id
		LEFT JOIN payment p ON p.order_id = o.id
		ORDER BY o.id DESC 
		LIMIT 2`)

	if err != nil {
		return fmt.Errorf("failed to get orders with joins: %w", err)
	}

	for _, row := range rows {
		order := row.Order
		order.Delivery = row.Delivery
		order.Payment = row.Payment

		err := s.DB.SelectContext(ctx, &order.Items, `
			SELECT chrt_id, track_number, price, rid, name, sale,
				   size, total_price, nm_id, brand, status
			FROM items WHERE order_id = $1`, order.ID)
		if err != nil {
			return fmt.Errorf("failed to get items for order %s: %w", order.OrderUID, err)
		}
		cache.Put(order)
	}

	return nil

}

// GetOrderById возвращает из БД данные в соответствии с запрошенным id
func (s *Storage) GetOrderById(ctx context.Context, orderUID string) (model.Order, error) {
	// времменная структура для заполнения ответа
	type orderRow struct {
		model.Order
		Delivery model.Delivery `db:"delivery"`
		Payment  model.Payment  `db:"payment"`
	}
	var row orderRow

	// Основной запрос с JOIN для delivery и payment
	err := s.DB.SelectContext(ctx, &row, `
	SELECT 
		o.*,
		d.name as "delivery.name",
		d.phone as "delivery.phone",
		d.zip as "delivery.zip",
		d.city as "delivery.city",
		d.address as "delivery.address",
		d.region as "delivery.region",
		d.email as "delivery.email",
		p.transaction as "payment.transaction",
		p.request_id as "payment.request_id",
		p.currency as "payment.currency",
		p.provider as "payment.provider",
		p.amount as "payment.amount",
		p.payment_dt as "payment.payment_dt",
		p.bank as "payment.bank",
		p.delivery_cost as "payment.delivery_cost",
		p.goods_total as "payment.goods_total",
		p.custom_fee as "payment.custom_fee"
		FROM orders o
		LEFT JOIN delivery d ON d.order_id = o.id
		LEFT JOIN payment p ON p.order_id = o.id
		WHERE o.order_uid = $1
		`, orderUID)

	if err != nil {
		return model.Order{}, fmt.Errorf("failed to get orders with joins: %w", err)
	}

	order := row.Order
	order.Delivery = row.Delivery
	order.Payment = row.Payment

	err = s.DB.SelectContext(ctx, &order.Items, `
			SELECT chrt_id, track_number, price, rid, name, sale,
				   size, total_price, nm_id, brand, status
			FROM items WHERE order_id = $1`, order.ID)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to get items for order %s: %w", order.OrderUID, err)
	}
	return order, nil

}
