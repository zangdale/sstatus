package sstatus

import (
	"errors"
	"sync"
	"time"
)

const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

var (
	ErrKeyNotFound  = errors.New("key not found")
	ErrKeyIsExist   = errors.New("key already exists")
	ErrKeyIsExpired = errors.New("key already expired")
)

type Data struct {
	Value          any
	IsNoExpiration bool
	ExpirationTime time.Time
}

type cache struct {
	defaultExpiration time.Duration
	mu                sync.RWMutex
	data              map[string]*Data
}

func NewCache(t time.Duration) *cache {
	return &cache{
		mu:                sync.RWMutex{},
		data:              make(map[string]*Data),
		defaultExpiration: t,
	}
}

func (c *cache) Add(k string, v any, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_d, ok := c.data[k]
	if ok {
		if _d.IsNoExpiration {
			return ErrKeyIsExist
		} else {
			if time.Since(_d.ExpirationTime) > 0 {
				delete(c.data, k)
			} else {
				return ErrKeyIsExist
			}
		}
	}

	data := &Data{Value: v}

	if d == DefaultExpiration {
		data.ExpirationTime = time.Now().Add(c.defaultExpiration)
	} else if d == NoExpiration {
		data.IsNoExpiration = true
	} else {
		data.ExpirationTime = time.Now().Add(d)
	}
	c.data[k] = data

	return nil
}

func (c *cache) MustAdd(k string, v any, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	data := &Data{Value: v}

	if d == DefaultExpiration {
		data.ExpirationTime = time.Now().Add(c.defaultExpiration)
	} else if d == NoExpiration {
		data.IsNoExpiration = true
	} else {
		data.ExpirationTime = time.Now().Add(d)
	}
	c.data[k] = data

	return nil
}

func (c *cache) Get(k string) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_d, ok := c.data[k]
	if !ok {
		return nil, ErrKeyNotFound
	}
	if !_d.IsNoExpiration {
		if time.Since(_d.ExpirationTime) > 0 {
			c.mu.Lock()
			delete(c.data, k)
			c.mu.Unlock()
			return nil, ErrKeyNotFound
		}
	}
	return _d.Value, nil
}

func (c *cache) Delete(k string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, k)
	return nil
}

func (c *cache) GetWithExpiration(k string) (any, time.Time, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_d, ok := c.data[k]
	if ok {
		if !_d.IsNoExpiration &&
			time.Since(_d.ExpirationTime) > 0 {
			c.mu.Lock()
			delete(c.data, k)
			c.mu.Unlock()
			return nil, time.Time{}, ErrKeyNotFound
		}
		return _d.Value, _d.ExpirationTime, nil
	}

	return nil, time.Time{}, ErrKeyNotFound
}
