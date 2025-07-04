package server

import (
	"cloud_market/internal/cache"
	"cloud_market/internal/storage"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

// Структура для работы
type Router struct {
	strg *storage.Storage
}

// 3. NewRouter  - функция-конструктор для создания экземпляра структуры Router
func NewRouter(strg *storage.Storage) *Router {

	return &Router{
		strg: strg,
	}
}

// 5. Routers - метод для создания экземпляра роутера chi.Mux, который соответствует интрефейсу http.Handler

func (r *Router) Routers() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/order/order_uid", r.IdHandler_Get)

	return router
}

//IdHandler_Get  - функция для возврата данных заказа по номеру UID из кэш.
// в случвае отсутствия данных в кэш, данные берутся из БД

func (rt *Router) IdHandler_Get(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var cache cache.Cache

		uid := r.FormValue("order_uid")

		_ = rt.strg.FillingCache(context.Background(), &cache)

		if cache.IsExist(uid) {
			answer, err := cache.Get(uid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response, err := json.MarshalIndent(answer, "", "\t")
			w.Write([]byte(response))
		} else {
			answer, err := rt.strg.GetOrderById(context.Background(), uid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response, err := json.MarshalIndent(answer, "", "\t")
			w.Write([]byte(response))
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}
