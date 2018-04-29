package cache

import "testing"

func TestCache(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Set("test", map[string]interface{}{"test": "test"}); err != nil {
		t.Fatal(err)
	}

	v, err := c.Get("test")
	if err != nil {
		t.Fatal(err)
	}

	if v.(map[string]interface{})["test"].(string) != "test" {
		t.Fatal("Expected test value to be test")
	}

	if err := c.RemovePrefix("test"); err != nil {
		t.Fatal(err)
	}

	if _, err := c.Get("test"); err == nil {
		t.Fatal(err)
	}
}
