package ipc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFlock1(t *testing.T) {
	flock, err := NewFlock("example")
	assert.NoError(t, err)
	t.Log(flock)
	err = flock.Lock()
	assert.NoError(t, err)
	time.Sleep(time.Hour)
}

func TestFlock2(t *testing.T) {
	flock, err := NewFlock("example")
	assert.NoError(t, err)
	t.Log(flock)
	err = flock.Lock(true)
	assert.EqualError(t, err, "cannot flock pathectory example - resource temporarily unavailable")
	err = flock.RLock(true)
	assert.EqualError(t, err, "cannot flock pathectory example - resource temporarily unavailable")
	err = flock.Lock()
	assert.NoError(t, err)
}
