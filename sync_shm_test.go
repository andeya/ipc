package ipc

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var data = []byte("TestSyncShm1: Write1")

const sizeIfCreate = 1 << 10

func TestSyncShm1(t *testing.T) {
	const a = 1 << 10
	syncShm, err := AttachSyncShm("example", sizeIfCreate)
	assert.NoError(t, err)
	err = syncShm.Write(data)
	assert.NoError(t, err)
	data2 := syncShm.Read()
	assert.Equal(t, data, data2)
}

func TestSyncShm2(t *testing.T) {
	syncShm, err := AttachSyncShm("example", sizeIfCreate)
	assert.NoError(t, err)
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data2 := syncShm.Read()
			assert.Equal(t, data, data2)
		}()
	}
	wg.Wait()
}

func TestSyncShm3(t *testing.T) {
	syncShm, err := AttachSyncShm("example", sizeIfCreate)
	assert.NoError(t, err)
	err = syncShm.Detach()
	assert.NoError(t, err)
	err = syncShm.Remove()
	assert.NoError(t, err)
}
