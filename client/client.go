package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/frozzare/go/http2"
	"github.com/frozzare/statscoll/stat"
	"github.com/pkg/errors"
)

// Client represents a new client.
type Client struct {
	client *http2.Client
	url    string
}

// New creates a new client.
func New(url string) *Client {
	return &Client{
		client: http2.NewClient(nil),
		url:    url,
	}
}

// Collect post a stat value to statscoll api.
func (c *Client) Collect(s *stat.Stat) error {
	buf, err := json.Marshal(s)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/collect", c.url), bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Wrap(errors.New(resp.Status), "statscoll")
	}

	return nil
}
