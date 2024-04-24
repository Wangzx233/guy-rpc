package center

import (
	"guy-rpc/codec"
	"io"
	"log"
	"net"
)

func registerCenter(listener net.Listener) {
	for {

		//接收部分
		conn, err := listener.Accept()
		if err != nil {
			log.Println("registerCenter : accept error:", err)
			return
		}

		defer func() { _ = conn.Close() }()

		c := codec.NewJsonCodec(conn)

		var h codec.Header
		err = c.ReadHeader(&h)
		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				log.Println("registerCenter : read header error", err)
			}
		}

		var adr string
		err = c.ReadBody(&adr)
		if err != nil {
			log.Println("registerCenter :read body err:", err)
		}

		methods[h.Method] = ":10010"

	}
	Wg.Done()
}

func StringsContains(array []string, val string) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}
