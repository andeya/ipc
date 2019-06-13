package main

import (
	"fmt"

	"github.com/henrylee2cn/goutil/ipc"
)

type Data struct {
	a *string
	b int
	c chan int
}

func main() {
	ipc, err := ipc.ShareByFtok("../ipc.go", 1, 0)
	if err != nil {
		panic(err)
	}
	defer ipc.Detach()

	data := *(*Data)(ipc.Pointer())

	fmt.Printf("a:%v, b:%v, c:%v", *data.a, data.b, data.c)
}
