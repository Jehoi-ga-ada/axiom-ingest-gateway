package domain

import "github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/shared/domain"

type APIKeyMetadata struct {
	TenantID domain.TenantID
	Active bool
}