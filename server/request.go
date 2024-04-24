package server

import (
	"errors"
	"guy-rpc/codec"
	"io"
	"log"
	"reflect"
	"sync"
)

type request struct {
	h    *codec.Header
	argv reflect.Value
	body reflect.Value
}

// 读取请求中的头部
func (server *Server) readRequestHeader(c codec.Codec) (*codec.Header, error) {
	var h codec.Header
	err := c.ReadHeader(&h)
	if err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server : read header error", err)
		}
		return nil, err
	}
	return &h, nil
}

// 读取请求
func (server *Server) readRequest(c codec.Codec) (*request, error) {
	//读取头部
	header, err := server.readRequestHeader(c)
	if err != nil {
		return nil, err
	}
	req := &request{h: header}

	req.argv, err = FindMethodIn(header.Method)
	if err != nil {
		log.Println("rpc server: find method err:", err)
	}

	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}

	//读取body
	err = c.ReadBody(argvi)
	if err != nil {
		log.Println("rpc server: read argv err:", err)
	}

	return req, nil
}

func (server *Server) sendResponse(c codec.Codec, header *codec.Header, body interface{}, mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()
	err := c.Write(header, body)
	if err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

func (server *Server) handleRequest(c codec.Codec, req *request, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	req.body = FindMethodOut(req.h.Method, req.argv)

	server.sendResponse(c, req.h, req.body.Interface(), mutex)
}

func FindMethodIn(method string) (argv reflect.Value, err error) {
	handler, ok := Handlers[method]
	if !ok {
		return reflect.ValueOf("服务未注册"), errors.New("服务未注册")
	}
	if handler.In.Kind() == reflect.Ptr {
		argv = reflect.New(handler.In.Elem())
	} else {
		argv = reflect.New(handler.In).Elem()
	}

	return argv, nil
}

func FindMethodOut(method string, arg reflect.Value) reflect.Value {
	handler := Handlers[method]
	call := handler.method.Func.Call([]reflect.Value{handler.self, arg})
	return call[0]
}
