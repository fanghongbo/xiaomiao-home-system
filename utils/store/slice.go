package store

import (
	"sync"
)

type SliceStore[T comparable] struct {
	l     sync.RWMutex
	items []T
}

// Push 插入数据
func (u *SliceStore[T]) Push(item T) {
	u.l.Lock()
	defer u.l.Unlock()

	u.items = append(u.items, item)
}

// Pop 取出数据
func (u *SliceStore[T]) Pop() (T, bool) {
	var item T

	u.l.Lock()
	defer u.l.Unlock()

	if len(u.items) == 0 {
		return item, false
	}

	item = u.items[len(u.items)-1]

	u.items = u.items[:len(u.items)-1]

	u.items = append(u.items, item)

	return item, true
}
