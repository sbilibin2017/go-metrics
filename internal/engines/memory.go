package engines

import (
	"sync"
)

type MemoryStorage[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

func NewMemoryStorage[K comparable, V any]() *MemoryStorage[K, V] {
	return &MemoryStorage[K, V]{
		data: make(map[K]V),
	}
}

type MemorySetter[K comparable, V any] struct {
	*MemoryStorage[K, V]
}

func NewMemorySetter[K comparable, V any](storage *MemoryStorage[K, V]) *MemorySetter[K, V] {
	return &MemorySetter[K, V]{MemoryStorage: storage}
}

func (m *MemorySetter[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

type MemoryGetter[K comparable, V any] struct {
	*MemoryStorage[K, V]
}

func NewMemoryGetter[K comparable, V any](storage *MemoryStorage[K, V]) *MemoryGetter[K, V] {
	return &MemoryGetter[K, V]{MemoryStorage: storage}
}

func (m *MemoryGetter[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, exists := m.data[key]
	return val, exists
}

type MemoryRanger[K comparable, V any] struct {
	*MemoryStorage[K, V]
}

func NewMemoryRanger[K comparable, V any](storage *MemoryStorage[K, V]) *MemoryRanger[K, V] {
	return &MemoryRanger[K, V]{MemoryStorage: storage}
}

func (m *MemoryRanger[K, V]) Range(f func(K, V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for key, value := range m.data {
		if !f(key, value) {
			break
		}
	}
}
