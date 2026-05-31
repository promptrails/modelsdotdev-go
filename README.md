# modelsdotdev-go

A Go client for the [models.dev](https://models.dev) catalog — a community
registry of LLM providers and models with their capabilities, token limits, and
pricing.

[![Go Reference](https://pkg.go.dev/badge/github.com/promptrails/modelsdotdev-go.svg)](https://pkg.go.dev/github.com/promptrails/modelsdotdev-go)
[![CI](https://github.com/promptrails/modelsdotdev-go/actions/workflows/ci.yml/badge.svg)](https://github.com/promptrails/modelsdotdev-go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/promptrails/modelsdotdev-go)](https://goreportcard.com/report/github.com/promptrails/modelsdotdev-go)

Typed structs over `https://models.dev/api.json`, an in-memory catalog with
lookup/filter helpers, an HTTP client with TTL caching, and an embedded offline
snapshot so it works with no network at all.

```go
import "github.com/promptrails/modelsdotdev-go"

model, _ := modelsdev.GetModelByID(ctx, "anthropic:claude-opus-4-8")
fmt.Println(model.Name, *model.Cost.Input) // Claude Opus 4.8 5

provider, _ := modelsdev.GetProviderByName(ctx, "OpenAI")
for id := range provider.Models {
    fmt.Println(id)
}
```

## Install

```bash
go get github.com/promptrails/modelsdotdev-go
```

Requires Go 1.26+. No third-party dependencies — standard library only.

## Usage

### Package-level helpers

Backed by a shared client (live fetch, 1h cache, offline fallback):

```go
model, err := modelsdev.GetModelByID(ctx, "openai:gpt-4o")
provider, err := modelsdev.GetProviderByName(ctx, "Anthropic")
```

### A configured client

```go
c := modelsdev.New(
    modelsdev.WithTTL(30*time.Minute),
    modelsdev.WithOfflineFallback(true),
    modelsdev.WithUserAgent("my-app/1.0"),
)

cat, err := c.Catalog(ctx) // fetched once, then cached until the TTL expires

// Every reasoning-capable model that also reports a cache-read price:
for _, m := range cat.Filter(func(m modelsdev.Model) bool {
    return m.Reasoning && m.Cost.CacheRead != nil
}) {
    fmt.Printf("%s:%s\n", m.Provider, m.ID)
}
```

`Catalog` serves the last good snapshot if a refetch fails, and (with
`WithOfflineFallback`) falls back to the embedded snapshot when there is nothing
cached yet. `Refresh` forces a refetch.

### Fully offline

```go
cat, err := modelsdev.Bundled() // embedded snapshot, never touches the network
```

## Catalog API

| Method | Returns |
| ------ | ------- |
| `Provider(id)` | a provider by ID (`"openai"`) |
| `ProviderByName(name)` | a provider by display name, case-insensitive (`"OpenAI"`) |
| `Model("provider:model")` | a model by fully qualified ID |
| `ModelOf(provider, model)` | a model by provider ID + model ID |
| `ProviderIDs()` | all provider IDs, sorted |
| `Models()` | every model, sorted by provider then model ID |
| `Filter(func(Model) bool)` | every model matching a predicate |

Missing lookups return `ErrNotFound`; a malformed qualified ID returns
`ErrInvalidID`.

## Data model

Each `Model` carries its capabilities (`Reasoning`, `ToolCall`, `Attachment`,
`StructuredOutput`, `Modalities`), `Limit` (context/output tokens), and `Cost`
(per-1M-token `Input`, `Output`, `CacheRead`, `CacheWrite`, `Reasoning`, audio,
and long-context tiers). Cost fields are pointers: `nil` means "not reported",
which is distinct from a reported zero.

## Updating the offline snapshot

```bash
make update-bundle   # curls https://models.dev/api.json into data/models.json
```

## License

MIT — see [LICENSE](LICENSE).
