package modelsdev

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"
)

// bundledJSON is an offline snapshot of the models.dev catalog, embedded so the
// library works without network access. Refresh it with `make update-bundle`.
//
//go:embed data/models.json
var bundledJSON []byte

var (
	bundleOnce sync.Once
	bundleCat  *Catalog
	bundleErr  error
)

// Bundled returns the catalog parsed from the embedded offline snapshot. It
// never touches the network and is parsed once, then cached for the process.
func Bundled() (*Catalog, error) {
	bundleOnce.Do(func() {
		var raw map[string]Provider
		if err := json.Unmarshal(bundledJSON, &raw); err != nil {
			bundleErr = fmt.Errorf("modelsdev: parse bundled catalog: %w", err)
			return
		}
		bundleCat = newCatalog(raw)
	})
	return bundleCat, bundleErr
}
