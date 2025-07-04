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
	router.Get("/order", r.IdHandler_Get)
	router.Options("/order", r.Opt)

	return router
}

func (rt *Router) Opt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
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
			switch err.Error() {
			case "no exists in database":
				sendError(w, "no exists in database", err)
				return
			default:
				sendError(w, "error capturing order from database", err)
				return
			}
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

	resp, err2 := json.Marshal(resptaskErr)
	if err2 != nil {
		http.Error(w, fmt.Sprintf("%s: %s", errText, err2.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write(resp)
}
