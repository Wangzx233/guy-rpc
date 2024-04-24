package server

const IdentityCode = 0x3bef5c

// Option 预检（基本上没用）
type Option struct {
	IdentityCode int //识别是否是guy-rpc的请求
	CodecType    string
}
