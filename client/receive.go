package client

import (
	"errors"
	"fmt"
	"guy-rpc/codec"
)

//接收返回的消息
func (client *Client) receive() {
	var err error
	for err == nil {
		var h codec.Header
		err = client.c.ReadHeader(&h)
		if err != nil {
			break
		}
		//每处理一个call，就从待处理map中删除一个
		call := client.removeCall(h.CallID)

		switch {
		case call == nil:
			err = client.c.ReadBody("")
		case h.Error != "":
			call.Error = fmt.Errorf(h.Method)
			err = client.c.ReadBody(nil)
			call.done()
		default:
			err = client.c.ReadBody(call.CallBack)
			if err != nil {
				call.Error = errors.New("read body error:" + err.Error())

			}
			call.done()
		}
	}

	client.StopCalls(err)
}
