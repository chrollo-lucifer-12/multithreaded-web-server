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

func (c *Cache) find(url string) *Cache

func (c *Cache) add_cache_element(data string, len int, url string) int

func (c *Cache) remove_cache_element(url string)
