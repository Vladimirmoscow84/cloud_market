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
			(order_uid, track_number)
		VALUES
			($1, $2);
	
		INSERT INTO orders
			(order_uid, track_number)
		VALUES
			($3, $4);

		COMMIT;

	`, order.OrderUID, order.TrackNumber)
	return nil
}
