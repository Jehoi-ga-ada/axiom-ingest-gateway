package service

import (
	"sync"

	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/auth/domain"
)

type MemoryKeyRegistry struct {
	mu sync.RWMutex
	keys map[string]*domain.APIKeyMetadata
}

func NewMemoryKeyRegistry() KeyRegistry {
	return &MemoryKeyRegistry{
		keys: make(map[string]*domain.APIKeyMetadata),
	}
}

func (r *MemoryKeyRegistry) Get(key []byte) (*domain.APIKeyMetadata, bool) {
    r.mu.RLock()
    meta, ok := r.keys[string(key)]
    r.mu.RUnlock()
	if !ok {
		return nil, false
	}
    return meta, true
}

func (r *MemoryKeyRegistry) Upsert(key string, meta domain.APIKeyMetadata) {
    r.mu.Lock()
    r.keys[key] = &meta
    r.mu.Unlock()
}