package storage

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

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
