package server

import (
	"reflect"
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

func Register(str interface{}) {

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
	}
}
