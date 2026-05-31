package modelsdev

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestClientCatalogCaches(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleJSON))
	}))
	defer srv.Close()

	c := New(WithURL(srv.URL), WithTTL(time.Hour))
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		cat, err := c.Catalog(ctx)
		if err != nil {
			t.Fatalf("Catalog: %v", err)
		}
		if _, err := cat.Model("openai:gpt-4o"); err != nil {
			t.Fatalf("lookup: %v", err)
		}
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Errorf("server hit %d times, want 1 (cached)", got)
	}
}

func TestClientRefresh(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&hits, 1)
		_, _ = w.Write([]byte(sampleJSON))
	}))
	defer srv.Close()

	c := New(WithURL(srv.URL))
	ctx := context.Background()
	if _, err := c.Catalog(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Refresh(ctx); err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadInt32(&hits); got != 2 {
		t.Errorf("server hit %d times, want 2 (refresh refetches)", got)
	}
}

func TestClientServesStaleOnError(t *testing.T) {
	var fail atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if fail.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(sampleJSON))
	}))
	defer srv.Close()

	c := New(WithURL(srv.URL), WithTTL(0)) // TTL 0 → always tries to refetch
	ctx := context.Background()
	if _, err := c.Catalog(ctx); err != nil {
		t.Fatalf("initial: %v", err)
	}
	fail.Store(true)
	cat, err := c.Catalog(ctx)
	if err != nil {
		t.Fatalf("want stale served, got error: %v", err)
	}
	if _, err := cat.Model("openai:gpt-4o"); err != nil {
		t.Errorf("stale catalog unusable: %v", err)
	}
}

func TestClientErrorNoCache(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	c := New(WithURL(srv.URL))
	if _, err := c.Catalog(context.Background()); err == nil {
		t.Fatal("want error on first fetch failure, got nil")
	}
}

func TestClientOfflineFallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := New(WithURL(srv.URL), WithOfflineFallback(true))
	cat, err := c.Catalog(context.Background())
	if err != nil {
		t.Fatalf("offline fallback should succeed: %v", err)
	}
	if len(cat.Providers) == 0 {
		t.Error("offline catalog is empty")
	}
}
