package main

import (
	"fmt"
	guy_rpc "guy-rpc"
	"log"
	"net"
)

func startServer() {
	num := Num{}
	guy_rpc.Heartbeat("127.0.0.1:9999", "math", "127.0.0.1:8080")

	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal("network error:", err)
	}

	guy_rpc.Register(&num, l.Addr().String())

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
	startServer()
	fmt.Println("server started")
	//time.Sleep(time.Second*2)
	//c, _ := guy_rpc.Dial("tcp", "", guy_rpc.DefaultOption, ":8000")
	//fmt.Println("client started")
	//defer func() { _ = c.Close() }()
	//
	//var wg sync.WaitGroup
	//for i := 0; i < 5; i++ {
	//	wg.Add(1)
	//	go func(i int) {
	//		defer wg.Done()
	//
	//		var reply int
	//		if err := c.SyncCall("Add", &reply, Arg{
	//			A: i,
	//			B: i + 5,
	//		}); err != nil {
	//			log.Println(err)
	//		}
	//		fmt.Println("reply:", i, "+", i+5, "=", reply)
	//
	//	}(i)
	//}
	//wg.Wait()
}

type Num struct{}

func (num *Num) Add(arg Arg) int {
	return arg.A + arg.B
}

func (num *Num) Name(struct{}) string {
	return "hello"
}
