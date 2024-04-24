package client

import (
	"log"
)

func (client *Client) send(call *Call) {
	client.sendMutex.Lock()
	defer client.sendMutex.Unlock()

	callID, err := client.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}

	client.header.Method = call.Method
	client.header.CallID = callID
	client.header.Error = ""

	err = client.c.Write(&client.header, call.Args)

	if err != nil {
		client.removeCall(callID)
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

func (client *Client) ASyncCall(Method string, callBack interface{}, done chan *Call, args interface{}) *Call {
	if done == nil {
		done = make(chan *Call, 1)
	} else if cap(done) == 0 {
		log.Panic("rpc client: done channel is unbuffered")
	}

	call := &Call{
		Method:   Method,
		Args:     args,
		CallBack: callBack,
		Done:     done,
	}
	client.send(call)
	return call
}

func (client *Client) SyncCall(Method string, callBack interface{}, args interface{}) error {
	call := <-client.ASyncCall(Method, callBack, make(chan *Call, 1), args).Done
	return call.Error
}
