package client

import (
	"encoding/json"
	"fmt"
	"guy-rpc/codec"
	err "guy-rpc/error"
	"guy-rpc/server"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
)

type Call struct {
	CallID   uint64      //识别id
	Method   string      //方法名
	Args     interface{} //函数参数
	CallBack interface{} //返回值
	Error    error       //错误返回
	Done     chan *Call  //用于同步的确认管道
}

//调用结束时，call.done告诉调用方
func (call *Call) done() {
	call.Done <- call
}

type Client struct {
	CallID    uint64         //请求编号
	c         codec.Codec    //序列化和反序列化消息
	option    *server.Option //预检
	mutex     sync.Mutex
	sendMutex sync.Mutex
	header    codec.Header
	pending   map[uint64]*Call
	isClose   bool //用户是否关闭
	stop      bool //服务器是否停止
}

// Close 关闭连接
func (client *Client) Close() error {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if client.isClose {
		return err.Closed
	}
	client.isClose = true
	return client.c.Close()
}

// IsAvailable 简单封装判断连接是否可用
func (client *Client) IsAvailable() bool {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	return !client.isClose && !client.stop
}

func Dial(network, address string, option *server.Option) (client *Client, err error) {

	if address == ":8000" {

		request, err := http.NewRequest("GET", "http://127.0.0.1:8000/_guy-rpc_/register", nil)
		if err != nil {
			fmt.Println(err)
		}
		request.Header.Set("Content-type", "application/json")
		request.Header.Set("X-GuyRpc-Servers","")
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			goto there
		}
		var get string

		get = response.Header.Get("X-GuyRpc-Servers")

		split := strings.Split(get, ",")

		address = split[0]

	}
	there:
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr(network, address)

	conn, err := net.DialTCP(network, nil, tcpAddr)

	if err != nil {
		return nil, err
	}

	defer func() {
		if client == nil {
			_ = conn.Close()
		}
	}()

	return NewClient(conn, option)
}

// NewClient 构建一个新的Client实例
func NewClient(conn net.Conn, option *server.Option) (*Client, error) {
	err := json.NewEncoder(conn).Encode(option)
	if err != nil {
		log.Println("rpc client: option error: ", err)
		_ = conn.Close()
		return nil, err
	}

	jsonCodec := codec.NewJsonCodec(conn)

	client := &Client{
		CallID:  1,
		c:       jsonCodec,
		option:  option,
		pending: make(map[uint64]*Call),
	}

	go client.receive()
	return client, nil
}
