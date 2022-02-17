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
	err = flock.ShareLock(true)
	assert.NoError(t, err)
	err = flock.ShareLock(true)
	assert.NoError(t, err)
	err = flock.UnlockAll()
	assert.NoError(t, err)
	err = flock.ExclusiveLock()
	assert.NoError(t, err)
	t.Log("Sleep...")
	time.Sleep(time.Hour)
}

func TestFlock2(t *testing.T) {
	flock, err := NewFlock("example")
	assert.NoError(t, err)
	t.Log(flock)
	err = flock.ExclusiveLock(true)
	assert.EqualError(t, err, "cannot flock pathectory example - resource temporarily unavailable")
	err = flock.ShareLock(true)
	assert.EqualError(t, err, "cannot flock pathectory example - resource temporarily unavailable")
	err = flock.ExclusiveLock()
	assert.NoError(t, err)
}
