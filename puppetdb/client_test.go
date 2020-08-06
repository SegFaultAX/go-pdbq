package puppetdb

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var hosts = []Host{
	{
		Name: "test-1",
	},
	{
		Name: "test-2",
	},
	{
		Name: "test-3",
	},
}

func TestHosts(t *testing.T) {
	assert := assert.New(t)
	s := newServer(t)
	c, _ := NewClient(s.URL)

	hs, err := c.Hosts("test-1")
	if err != nil {
		t.Fatal("unexpected error fetching hosts", err)
	}

	assert.ElementsMatch(hosts, hs[0:1])
}

func newServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pdb/query/v4" {
			t.Error("expected PQL path, got", r.URL.Path)
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(hosts)
	}))
}
