package server

import (
	"net/http"
	"reflect"
	"time"
)


type Handler struct {
	self   reflect.Value
	Func   reflect.Value
	In     reflect.Type
	method reflect.Method
	NumIn  int
	NumOut int
}

var Handlers = make(map[string]*Handler)

func Register(str interface{},localAdr string,centerAdr ...string) {

	v := reflect.ValueOf(str)
	t := reflect.TypeOf(str)
	for i := 0; i < v.NumMethod(); i++ {
		name := t.Method(i).Name
		method := v.Method(i)
		m := t.Method(i)
		//Handlers[name].In = new([]reflect.Type)
		//Handlers[name].Out = make([]reflect.Type, method.Type().NumOut())
		//
		//for j := 0; i < method.Type().NumIn(); j++ {
		//	//Handlers[name].In = append(Handlers[name].In, method.Type().In(i))
		//	Handlers[name].In[i]=method.Type().In(i)
		//}
		//
		//for j := 0; i < method.Type().NumOut(); j++ {
		//	Handlers[name].Out[i] =method.Type().Out(i)
		//}

		args := make([]string, 0, m.Type.NumIn())
		returns := make([]string, 0, m.Type.NumOut())
		// j 从 1 开始，第 0 个入参是 wg 自己。
		for j := 1; j < m.Type.NumIn(); j++ {
			args = append(args, m.Type.In(j).Name())
		}
		for j := 0; j < m.Type.NumOut(); j++ {
			returns = append(returns, m.Type.Out(j).Name())
		}

		Handlers[name] = &Handler{
			self:   v,
			Func:   method,
			method: m,
			In:     m.Type.In(1),
			NumIn:  method.Type().NumIn(),
			NumOut: method.Type().NumOut(),
		}

		if centerAdr!=nil {
			registerCenter(localAdr,centerAdr[0])
		}
	}

}

func registerCenter(localAdr string,centerAdr string)  {

		request, err := http.NewRequest("POST", "http://"+centerAdr+"/_guy-rpc_/register", nil)
		if err != nil {
			return
		}
		request.Header.Set("Content-type", "application/json")

		request.Header.Set("X-GuyRpc-Servers",localAdr)
		client := &http.Client{}
		client.Timeout = time.Second

		_, err = client.Do(request)
		if err != nil {
			client.CloseIdleConnections()
			return
		}


	//conn, err := net.Dial("tcp", ":10086")
	//if err != nil {
	//	//说明服务中心未启动，直接return
	//	return
	//}
	//
	//c := codec.NewJsonCodec(conn)
	//
	//var h = codec.Header{Method: method}
	//err = c.Write(&h, adr)
	//if err != nil {
	//	log.Println("server register err: write err:",err)
	//}

}
