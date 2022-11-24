package frame

import (
	"encoding/binary"
	"errors"
	"io"
)

/*
Frame定义

frameHeader + framePayload(packet)

frameHeader
	4 bytes: length 整型，帧总长度(含头及payload)

framePayload
	Packet
*/

// Payload 定义载荷数据类型
type Payload []byte

// StreamFrameCodec 接口
type StreamFrameCodec interface {
	Encode(io.Writer, Payload) error   // data -> frame，并写入io.Writer
	Decode(io.Reader) (Payload, error) // 从io.Reader中提取frame payload，并返回给上层
}

var ErrShortWrite = errors.New("short write")
var ErrShortRead = errors.New("short read")

type Codec struct {
}

// NewCodec 创建 Frame 编码解码器
func NewCodec() StreamFrameCodec {
	return &Codec{}
}

// Encode Frame 层的编码
// 进入这里之前 w 应该已经有 frameHeader 了 以也就是 已经写入了 4 位的 totalLen
func (c *Codec) Encode(w io.Writer, framePayload Payload) error {
	// 复制出来一份
	var f = framePayload

	totalLen := int32(len(framePayload)) + 4
	// 把 totalLen 写入了 w 中 0x0 0x0 0x0 0x9 -> BigEndian | littleEndian -> 0x9 0x0 0x0 0x0
	// 把 int32的数字写入 w
	err := binary.Write(w, binary.BigEndian, &totalLen) // write the frame payload to outbound stream
	if err != nil {
		return err
	}
	// 然后在写数据
	n, err := w.Write(f)
	if n != len(framePayload) {
		return ErrShortWrite
	}
	return nil
}

// Decode Frame 层的解码
// 进入这里之前 r 应该已经有 frameHeader 了 以也就是 已经读取了 4 位的 totalLen
func (c *Codec) Decode(r io.Reader) (Payload, error) {
	var totalLen int32
	err := binary.Read(r, binary.BigEndian, &totalLen)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, totalLen-4)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	if n != int(totalLen-4) {
		return nil, ErrShortRead
	}
	return buf, nil
}
