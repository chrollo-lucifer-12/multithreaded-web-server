package main

import (
	"sync"
	"time"
)

type CacheStore interface {
	find(url string) *CacheElement
	add_cache_element(data string, len int, url string) int
	remove_cache_element(url string)
}

type CacheElement struct {
	data     string
	len      int
	url      string
	lru_time time.Time
	next     *CacheElement
	prev     *CacheElement
}

type Cache struct {
	head       *CacheElement
	tail       *CacheElement
	cache_size int
	mu         sync.Mutex
}

func (c *Cache) move_to_front(elem *CacheElement) {
	if elem == c.head {
		return
	}
	prev_element := elem.prev
	next_element := elem.next

	prev_element.next = next_element
	next_element.prev = prev_element

	c.head.prev = elem
	elem.next = c.head

	c.head = elem
}

func (c *Cache) find(url string) *CacheElement {

	c.mu.Lock()
	defer c.mu.Unlock()

	temp := c.head

	for temp != nil {
		if temp.url == url {
			c.move_to_front(temp)
			return temp
		}
		temp = temp.next
	}
	return nil
}

func (c *Cache) add_cache_element(data string, len int, url string) int {
	newElement := &CacheElement{
		data:     data,
		len:      len,
		url:      url,
		lru_time: time.Now(),
	}

	newElement.prev = nil
	newElement.next = c.head

	c.head = newElement

	return len
}

func (c *Cache) remove_cache_element(url string) { return }
