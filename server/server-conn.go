package server

import (
	"encoding/json"
	"guy-rpc/codec"
	"io"
	"log"
	"sync"
)

// ServeConn 处理接收到的消息
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()
	var option Option

	//接收传来的信息的option并反序列化后存入option
	err := json.NewDecoder(conn).Decode(&option)
	if err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}
	//判断是否是guy-rpc的消息
	if option.IdentityCode != IdentityCode {
		log.Println("非guy-rpc消息")
		return
	}


	//接收并处理conn中剩余的信息（header和body）
	server.serveCodec(codec.NewJsonCodec(conn))

}

func (server *Server) serveCodec(c codec.Codec) {

	mutex := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := server.readRequest(c)
		if err != nil {
			if req == nil {

				break //无请求就关闭连接
			}
			req.h.Error = err.Error()
			server.sendResponse(c, req.h, struct{}{}, mutex)
			continue
		}
		wg.Add(1)
		go server.handleRequest(c, req, mutex, wg)

	}
	wg.Wait()
	_ = c.Close()
}
