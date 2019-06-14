package main

import (
	"log"
	"reflect"

	"github.com/henrylee2cn/ipc"
)

var expected = "henrylee2cn"

func main() {
	key, err := ipc.Ftok("../../ipc.go", 2)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("key:%d", key)

	shmid, err := ipc.Shmget(key, 0, ipc.IPC_R)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := ipc.Shmctl(shmid, ipc.IPC_RMID)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// attach
	shmaddr, err := ipc.Shmat(shmid, ipc.SHM_REMAP)
	if err != nil {
		log.Fatal(err)
	}

	// read
	data := ipc.Shmread(shmaddr)
	actual := string(data)
	if !reflect.DeepEqual(expected, actual) {
		log.Fatalf("got:%#vwant:%#v", actual, expected)
	}
}
