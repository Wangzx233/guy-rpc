package guy_rpc

import (
	"guy-rpc/server"
	"net"
)

// Accept 对外封装一下server的accept
func Accept(lis net.Listener)  {
	server.DefaultServer.Accept(lis)
}

// DefaultOption 对外封装默认预检
var DefaultOption = &server.Option{
	IdentityCode: server.IdentityCode,
	CodecType:    "application/json",
}
