package ipc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expected = &Msgp{
	Mtype: 9,
	Mtext: []byte("henrylee2cn"),
}

func TestMsq_snd(t *testing.T) {
	key, err := Ftok("go", 2)
	assert.Nil(t, err)
	t.Log(key)

	msqid, err := Msgget(key, IPC_CREAT|IPC_RW)
	if err != nil {
		t.Fatal(err)
	}

	err = Msgsnd(msqid, expected, MSG_BLOCK)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMsq_recv(t *testing.T) {
	key, err := Ftok("go", 2)
	assert.Nil(t, err)
	t.Log(key)

	msqid, err := Msgget(key, IPC_R)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := Msgctl(msqid, IPC_RMID)
		if err != nil {
			t.Fatal(err)
		}
	}()

	actual, err := Msgrcv(msqid, 0, IPC_NOWAIT)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, actual)
}

func BenchmarkMsgp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var expected = Msgp{
			Mtype: 9,
			Mtext: []byte("henrylee2cn"),
		}
		ptr, textSize := expected.marshal()
		if len(expected.Mtext) != textSize {
			b.Fatalf("expected:%v, actual:%v", len(expected.Mtext), textSize)
		}
		var actual Msgp
		err := actual.unmarshal(textSize, ptr)
		if err != nil {
			b.Fatal(err)
		}
		if !reflect.DeepEqual(expected, actual) {
			b.Fatalf("expected:%v, actual:%v", expected, actual)
		}
	}
}
