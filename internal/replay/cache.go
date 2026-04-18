// Package replay implements a nonce cache to prevent replay attacks.
package replay

import (
	"sync"
	"time"
)

// NonceCache tracks seen nonces with their arrival time.
type NonceCache struct {
    mu     sync.Mutex
    entries map[string]time.Time
}

// NewNonceCache initializes the map using make.
func NewNonceCache() *NonceCache {
	return &NonceCache{
		entries: make(map[string]time.Time),
	}
}

// Add checks if the nonce exists using the "comma ok" idiom.
func (c *NonceCache) Add(nonce string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.entries[nonce]; exists {
		return false
	}
	c.entries[nonce] = time.Now()
	return true
}

// Purge iterates over the map and removes entries older than the TTL.
func (c *NonceCache) Purge(ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for nonce, timestamp := range c.entries { 
		if time.Since(timestamp) > ttl {      
			delete(c.entries, nonce)          
		}
	}
}
