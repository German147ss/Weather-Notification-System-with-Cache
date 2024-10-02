package adapters

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type MemoryCacheService struct {
	data  map[string]string
	mutex sync.RWMutex
}

func NewMemoryCacheService() *MemoryCacheService {
	return &MemoryCacheService{
		data: make(map[string]string),
	}
}

func (m *MemoryCacheService) Get(ctx context.Context, key string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	val, ok := m.data[key]
	if !ok {
		return "", redis.Nil
	}
	return val, nil
}

func (m *MemoryCacheService) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Ignore the expiration for the in-memory store
	m.data[key] = value
	return nil
}
