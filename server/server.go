package server

import (
	"log"
	"net"
)

// Server 是RPC的服务
type Server struct{}

// NewServer 构造一个Server
func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

// Accept 接收所有传过来的请求并进行处理
func (server *Server) Accept(listener net.Listener) {
	for {
		//接收部分
		conn, err := listener.Accept()
		if err != nil {
			log.Println("rpc server: accept error:",err)
			return
		}

		//进行处理
		server.ServeConn(conn)
	}
}


