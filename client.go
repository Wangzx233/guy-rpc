package guy_rpc

import (
	"guy-rpc/client"
	"guy-rpc/server"
)

func Dial(network, address string, option *server.Option, centerAdr ...string) (c *client.Client, err error) {
	if centerAdr != nil {
		return client.Dial(network, address, option, centerAdr[0])
	} else {
		return client.Dial(network, address, option)
	}

}
