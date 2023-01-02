// package cache implements in-memory cache to store TimeSeries structure in sorted list
package cache

import (
	"strconv"
	"sync"
	"time"

	"github.com/ggrafu/sticker/utils"
)

const LOCAL_CACHE_TTL time.Duration = time.Minute

type Cache struct {
	Updated time.Time
	Head    *CacheNode
	mu      sync.RWMutex
}

type CacheNode struct {
	next  *CacheNode
	prev  *CacheNode
	date  time.Time
	value float32
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Clear() {
	c.Head = nil
}

func (c *Cache) Add(date time.Time, val float32) {
	n := c.Head
	if n == nil {
		c.Head = &CacheNode{
			date:  date,
			value: val,
		}
		return
	}
	// insert to the left
	if n.date.Before(date) {
		c.Head = &CacheNode{
			date:  date,
			value: val,
			next:  n,
		}
		n.prev = c.Head
		return
	}

	// traverse until found the older date
	for ; n.next != nil && n.next.date.After(date); n = n.next {
	}

	n.next = &CacheNode{
		prev:  n.prev,
		next:  n.next,
		date:  date,
		value: val,
	}
	n.next.prev = n
	return
}

// function GetLastElements returns most recent N elements from the cache
func (c *Cache) GetLastElements(n int) []float32 {
	if n < 0 {
		return []float32{}
	}
	result := make([]float32, n)
	node := c.Head
	i := 0
	for ; i < n && node != nil; i++ {
		result[i] = node.value
		node = node.next
	}
	return result[:i]
}

// function IsOutdated returns true if the cache was updated more then CACHE_INVALIDATION_TIME ago
// or if cache was not used before
func (c *Cache) IsOutdated() bool {
	return c.Updated.IsZero() || c.Updated.Add(LOCAL_CACHE_TTL).Before(time.Now())
}

// function Update reassembles cache from TimeSeries data structure
func (c *Cache) Update(ts *utils.APIData) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	c.Updated = time.Now()
	c.Head = nil

	for k, v := range ts.TimeSeries {
		f, err := strconv.ParseFloat(v.Close, 32)
		if err != nil {
			return err
		}
		date, err := time.Parse("2006-01-02", k)
		if err != nil {
			return err
		}
		c.Add(date, float32(f))
	}
	return nil
}
