package guy_rpc

import (
	"guy-rpc/server"
	"net"
)

// Accept 对外封装一下server的accept
func Accept(l net.Listener) {
	server.DefaultServer.Accept(l)
}

// DefaultOption 对外封装默认预检
var DefaultOption = &server.Option{
	IdentityCode: server.IdentityCode,
	CodecType:    "application/json",
}

func Register(str interface{}, localAdr string, centerAdr ...string) {

	if centerAdr != nil {
		server.Register(str, localAdr, centerAdr[0])
	} else {
		server.Register(str, localAdr)
	}

}
