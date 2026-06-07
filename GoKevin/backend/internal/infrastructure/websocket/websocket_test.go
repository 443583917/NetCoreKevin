package websocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHub_NewHub(t *testing.T) {
	hub := NewHub()
	assert.NotNil(t, hub)
	assert.NotNil(t, hub.clients)
	assert.NotNil(t, hub.Broadcast)
	assert.NotNil(t, hub.Register)
	assert.NotNil(t, hub.Unregister)
}

func TestHub_GetClientCount(t *testing.T) {
	hub := NewHub()
	assert.Equal(t, 0, hub.GetClientCount())
}
