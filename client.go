package guy_rpc

import (
	"guy-rpc/client"
	"guy-rpc/server"
)

func Dial(network, address string, option *server.Option) (c *client.Client, err error)  {
	return client.Dial(network,address,option)
}