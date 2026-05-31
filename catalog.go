package modelsdev

import (
	"sort"
	"strings"
)

// Catalog is an in-memory snapshot of the models.dev registry. It is immutable
// once built and safe for concurrent reads.
type Catalog struct {
	// Providers is keyed by provider ID (Provider.ID).
	Providers map[string]Provider
}

// newCatalog normalizes a decoded api.json map into a Catalog, backfilling the
// ID/Provider fields that models.dev keys implicitly by map position.
func newCatalog(raw map[string]Provider) *Catalog {
	providers := make(map[string]Provider, len(raw))
	for id, p := range raw {
		if p.ID == "" {
			p.ID = id
		}
		models := make(map[string]Model, len(p.Models))
		for mid, m := range p.Models {
			if m.ID == "" {
				m.ID = mid
			}
			if m.Provider == "" {
				m.Provider = p.ID
			}
			models[mid] = m
		}
		p.Models = models
		providers[p.ID] = p
	}
	return &Catalog{Providers: providers}
}

// Provider returns the provider with the given ID (e.g. "openai").
func (c *Catalog) Provider(id string) (Provider, error) {
	p, ok := c.Providers[id]
	if !ok {
		return Provider{}, ErrNotFound
	}
	return p, nil
}

// ProviderByName returns the provider whose Name matches name, case-insensitively
// (e.g. "OpenAI"). The first match in sorted-ID order is returned.
func (c *Catalog) ProviderByName(name string) (Provider, error) {
	for _, id := range c.ProviderIDs() {
		if strings.EqualFold(c.Providers[id].Name, name) {
			return c.Providers[id], nil
		}
	}
	return Provider{}, ErrNotFound
}

// Model returns a model by its fully qualified ID in "provider:model" form,
// e.g. "openai:gpt-4o" or "amazon-bedrock:anthropic.claude-sonnet".
func (c *Catalog) Model(qualifiedID string) (Model, error) {
	provider, model, ok := strings.Cut(qualifiedID, ":")
	if !ok || provider == "" || model == "" {
		return Model{}, ErrInvalidID
	}
	return c.ModelOf(provider, model)
}

// ModelOf returns a model by its provider ID and provider-local model ID.
func (c *Catalog) ModelOf(providerID, modelID string) (Model, error) {
	p, ok := c.Providers[providerID]
	if !ok {
		return Model{}, ErrNotFound
	}
	m, ok := p.Models[modelID]
	if !ok {
		return Model{}, ErrNotFound
	}
	return m, nil
}

// ProviderIDs returns all provider IDs in lexical order.
func (c *Catalog) ProviderIDs() []string {
	ids := make([]string, 0, len(c.Providers))
	for id := range c.Providers {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// Models returns every model in the catalog, sorted by provider ID then model
// ID for deterministic iteration. Each returned model has its ID and Provider
// fields populated.
func (c *Catalog) Models() []Model {
	var out []Model
	for _, pid := range c.ProviderIDs() {
		p := c.Providers[pid]
		ids := make([]string, 0, len(p.Models))
		for mid := range p.Models {
			ids = append(ids, mid)
		}
		sort.Strings(ids)
		for _, mid := range ids {
			out = append(out, p.Models[mid])
		}
	}
	return out
}

// Filter returns every model for which keep reports true, in the same
// deterministic order as Models.
func (c *Catalog) Filter(keep func(Model) bool) []Model {
	var out []Model
	for _, m := range c.Models() {
		if keep(m) {
			out = append(out, m)
		}
	}
	return out
}
