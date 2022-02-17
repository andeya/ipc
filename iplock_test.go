package ipc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIpLock(t *testing.T) {
	ipLock, err := NewIPLock("example")
	assert.NoError(t, err)
	ipLock.RLock()
	ipLock.RLock()
	ipLock.RUnlock()
	ipLock.RUnlock()
	ipLock.Lock()
	ipLock.Unlock()
}
