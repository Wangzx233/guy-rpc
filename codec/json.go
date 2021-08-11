package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type JsonCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

func (j *JsonCodec) Close() error {
	return j.conn.Close()
}

func (j *JsonCodec) ReadHeader(h *Header) error {
	return j.dec.Decode(h)
}

func (j *JsonCodec) ReadBody(i interface{}) error {
	return j.dec.Decode(i)
}

func (j *JsonCodec) Write(header *Header, i interface{}) error {
	defer func() {
		//最后将缓存写入
		err := j.buf.Flush()
		if err != nil {
			//出错则关闭io
			_ = j.Close()

		}
	}()
	//序列化Header
	err := j.enc.Encode(header)
	if err != nil {
		log.Println("rpc codec: json error encoding header:", err)
		return err
	}

	//序列化Body

	err = j.enc.Encode(i)
	if err != nil {
		log.Println("rpc codec: json error encoding body:", err)
	}
	return err
}

func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &JsonCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}
