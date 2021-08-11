package main

import (
	guy_rpc "guy-rpc"
	"guy-rpc/server"
	"log"
	"net"
	"time"
)

func startServer(addr chan string) {
	num := Num{}
	server.Register(&num)

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}


	addr <- l.Addr().String()
	guy_rpc.Accept(l)
}

type Arg struct {
	A int
	B int
}

func main() {

	addr := make(chan string)
	go startServer(addr)
	c, _ := guy_rpc.Dial("tcp", <-addr, guy_rpc.DefaultOption)

	defer func() { _ = c.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	//var wg sync.WaitGroup
	//for i := 0; i < 5; i++ {
	//	wg.Add(1)
	//	go func(i int) {
	//		defer wg.Done()

	var reply int
	if err := c.ASyncCall("Add", &reply,nil, Arg{
		A: 1,
		B: 3,
	}); err != nil {
		log.Println(err)
	}
	log.Println("reply:", reply)
	var back string
	if err := c.SyncCall("Name", &back,"wzx"); err != nil {
		log.Println(err)
	}
	log.Println("reply:", back)
	//	}(i)
	//}
	//wg.Wait()
}

type Num struct{}

func (num *Num) Add(arg Arg) int {
	return arg.A + arg.B
}

func (num *Num) Name(hi string) string {
	return "hi" + hi
}
