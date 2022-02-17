package ipc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var data = []byte("TestSyncShm1: Write1")

func TestSyncShm1(t *testing.T) {
	const a = 1 << 10
	syncShm, err := AttachSyncShm("example", 1<<10)
	assert.NoError(t, err)
	err = syncShm.Write(data)
	assert.NoError(t, err)
	data2 := syncShm.Read()
	assert.Equal(t, data, data2)
}

func TestSyncShm2(t *testing.T) {
	syncShm, err := AttachSyncShm("example", 0)
	assert.NoError(t, err)
	data2 := syncShm.Read()
	assert.Equal(t, data, data2)
	t.Log(string(data))
}

func TestSyncShm3(t *testing.T) {
	syncShm, err := AttachSyncShm("example", 0)
	assert.NoError(t, err)
	err = syncShm.Detach()
	assert.NoError(t, err)
	err = syncShm.Remove()
	assert.NoError(t, err)
}
