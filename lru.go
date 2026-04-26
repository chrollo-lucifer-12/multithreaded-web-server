package main

import (
	"sync"
	"time"
)

type CacheStore interface {
	find(url string) *Cache
	add_cache_element(data string, len int, url string) int
	remove_cache_element(url string)
}

type CacheElement struct {
	data     string
	len      int
	url      string
	lru_time time.Time
	next     *Cache
}

type Cache struct {
	head       *CacheElement
	cache_size int
	mu         sync.Mutex
}

func (c *Cache) find(url string) *Cache {
	return nil
}

func (c *Cache) add_cache_element(data string, len int, url string) int {
	return 0
}

func (c *Cache) remove_cache_element(url string) { return }
