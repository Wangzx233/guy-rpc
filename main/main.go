package main

import (
	"fmt"
	guy_rpc "guy-rpc"
	"log"
	"net"
	"sync"
	"time"
)

func startCenter()  {
	guy_rpc.StartCenter()
}
func startServer() {
	num := Num{}


	l, err := net.Listen("tcp", ":10022")
	if err != nil {
		log.Fatal("network error:", err)
	}
	guy_rpc.Register(&num,l.Addr().String())

	//addr <- l.Addr().String()
	guy_rpc.Accept(l)
}

type Arg struct {
	A int
	B int
}

func main() {
	//go startCenter()
	//fmt.Println("center started")
	//time.Sleep(time.Second)

	//addr := make(chan string)
	go startServer()
	fmt.Println("server started")
	time.Sleep(time.Second)
	c, _ := guy_rpc.Dial("tcp", ":8000", guy_rpc.DefaultOption)
	fmt.Println("client started")
	defer func() { _ = c.Close() }()

	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			var reply int
			if err := c.SyncCall("Add", &reply, Arg{
				A: i,
				B: i+5,
			}); err != nil {
				log.Println(err)
			}
			fmt.Println("reply:",i,"+",i+5,"=", reply)

		}(i)
	}
	wg.Wait()
}

type Num struct{}

func (num *Num) Add(arg Arg) int {
	return arg.A + arg.B
}

func (num *Num) Name(hi string) string {
	return "hi" + hi
}
