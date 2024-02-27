package repository

import (
	"context"
	"github.com/yngk19/wb_l0task/internal/repository/models"
	"sync"
)

type Cache struct {
	mtx sync.RWMutex
	Capacity int
	Data     map[int]interface{}
}

type OrderGetter interface {
	GetOrdersByLimit(context.Context, int) ([]models.Order, error)
}

func (c *Cache) Put(id int, value interface{}) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if len(c.Data) == c.Capacity {
		for key, _ := range c.Data {
			if key != id {
				delete(c.Data, key)
				break
			}
		}
	}
	_, ok := c.Data[id]
	if !ok {
		c.Data[id] = value
	}
	return ok
}

func (c *Cache) Get(id int) interface{} {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	value, ok := c.Data[id]
	if !ok {
		return nil
	}
	delete(c.Data, id)
	return value
}

func NewCache(cap int) *Cache {
	return &Cache{
		mtx: sync.RWMutex{},
		Capacity: cap,
		Data:     make(map[int]interface{}),
	}
}

func (c *Cache) GetFromDB(ctx context.Context, orderGetter OrderGetter) {
	orders, err := orderGetter.GetOrdersByLimit(ctx, c.Capacity-len(c.Data))
	if err != nil {
		return
	}
	for _, order := range orders {
		c.Put(order.ID, order.Data)
	}
}
