// Package modelsdev is a Go client for the models.dev catalog — a community
// registry of LLM providers and models with their capabilities, token limits,
// and pricing.
//
// It mirrors the data exposed at https://models.dev/api.json with typed Go
// structs, an in-memory Catalog with lookup/filter helpers, and an HTTP Client
// with TTL caching.
//
// The quickest start uses the package-level helpers backed by a shared client:
//
//	model, err := modelsdev.GetModelByID(ctx, "anthropic:claude-opus-4-8")
//	provider, err := modelsdev.GetProviderByName(ctx, "OpenAI")
//
// For control over the HTTP client, URL, cache TTL, or offline fallback,
// construct a Client:
//
//	c := modelsdev.New(modelsdev.WithTTL(30*time.Minute), modelsdev.WithOfflineFallback(true))
//	cat, err := c.Catalog(ctx)
//	for _, m := range cat.Filter(func(m modelsdev.Model) bool { return m.Reasoning }) {
//		fmt.Println(m.Provider, m.ID)
//	}
//
// To work without a network call, parse a snapshot you manage yourself (read
// from disk, embedded via go:embed, …):
//
//	cat, err := modelsdev.Parse(mySnapshotBytes)
package modelsdev
