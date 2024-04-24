package center

import (
	"guy-rpc/codec"
	"log"
	"net"
	"reflect"
	"sync"
)

type request struct {
	h    *codec.Header
	argv reflect.Value
	body reflect.Value
}

var Wg = new(sync.WaitGroup)

func StartCenter() {

	registerAdr, err := net.Listen("tcp", ":10086")
	if err != nil {
		log.Println(err)
	}

	Wg.Add(1)
	go registerCenter(registerAdr)

	client, err := net.Listen("tcp", ":10010")
	if err != nil {
		log.Println(err)
	}

	Wg.Add(1)
	go receiveClient(client)

	Wg.Wait()
}
