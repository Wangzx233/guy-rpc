package codec

import (
	"io"
)

type Header struct {
	Method string //请求的服务名和方法名
	CallID uint64 //请求的id
	Error  string //如果服务端发生错误，错误写入error
}

// Codec 的接口
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(value interface{}) error
	Write(*Header, interface{}) error
}
