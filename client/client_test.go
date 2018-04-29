package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"success":true}`))
	}))

	defer server.Close()

	c := New(server.URL)

	err := c.Collect(Stat{
		Metric: "test",
		Value:  23.0,
	})

	if err != nil {
		t.Errorf("Expected: nil, got: %s", err)
	}
}
