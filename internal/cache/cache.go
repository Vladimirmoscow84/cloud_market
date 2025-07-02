package cache

import (
	"cloud_market/internal/model"
	"errors"
	"sync"
)

type Cashe struct {
	data map[string]model.Order
	mu   sync.RWMutex
}

// NewCache возвращает структуру с инициализированным кэшом
func NewCashe() *Cashe {
	return &Cashe{
		data: make(map[string]model.Order),
		mu:   sync.RWMutex{},
	}
}

// Get возвращает данные из кэша в соответствии с запрошенным id  или ошибку, если указанного id нет в кэш
func (c *Cashe) Get(id string) (model.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	order, ok := c.data[id]
	if !ok {
		return model.Order{}, errors.New("the cache does not contain the order")
	}
	return order, nil
}

// Put добавляет заказ в кэш
func (c *Cashe) Put(order model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[order.OrderUID] = order

}

// IsExist проверяет наличие заказа в кэш по его id
func (c *Cashe) IsExist(id string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.data[id]; !ok {
		return false
	}
	return true
}
