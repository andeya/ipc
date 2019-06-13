package ipc_test

import (
	"testing"

	"github.com/henrylee2cn/ipc"
	"github.com/stretchr/testify/assert"
)

func TestMsq(t *testing.T) {
	key, err := ipc.Ftok("ipc.go", 1)
	assert.Nil(t, err)

	msqid, err := ipc.Msgget(key, ipc.IPC_CREAT|ipc.ACCESS_RDWR)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := ipc.Msgctl(msqid, ipc.IPC_RMID)
		if err != nil {
			t.Fatal(err)
		}
	}()

	var expected = &ipc.Msgp{
		Mtype: 9,
		Mtext: []byte("henrylee2cn"),
	}

	err = ipc.Msgsnd(msqid, expected, ipc.MSG_BLOCK)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := ipc.Msgrcv(msqid, 0, ipc.IPC_NOWAIT)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, actual)
}