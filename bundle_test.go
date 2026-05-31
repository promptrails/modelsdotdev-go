package modelsdev

import "testing"

func TestBundledParses(t *testing.T) {
	cat, err := Bundled()
	if err != nil {
		t.Fatalf("Bundled: %v", err)
	}
	if len(cat.Providers) == 0 {
		t.Fatal("bundled catalog has no providers")
	}
	// Sanity: a couple of well-known providers should be present.
	for _, id := range []string{"openai", "anthropic"} {
		if _, err := cat.Provider(id); err != nil {
			t.Errorf("bundled catalog missing provider %q: %v", id, err)
		}
	}
	// Every model should have its ID and Provider backfilled.
	for _, m := range cat.Models() {
		if m.ID == "" || m.Provider == "" {
			t.Fatalf("model with empty ID/Provider: %+v", m)
		}
	}
}
