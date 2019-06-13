package main

import (
	"time"

	"github.com/henrylee2cn/goutil/ipc"
)

type Data struct {
	a *string
	b int
	c chan int
}

func main() {
	ipc, err := ipc.ShareByFtok("../ipc.go", 1, 100)
	if err != nil {
		panic(err)
	}
	defer ipc.Detach()
	a := "111"
	data := Data{
		a: &a,
		b: 12,
		c: make(chan int, 1),
	}
	data.c <- 3
	*(*Data)(ipc.Pointer()) = data
	time.Sleep(10000e9)
}
