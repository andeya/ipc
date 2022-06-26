package main

import (
	"log"

	"github.com/andeya/ipc"
)

var expected = "andeya"

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
	err = ipc.Shmwrite(shmaddr, []byte(expected))
	if err != nil {
		log.Fatal(err)
	}

	// detach
	err = ipc.Shmdt(shmaddr)
	if err != nil {
		log.Fatal(err)
	}
}
