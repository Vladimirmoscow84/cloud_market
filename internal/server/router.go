package server

import (
	"cloud_market/internal/storage"

	"github.com/go-chi/chi"
)

// Структура для работы
type Router struct {
	strg *storage.Storage
}

// 3. созадется функция-конструктор для создания экземпляра структуры Router
func NewRouter(strg *storage.Storage) *Router {

	return &Router{
		strg: strg,
	}
}

// 5. метод для создания экземпляра роутера chi.Mux, который соответствует интрефейсу http.Handler
func (r *Router) Routers() *chi.Mux {
	router := chi.NewRouter()

	return router
}
