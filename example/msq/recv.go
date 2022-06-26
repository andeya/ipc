package main

import (
	"log"
	"reflect"

	"github.com/andeya/ipc"
)

var expected = &ipc.Msgp{
	Mtype: 9,
	Mtext: []byte("andeya"),
}

func main() {
	key, err := ipc.Ftok("../../ipc.go", 1)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("key:%d", key)

	msqid, err := ipc.Msgget(key, ipc.IPC_CREAT|ipc.IPC_RW)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := ipc.Msgctl(msqid, ipc.IPC_RMID)
		if err != nil {
			log.Fatal(err)
		}
	}()

	actual, err := ipc.Msgrcv(msqid, 9, ipc.MSG_BLOCK)
	if err != nil {
		log.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		log.Fatalf("got:%#vwant:%#v", actual, expected)
	}
}
