package packet

import (
	"bytes"
	"fmt"
	"sync"
)

// Packet协议定义

/*
### packet header
1 byte: commandID 类型

### submit packet

8字节 ID 字符串
任意字节 payload

### submit ack packet

8字节 ID 字符串
1字节 result
*/

const (
	CommandConn   = iota + 0x01 // 连接请求包（值为0x01）
	CommandSubmit               // 消息请求包（值为0x02）
)

const (
	CommandConnAck   = iota + 0x80 // 连接响应包（值为0x81）
	CommandSubmitAck               // 消息响应包（值为0x82）
)

const (
	OkResponse = "OK"
)

type Packet interface {
	Decode([]byte) error     // []byte -> struct
	Encode() ([]byte, error) //  struct -> []byte
}

type Submit struct {
	ID      string // ID 消息请求包的ID
	Payload []byte // Payload 消息请求包的具体信息
}

// Decode 解析 Packet 中的信息
func (s *Submit) Decode(packetBody []byte) error {
	s.ID = string(packetBody[:8]) // 取前 8 个字符 转换成字符串
	s.Payload = packetBody[8:]    // 取剩下所有的 具体内容
	return nil
}

// Encode 编译 Packet 中的信息
func (s *Submit) Encode() ([]byte, error) {
	// return []byte(s.ID + string(s.Payload)), nil
	// 这个地方需要补齐8位s.ID[:8]
	return bytes.Join([][]byte{[]byte(s.ID[:8]), s.Payload}, nil), nil
}

type SubmitAck struct {
	ID     string // 消息响应包Id
	Result uint8  // 结果 ack 的result 是 0/1
}

func (s *SubmitAck) Decode(packetBody []byte) error {
	s.ID = string(packetBody[:8]) // 取得ID
	s.Result = packetBody[8]      // 取得结果
	return nil
}

func (s *SubmitAck) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.ID[:8]), []byte{s.Result}}, nil), nil
}

// Con 连接请求包
type Con struct {
	ID      string
	Payload []byte
}

func (c *Con) Decode(connBody []byte) error {
	c.ID = string(connBody[:8])
	c.Payload = connBody[8:]
	return nil
}

func (c *Con) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(c.ID[:8]), c.Payload}, nil), nil
}

// ConAck 连接请求包
type ConAck struct {
	ID     string // 连接响应包Id
	Result uint8  // 结果 ack 的result 是 0/1
}

func (c *ConAck) Decode(connBody []byte) error {
	c.ID = string(connBody[:8]) // 取得id
	c.Result = connBody[8]
	return nil
}

func (c *ConAck) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(c.ID[:8]), []byte{c.Result}}, nil), nil
}

var SubmitPool = sync.Pool{
	New: func() interface{} {
		return &Submit{}
	},
}

// Decode 根据 frame 的解析结果，继续解析
func Decode(packet []byte) (Packet, error) {
	commandID := packet[0] // 1 byte: commandID 类型
	pktBody := packet[1:]

	switch commandID {
	case CommandConn:
		c := &Con{}
		err := c.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return c, nil
	case CommandConnAck:
		c := &ConAck{}
		err := c.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return c, nil
	case CommandSubmit:
		// s := Submit{}
		// 优化堆内存 使用池化 Submit
		s := SubmitPool.Get().(*Submit) // get submit pool
		err := s.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return s, nil
	case CommandSubmitAck:
		s := &SubmitAck{}
		err := s.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return s, nil
	default:
		return nil, fmt.Errorf("unknown commandID [%d]", commandID)
	}
}

// Encode 编译 Packet 将结果向上传递给 Frame
func Encode(p Packet) ([]byte, error) {
	var (
		commandID uint8
		pktBody   []byte
		err       error
	)

	switch t := p.(type) {
	case *Con:
		commandID = CommandConn
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *ConAck:
		commandID = CommandConnAck
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *Submit:
		commandID = CommandSubmit
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *SubmitAck:
		commandID = CommandSubmitAck
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown type [%s]", t)
	}
	return bytes.Join([][]byte{{commandID}, pktBody}, nil), nil
}
