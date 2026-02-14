package service

import "github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/auth/domain"

type KeyRegistry interface {
	Get(key []byte) (*domain.APIKeyMetadata, bool)
	Upsert(key string, metadata domain.APIKeyMetadata)
}