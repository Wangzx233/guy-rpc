# guy-rpc

#### 功能

- [x] client端 远程调用
- [x] server端 注册服务
- [x] 服务中心
- [x] 心跳包
- [x] 支持并发调用

## 快速开始

```
go get -u github.com/Wangzx233/guy-rpc
```

#### 服务端

```go
func main(){
    num := Num{}
	
	//注册Num{}结构体下的所有方法
	guy_rpc.Register(&num)
	
	//地址随意
	lis, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal("network error:", err)
	}
    
    //发送心跳包（可选），只能在有服务中心时使用。"/_guy-rpc_/register"为固定值
	register.Heartbeat("http://127.0.0.1:8000/_guy-rpc_/register", "tcp@"+l.Addr().String(), 0)
    
    //开始监听
	guy_rpc.Accept(lis)
}
	


type Num struct{}
func (num *Num) Add(arg Arg) int {
	return arg.A + arg.B
}
```

#### 客户端

```go

func main(){
    //第二个参数为服务端地址，第四个参数为服务中心地址。根据需要p2p模式或cs模式选填一个即可，无服务中心时第四个参数不能填
    c, _ := guy_rpc.Dial("tcp", "", 	guy_rpc.DefaultOption,":8000")
    
	defer func() { _ = c.Close() }()

	var wg sync.WaitGroup
	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			var reply int
            //如果入参为多个，请打包为结构体，出参同理。无出参或者入参用空结构体占位即可。
			if err := c.SyncCall("Add", &reply, Arg{
				A: i,
				B: i+5,
			}); err != nil {
				log.Println(err)
			}
            fmt.Println("reply:",i,"+",i+5,"=", reply)
            
		}(i)
	}
	wg.Wait()
}

```

#### 服务中心（可选）

```go
func main()  {
	guy_rpc.HandleHTTP()
    
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Println(err)
	}
}
```

