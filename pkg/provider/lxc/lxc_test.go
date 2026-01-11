package lxc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLXCProvider(t *testing.T) {
	provider, err := NewLXCProvider()
	if err != nil {
		t.Skip("LXC not available:", err)
	}

	assert.NotNil(t, provider)
	assert.Equal(t, "lxc", provider.Name())
}

func TestLXCProvider_Name(t *testing.T) {
	provider, err := NewLXCProvider()
	if err != nil {
		t.Skip("LXC not available:", err)
	}

	assert.Equal(t, "lxc", provider.Name())
}
