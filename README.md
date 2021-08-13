# guy-rpc

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
    //如果需要使用服务中心，地址应为服务中心地址
    c, _ := guy_rpc.Dial("tcp", ":8000", 	guy_rpc.DefaultOption)
    
	defer func() { _ = c.Close() }()

	var wg sync.WaitGroup
	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			var reply int
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
	register.HandleHTTP()
    
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Println(err)
	}
}
```

