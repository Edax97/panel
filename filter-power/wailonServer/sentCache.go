package wailonServer

import (
	"encoding/gob"
	"os"
	"time"
)

type SentCache struct {
	sentMap  map[string]time.Time
	diskPath string
}

func NewSentCache(diskPath string) *SentCache {
	cache := &SentCache{
		diskPath: diskPath,
		sentMap:  make(map[string]time.Time),
	}
	cache.loadCache()
	return cache
}

func (c *SentCache) saveCache() {
	f, err := os.Create(c.diskPath)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()
	encoder := gob.NewEncoder(f)
	err = encoder.Encode(c.sentMap)
	if err != nil {
		return
	}

}

func (c *SentCache) loadCache() {
	f, err := os.Open(c.diskPath)
	if err != nil {
		c.sentMap = make(map[string]time.Time)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	var data map[string]time.Time
	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&data)
	if err != nil {
		c.sentMap = make(map[string]time.Time)
		return
	}
	c.sentMap = data
}

func (c *SentCache) hasSent(imei string, t time.Time) bool {
	sent, ok := c.sentMap[imei]
	if !ok {
		return false
	}
	if t.Before(sent) {
		return true
	}
	return false
}

func (c *SentCache) updateSent(imei string, sent time.Time) bool {
	c.sentMap[imei] = sent
	c.saveCache()
	return true
}
