package main

import (
	"log"

	"github.com/henrylee2cn/ipc"
)

var expected = &ipc.Msgp{
	Mtype: 9,
	Mtext: []byte("henrylee2cn"),
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

	err = ipc.Msgsnd(msqid, expected, ipc.MSG_BLOCK)
	if err != nil {
		log.Fatal(err)
	}
}
