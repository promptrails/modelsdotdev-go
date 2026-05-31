package modelsdev

import (
	"context"
	"sync"
)

// defaultClient backs the package-level helpers. It is created lazily on first
// use with default options (live fetch, 1h TTL, offline fallback enabled).
var (
	defaultOnce   sync.Once
	defaultClient *Client
)

// Default returns the shared Client used by the package-level helpers. Callers
// who want different behavior should construct their own Client with New.
func Default() *Client {
	defaultOnce.Do(func() {
		defaultClient = New()
	})
	return defaultClient
}

// GetModelByID returns a model by its fully qualified ID in "provider:model"
// form, e.g. "anthropic:claude-opus-4-8", using the shared Default client.
func GetModelByID(ctx context.Context, qualifiedID string) (Model, error) {
	cat, err := Default().Catalog(ctx)
	if err != nil {
		return Model{}, err
	}
	return cat.Model(qualifiedID)
}

// GetProviderByName returns a provider by its display name (e.g. "OpenAI"),
// case-insensitively, using the shared Default client.
func GetProviderByName(ctx context.Context, name string) (Provider, error) {
	cat, err := Default().Catalog(ctx)
	if err != nil {
		return Provider{}, err
	}
	return cat.ProviderByName(name)
}
