package config

import "testing"

func TestReadFile(t *testing.T) {
	c, err := ReadFile("./config.yml")

	if err != nil {
		t.Errorf("Expected: nil, got: %v", err)
	}

	if c == nil {
		t.Errorf("Expected: struct, got: %v", c)
	}

	if c.Port != 9300 {
		t.Errorf("Expected: 9300, got: %v", c.Port)
	}
}
