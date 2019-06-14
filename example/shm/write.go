package main

import (
	"log"

	"github.com/henrylee2cn/ipc"
)

var expected = "henrylee2cn"

func main() {
	key, err := ipc.Ftok("../../ipc.go", 2)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("key:%d", key)

	shmid, err := ipc.Shmget(key, 32, ipc.IPC_CREAT|ipc.IPC_RW)
	if err != nil {
		log.Fatal(err)
	}

	// attach
	shmaddr, err := ipc.Shmat(shmid, ipc.SHM_REMAP)
	if err != nil {
		log.Fatal(err)
	}

	// write
	ipc.Shmwrite(shmaddr, []byte(expected))

	// detach
	err = ipc.Shmdt(shmaddr)
	if err != nil {
		log.Fatal(err)
	}
}
