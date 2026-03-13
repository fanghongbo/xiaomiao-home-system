package store

import (
	"errors"
	"sync"
)

var (
	ErrNotFountStoreKey = errors.New("store key not found")
)

type MapStore[K string, V any] struct {
	l sync.RWMutex
	m map[K]V
}

// Load 加载
func (u *MapStore[K, V]) Load(key K) (V, error) {
	var (
		item V
		ok   bool
	)

	u.l.RLock()
	defer u.l.RUnlock()

	if u.m != nil {
		if item, ok = u.m[key]; ok {
			return item, nil
		}
	}

	return item, ErrNotFountStoreKey
}

// Store 缓存
func (u *MapStore[K, V]) Store(key K, item V) {
	u.l.Lock()
	defer u.l.Unlock()

	if u.m == nil {
		u.m = make(map[K]V)
	}

	u.m[key] = item
}
