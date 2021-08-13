package center

import (
	"encoding/json"
	"fmt"
	"guy-rpc/codec"
	"guy-rpc/server"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

var methods = make(map[string]string)
var mutex = new(sync.Mutex)
var fromClient = make(map[uint64]string)

func receiveClient(listener net.Listener) {
	for {

		//接收部分
		conn, err := listener.Accept()
		if err != nil {
			log.Println("registerCenter : accept error:", err)
			return
		}

		defer func() { _ = conn.Close() }()

		//接收传来的信息的option并反序列化后存入option
		var option = server.Option{}
		err = json.NewDecoder(conn).Decode(&option)
		if err != nil {
			log.Println("rpc server: options error: ", err)
		}

		//判断是否是guy-rpc的消息
		if option.IdentityCode != server.IdentityCode {
			log.Println("非guy-rpc消息")
		}
		c := codec.NewJsonCodec(conn)
		fmt.Println(option)

		for {
			var h codec.Header
			err = c.ReadHeader(&h)
			if err != nil {

				if err != io.EOF && err != io.ErrUnexpectedEOF {
					log.Println("registerCenter : read header error", err)
				}
				break
			}

			////判断方法是否注册
			//adr, ok := methods[h.Method]
			//if !ok {
			//	log.Println("方法未注册")
			//	break
			//}
			//记录客户端地址
			fromClient[h.CallID] = conn.RemoteAddr().String()
			req := &request{h: &h}

			req.argv, err = server.FindMethodIn(h.Method)
			if err != nil {
				log.Println("rpc center: find method err:", err)
			}

			argvi := req.argv.Interface()
			if req.argv.Type().Kind() != reflect.Ptr {
				argvi = req.argv.Addr().Interface()
			}

			err = c.ReadBody(argvi)

			if err != nil {
				log.Println("rpc center: read body err:", err)
			}

			//con, err := net.Dial("tcp", ":10011")
			//defer con.Close()
			//if err != nil {
			//	log.Println("rpc center: dial err:", err)
			//}
			con, err := net.Dial("tcp", ":10010")
			if err != nil {
				log.Println(err)
			}
			err = json.NewEncoder(con).Encode(option)
			if err != nil {
				log.Println("rpc client: option error: ", err)
				_ = con.Close()
			}
			co := codec.NewJsonCodec(con)
			err = co.Write(&h, argvi)
			if err != nil {
				log.Println("rpc center: write err:", err)
			}

			receiveServer(con)

		}

		_ = c.Close()
	}
}

func receiveServer(conn net.Conn) {

	//接收部分
	//conn, err := listener.Accept()
	//if err != nil {
	//	log.Println("registerCenter : accept error:", err)
	//	return
	//}

	c := codec.NewJsonCodec(conn)

	var h codec.Header
	err := c.ReadHeader(&h)

	if err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("registerCenter : read header error", err)
		}

	}

	req := &request{h: &h}

	req.argv, err = server.FindMethodIn(h.Method)
	if err != nil {
		log.Println("rpc center: find method err:", err)
	}

	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}
	err = c.ReadBody(argvi)
	if err != nil {
		log.Println("rpc center: read body err:", err)
	}

	conn, err = net.Dial("tcp", fromClient[h.CallID])
	c = codec.NewJsonCodec(conn)

	err = c.Write(&h, argvi)
	if err != nil {
		log.Println("rpc center: write err:", err)
	}


	_ = c.Close()

}
