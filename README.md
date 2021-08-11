# guy-rpc

## 快速开始

```
go get -u github.com/Wangzx233/guy-rpc
```

#### 服务端

```go
	num := Num{}
	
	//注册Num{}结构体下的所有方法
	guy_rpc.Register(&num)
	
	lis, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal("network error:", err)
	}
	
	guy_rpc.Accept(lis)
```

#### 客户端

```go
c, _ := guy_rpc.Dial("tcp", ":9999", guy_rpc.DefaultOption)
defer func() { _ = c.Close() }()


var reply int

	//同步调用
	if err := c.SyncCall("Add", &reply, Arg{
		A: 1,
		B: 3,
	}); err != nil {
		log.Println(err)
	}
	log.Println("reply:", reply)

	//异步调用
	if err = c.ASyncCall("Add", &reply, nil,Arg{
		A: 1,
		B: 3,
	}); err != nil {
		log.Println(err)
	}
	log.Println("reply:", reply)
```

