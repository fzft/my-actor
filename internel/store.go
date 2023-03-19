package internel

import "github.com/fzft/my-actor/pkg"

// Storer is the interface that wraps the basic Store methods.
type Storer interface {
	Put(key string, value any)
	Get(key string) (any, bool)
}

// MemoryStore wraps the KeyValueStore
type MemoryStore struct {
	store *pkg.KeyValueStore[string, any]
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{store: pkg.NewKeyValueStore[string, any]()}
}

// Put ...
func (m *MemoryStore) Put(key string, value any) {
	m.store.Put(key, value)
}

// Get ...
func (m *MemoryStore) Get(key string) (any, bool) {
	return m.store.Get(key)
}

// Pool ...
type Pool interface {
	Put(key any, value any)
	Get(key any) []any
}
