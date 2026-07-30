package feedschema

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

// TestVendoredSchemaParsesLiveFeed is a guard against silent drift: the copy
// compiling proves nothing about whether it still matches what the feed
// actually serves.
func TestVendoredSchemaParsesLiveFeed(t *testing.T) {
	c := &http.Client{Timeout: 20 * time.Second}
	resp, err := c.Get("https://feed.slash0.io/v1/index.json")
	if err != nil {
		t.Skipf("live feed unreachable: %v", err)
	}
	defer resp.Body.Close()
	// Skip rather than fail on a bad response: a CDN blip is not schema drift,
	// and this test exists to catch drift.
	if resp.StatusCode != http.StatusOK {
		t.Skipf("live feed returned %s", resp.Status)
	}

	var idx Index
	if err := json.NewDecoder(resp.Body).Decode(&idx); err != nil {
		t.Fatalf("decode live index: %v", err)
	}
	if len(idx.Services) == 0 {
		t.Fatal("live index decoded to zero services")
	}
	withPub := 0
	for _, s := range idx.Services {
		if s.Publication != nil {
			withPub++
			if s.Publication.DocumentType == "" || s.Publication.PollMode == "" {
				t.Errorf("%s: publication present but documentType/pollMode empty", s.Slug)
			}
		}
	}
	if withPub == 0 {
		t.Fatal("no service carried a publication block; the vendored schema is out of sync")
	}
	t.Logf("decoded %d services, %d with publication metadata", len(idx.Services), withPub)
}
