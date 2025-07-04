package server

import (
	"cloud_market/internal/cache"
	"cloud_market/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

// Структура для работы
type Router struct {
	strg  *storage.Storage
	cache *cache.Cache
}

// структура для вывода ответа в случае ошибки
type respTask struct {
	Error string `json:"error,omitempty"`
}

// 3. NewRouter  - функция-конструктор для создания экземпляра структуры Router
func NewRouter(strg *storage.Storage, cache *cache.Cache) *Router {

	return &Router{
		strg:  strg,
		cache: cache,
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

	if r.Method != "GET" {
		sendError(w, "error method", errors.New("error method"))
		return
	}

	uid := r.FormValue("order_uid")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var response []byte
	if rt.cache.IsExist(uid) {
		answer, err := rt.cache.Get(uid)
		if err != nil {
			sendError(w, "error capturing order from cache", err)
			return
		}
		response, err = json.MarshalIndent(answer, "", "\t")
		if err != nil {
			sendError(w, "error encoding json", err)
			return
		}
	} else {
		answer, err := rt.strg.GetOrderById(context.Background(), uid)
		if err != nil {
			sendError(w, "error capturing order from database", err)
			return
		}
		response, err = json.MarshalIndent(answer, "", "\t")
		if err != nil {
			sendError(w, "error encoding json", err)
			return
		}

	}
	w.Write(response)
}

// sendError - сериализация и отправка ошибки в формате JSON
func sendError(w http.ResponseWriter, errText string, err error) {
	var resptaskErr respTask
	resptaskErr.Error = errText
	fmt.Printf("%s: %s\n", errText, err.Error())

	resp, err := json.Marshal(resptaskErr)
	if err != nil {
		fmt.Printf("error encoding error json: %s\n", err.Error())
		http.Error(w, fmt.Sprintf("%s: %s", errText, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write(resp)
}
