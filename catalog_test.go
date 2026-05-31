package modelsdev

import (
	"encoding/json"
	"errors"
	"testing"
)

const sampleJSON = `{
  "openai": {
    "id": "openai",
    "name": "OpenAI",
    "env": ["OPENAI_API_KEY"],
    "models": {
      "gpt-4o": {
        "name": "GPT-4o",
        "reasoning": false,
        "tool_call": true,
        "modalities": {"input": ["text", "image"], "output": ["text"]},
        "limit": {"context": 128000, "output": 16384},
        "cost": {"input": 2.5, "output": 10, "cache_read": 1.25}
      }
    }
  },
  "anthropic": {
    "name": "Anthropic",
    "models": {
      "claude-opus-4-8": {
        "name": "Claude Opus 4.8",
        "reasoning": true,
        "tool_call": true,
        "modalities": {"input": ["text", "image", "pdf"], "output": ["text"]},
        "limit": {"context": 1000000, "output": 128000},
        "cost": {"input": 5, "output": 25, "cache_read": 0.5}
      }
    }
  }
}`

func testCatalog(t *testing.T) *Catalog {
	t.Helper()
	var raw map[string]Provider
	if err := json.Unmarshal([]byte(sampleJSON), &raw); err != nil {
		t.Fatalf("unmarshal sample: %v", err)
	}
	return newCatalog(raw)
}

func TestNewCatalogBackfillsIDs(t *testing.T) {
	c := testCatalog(t)
	// anthropic has no explicit "id" in the sample → backfilled from the map key.
	p, err := c.Provider("anthropic")
	if err != nil {
		t.Fatalf("Provider: %v", err)
	}
	if p.ID != "anthropic" {
		t.Errorf("provider ID = %q, want anthropic", p.ID)
	}
	m := p.Models["claude-opus-4-8"]
	if m.ID != "claude-opus-4-8" {
		t.Errorf("model ID = %q, want claude-opus-4-8", m.ID)
	}
	if m.Provider != "anthropic" {
		t.Errorf("model Provider = %q, want anthropic", m.Provider)
	}
}

func TestProviderByName(t *testing.T) {
	c := testCatalog(t)
	p, err := c.ProviderByName("openai") // case-insensitive
	if err != nil {
		t.Fatalf("ProviderByName: %v", err)
	}
	if p.ID != "openai" {
		t.Errorf("got %q, want openai", p.ID)
	}
	if _, err := c.ProviderByName("nope"); !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestModelQualifiedID(t *testing.T) {
	c := testCatalog(t)
	m, err := c.Model("openai:gpt-4o")
	if err != nil {
		t.Fatalf("Model: %v", err)
	}
	if m.Cost.CacheRead == nil || *m.Cost.CacheRead != 1.25 {
		t.Errorf("cache_read not parsed: %+v", m.Cost.CacheRead)
	}
	if !m.SupportsVision() {
		t.Error("gpt-4o should support vision")
	}

	for _, bad := range []string{"openai", "openai:", ":gpt-4o", ""} {
		if _, err := c.Model(bad); !errors.Is(err, ErrInvalidID) {
			t.Errorf("Model(%q): want ErrInvalidID, got %v", bad, err)
		}
	}
	if _, err := c.Model("openai:ghost"); !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestModelsDeterministicOrder(t *testing.T) {
	c := testCatalog(t)
	models := c.Models()
	if len(models) != 2 {
		t.Fatalf("got %d models, want 2", len(models))
	}
	// Sorted by provider ID: anthropic before openai.
	if models[0].Provider != "anthropic" || models[1].Provider != "openai" {
		t.Errorf("order = %s,%s; want anthropic,openai", models[0].Provider, models[1].Provider)
	}
}

func TestFilter(t *testing.T) {
	c := testCatalog(t)
	reasoning := c.Filter(func(m Model) bool { return m.Reasoning })
	if len(reasoning) != 1 || reasoning[0].ID != "claude-opus-4-8" {
		t.Errorf("reasoning filter = %+v", reasoning)
	}
}
