package ipc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedShm = []byte("henrylee2cn")

func TestShmwrite(t *testing.T) {
	key, err := Ftok("ipc.go", 5)
	assert.Nil(t, err)

	shmid, err := Shmget(key, 32, IPC_CREAT|IPC_RW)
	if err != nil {
		t.Fatal(err)
	}

	// attach
	shmaddr, err := Shmat(shmid, SHM_REMAP)
	if err != nil {
		t.Fatal(err)
	}

	// write
	err = Shmwrite(shmaddr, expectedShm)
	if err != nil {
		t.Fatal(err)
	}

	// detach
	err = Shmdt(shmaddr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestShmread(t *testing.T) {
	key, err := Ftok("ipc.go", 5)
	assert.Nil(t, err)

	shmid, err := Shmget(key, 0, IPC_R)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := Shmctl(shmid, IPC_RMID)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// attach
	shmaddr, err := Shmat(shmid, SHM_REMAP)
	if err != nil {
		t.Fatal(err)
	}

	// read
	actual := Shmread(shmaddr)
	assert.Equal(t, expectedShm, actual)
}
