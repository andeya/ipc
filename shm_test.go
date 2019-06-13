package ipc_test

import (
	"testing"

	"github.com/henrylee2cn/ipc"
	"github.com/stretchr/testify/assert"
)

func TestShm(t *testing.T) {
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

	// write
	expected := ipc.Msgp{
		Mtype: 9,
		Mtext: []byte("henrylee2cn"),
	}
	*(*ipc.Msgp)(shmaddr) = expected

	// detach
	err = ipc.Shmdt(shmaddr)
	if err != nil {
		t.Fatal(err)
	}

	// attach
	shmaddr, err = ipc.Shmat(shmid, ipc.SHM_REMAP)
	if err != nil {
		t.Fatal(err)
	}

	// read
	actual := *(*ipc.Msgp)(shmaddr)

	assert.Equal(t, expected, actual)
}
