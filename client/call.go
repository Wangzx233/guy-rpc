package client

import err "guy-rpc/error"

func (client *Client) registerCall(call *Call) (uint64, error) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if client.isClose || client.stop {
		return 0, err.Closed
	}
	call.CallID = client.CallID
	client.pending[call.CallID] = call
	client.CallID++
	return call.CallID, nil
}

//删除等待中的call
func (client *Client) removeCall(callID uint64) *Call {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	call := client.pending[callID]
	delete(client.pending, callID)
	return call
}

// StopCalls 服务端或者客户端发生错误时调用，停止所有的call
func (client *Client) StopCalls(err error) {
	client.sendMutex.Lock()
	defer client.sendMutex.Unlock()
	for _, call := range client.pending {
		call.Error=err
		call.done()
	}
}

