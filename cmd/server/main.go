package main

import (
	"fmt"
	"github.com/CoderI421/tcp-service/frame"
	"github.com/CoderI421/tcp-service/packet"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error: ", err)
		return
	}

	fmt.Println("server listening on(*:8888)")

	for {
		accept, aErr := listen.Accept()
		if aErr != nil {
			fmt.Println("accept error:", aErr)
			break
		}

		// start a new goroutine to handle
		// the new connection.
		// 每个客户端连接，由一个协成进行处理
		go handleConn(accept)
	}
}

// handleConn 第一层，解析 Frame 层
func handleConn(c net.Conn) {
	// c -> 和每个客户端的连接
	defer c.Close()
	frameCodec := frame.NewCodec()

	for {
		// read from the connection

		// decode the frame to get the payload
		// is undecoded packet
		framePayload, err := frameCodec.Decode(c)
		if err != nil {
			fmt.Println("handleConn: frame decode error:", err)
			return
		}

		// do something with the packet
		// packet层的响应
		ackFramePayload, err := handlePacket(framePayload)
		if err != nil {
			fmt.Println("handleConn: frame decode error:", err)
			return
		}

		// write ack frame to the connection
		// Frame 层编码
		err = frameCodec.Encode(c, ackFramePayload)
		if err != nil {
			fmt.Println("handleConn: frame encode error:", err)
			return
		}
	}
}

// handlePacket 第二层，解析 packet 层
func handlePacket(framePayload []byte) (ackFramePayload []byte, err error) {
	var p packet.Packet
	// 解析后，获取 packet 实例 或是 submit submitAck conn connAck
	p, err = packet.Decode(framePayload)
	if err != nil {
		fmt.Println("handleConn: packet decode error:", err)
		return
	}

	switch p.(type) {
	case *packet.Submit:
		// 获取请求信息
		submit := p.(*packet.Submit)
		fmt.Printf("recv submit: id = %s, payload=%s\n", submit.ID, string(submit.Payload))
		// 根据请求信息，响应信息
		submitAck := &packet.SubmitAck{
			ID:     submit.ID,
			Result: 0,
		}
		ackFramePayload, err = packet.Encode(submitAck)
		if err != nil {
			fmt.Println("handleConn: packet decode error:", err)
			return nil, err
		}
		return ackFramePayload, nil
	case *packet.Con:
		// 获取请求信息
		conn := p.(*packet.Con)
		fmt.Printf("recv conn: id = %s, payload=%s\n", conn.ID, string(conn.Payload))
		connAck := &packet.Con{
			ID:      conn.ID,
			Payload: nil,
		}
		ackFramePayload, err = packet.Encode(connAck)
		if err != nil {
			fmt.Println("handleConnAck:packet decode error:", err)
			return nil, err
		}
		return ackFramePayload, nil
	default:
		return nil, fmt.Errorf("unknown packet type")
	}
}
