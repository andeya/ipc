package ipc_test

import (
	"testing"

	"github.com/henrylee2cn/ipc"
	"github.com/stretchr/testify/assert"
)

var expectedShm = []byte("henrylee2cn")

func TestShmwrite(t *testing.T) {
	key, err := ipc.Ftok("ipc.go", 2)
	assert.Nil(t, err)

	shmid, err := ipc.Shmget(key, 32, ipc.IPC_CREAT|ipc.IPC_RW)
	if err != nil {
		t.Fatal(err)
	}

	// attach
	shmaddr, err := ipc.Shmat(shmid, ipc.SHM_REMAP)
	if err != nil {
		t.Fatal(err)
	}

	// write
	ipc.Shmwrite(shmaddr, expectedShm)

	// detach
	err = ipc.Shmdt(shmaddr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestShmread(t *testing.T) {
	key, err := ipc.Ftok("ipc.go", 2)
	assert.Nil(t, err)

	shmid, err := ipc.Shmget(key, 32, ipc.IPC_CREAT|ipc.IPC_RW)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := ipc.Shmctl(shmid, ipc.IPC_RMID)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// attach
	shmaddr, err := ipc.Shmat(shmid, ipc.SHM_REMAP)
	if err != nil {
		t.Fatal(err)
	}

	// read
	actual := ipc.Shmread(shmaddr)
	assert.Equal(t, expectedShm, actual)
}
