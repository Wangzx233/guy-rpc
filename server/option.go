package server

const IdentityCode = 0x3bef5c

// Option 预检
type Option struct {
	IdentityCode int    //识别是否是guy-rpc的请求
	CodecType    string //序列化和反序列化的方式（目前只有JSON）
}


