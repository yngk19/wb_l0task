package repository

import (
	"sync"
	"time"
)

type Cache struct {
	Capacity int
	TTL      time.Duration
	Data     map[int]interface{}
}

func (c *Cache) Serve() {
	return
}

func (c *Cache) Put(id int, value interface{}) {
	var mtx sync.Mutex
	mtx.Lock()
	defer mtx.Unlock()
	if len(c.Data) == c.Capacity {
		for key, _ := range c.Data {
			if key != id {
				delete(c.Data, key)
				break
			}
		}
	}
	c.Data[id] = value
}

func (c *Cache) Get(id int) interface{} {
	var mtx sync.Mutex
	mtx.Lock()
	defer mtx.Unlock()
	value, ok := c.Data[id]
	if !ok {
		return nil
	}
	delete(c.Data, id)
	return value
}

func NewCache(cap int, ttl time.Duration) *Cache {
	return &Cache{
		Capacity: cap,
		TTL:      ttl,
	}
}
