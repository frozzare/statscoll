package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frozzare/statscoll/stat"
)

func TestClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"success":true}`))
	}))

	defer server.Close()

	c := New(server.URL)

	err := c.Collect(&stat.Stat{
		Count:  23,
		Metric: "test",
	})

	if err != nil {
		t.Errorf("Expected: nil, got: %s", err)
	}
}
