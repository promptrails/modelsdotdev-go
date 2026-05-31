package modelsdev

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// DefaultAPIURL is the models.dev machine-readable catalog endpoint.
const DefaultAPIURL = "https://models.dev/api.json"

// defaultTTL is how long a fetched catalog is cached before a refetch.
const defaultTTL = time.Hour

// Client fetches and caches the models.dev catalog. It is safe for concurrent
// use. The zero value is not usable — construct one with New.
type Client struct {
	httpClient      *http.Client
	url             string
	userAgent       string
	ttl             time.Duration
	offlineFallback bool

	mu        sync.Mutex
	cache     *Catalog
	fetchedAt time.Time
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets the HTTP client used for fetching. Defaults to a client
// with a 20s timeout.
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		if h != nil {
			c.httpClient = h
		}
	}
}

// WithURL overrides the catalog URL. Defaults to DefaultAPIURL.
func WithURL(url string) Option {
	return func(c *Client) {
		if url != "" {
			c.url = url
		}
	}
}

// WithUserAgent sets the User-Agent header sent with catalog requests.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		if ua != "" {
			c.userAgent = ua
		}
	}
}

// WithTTL sets how long a fetched catalog is cached before it is refetched. A
// non-positive ttl disables caching (every call refetches).
func WithTTL(ttl time.Duration) Option {
	return func(c *Client) { c.ttl = ttl }
}

// WithOfflineFallback makes Catalog fall back to the embedded snapshot (see
// Bundled) when a live fetch fails and nothing is cached yet. This trades
// freshness for availability.
func WithOfflineFallback(enabled bool) Option {
	return func(c *Client) { c.offlineFallback = enabled }
}

// New constructs a Client with the given options.
func New(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: 20 * time.Second},
		url:        DefaultAPIURL,
		userAgent:  "modelsdotdev-go",
		ttl:        defaultTTL,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Catalog returns the models.dev catalog, fetching it on first use and on cache
// expiry. Concurrent callers share a single in-flight result via the lock.
func (c *Client) Catalog(ctx context.Context) (*Catalog, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache != nil && c.ttl > 0 && time.Since(c.fetchedAt) < c.ttl {
		return c.cache, nil
	}

	cat, err := c.fetch(ctx)
	if err != nil {
		if c.cache != nil {
			return c.cache, nil // serve stale rather than fail
		}
		if c.offlineFallback {
			if b, berr := Bundled(); berr == nil {
				c.cache = b
				return b, nil
			}
		}
		return nil, err
	}

	c.cache = cat
	c.fetchedAt = c.now()
	return cat, nil
}

// Refresh forces a fetch, replacing any cached catalog on success.
func (c *Client) Refresh(ctx context.Context) (*Catalog, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cat, err := c.fetch(ctx)
	if err != nil {
		return nil, err
	}
	c.cache = cat
	c.fetchedAt = c.now()
	return cat, nil
}

// now is overridable in tests; defaults to time.Now.
func (c *Client) now() time.Time { return time.Now() }

func (c *Client) fetch(ctx context.Context) (*Catalog, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("modelsdev: fetch catalog: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("modelsdev: catalog returned status %d", resp.StatusCode)
	}

	var raw map[string]Provider
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("modelsdev: decode catalog: %w", err)
	}
	return newCatalog(raw), nil
}
