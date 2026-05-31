package modelsdev

import (
	"encoding/json"
	"fmt"
)

// Parse builds a Catalog from a raw models.dev api.json document. Use it to
// load a snapshot you manage yourself (embedded via go:embed, read from disk,
// etc.) when you want catalog data without a network call.
func Parse(data []byte) (*Catalog, error) {
	var raw map[string]Provider
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("modelsdev: parse catalog: %w", err)
	}
	return newCatalog(raw), nil
}
