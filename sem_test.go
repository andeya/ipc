package ipc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemctl(t *testing.T) {
	semid, err := Semget(1, 1, IPC_CREAT|1023)
	assert.NoError(t, err)
	ds, err := SemAllInfo(semid)
	assert.NoError(t, err)
	t.Log(ds)
}
