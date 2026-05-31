package modelsdev

import "testing"

func TestParse(t *testing.T) {
	cat, err := Parse([]byte(sampleJSON))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(cat.Providers) != 2 {
		t.Fatalf("got %d providers, want 2", len(cat.Providers))
	}
	m, err := cat.Model("anthropic:claude-opus-4-8")
	if err != nil {
		t.Fatalf("Model: %v", err)
	}
	// IDs and provider back-reference are backfilled during parse.
	if m.ID != "claude-opus-4-8" || m.Provider != "anthropic" {
		t.Errorf("backfill failed: id=%q provider=%q", m.ID, m.Provider)
	}
}

func TestParseInvalid(t *testing.T) {
	if _, err := Parse([]byte("not json")); err == nil {
		t.Fatal("want error for invalid JSON, got nil")
	}
}
