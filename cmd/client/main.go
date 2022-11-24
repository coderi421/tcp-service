package main

import (
	"fmt"
	"github.com/CoderI421/tcp-service/frame"
	"github.com/CoderI421/tcp-service/packet"
	"github.com/lucasepe/codename"
	"net"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var num = 5

	wg.Add(num)

	// 5 个协成并发客户端
	for i := 0; i < num; i++ {
		go func(i int) {
			defer wg.Done()
			startClient(i)
		}(i + 1)
	}
	wg.Wait()
}

func startClient(clientId int) {
	// 客户端关闭通知
	quit := make(chan struct{})
	// 接收服务端响应完毕的信号通知
	done := make(chan struct{})
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}

	defer conn.Close()
	fmt.Printf("[client %d]: dial ok", clientId)

	// 生成随机的 payload
	rng, err := codename.DefaultRNG()
	if err != nil {
		panic(err)
	}

	frameCodec := frame.NewCodec()
	var counter int

	// 处理 server 的响应信息
	go func() {
		for {
			select {
			// 如果客户端关闭，这里要处理完信息，再关闭
			case <-quit:
				// 处理完 server 响应信息的信号
				done <- struct{}{}
				return
			// 没有 default 则无法关闭
			default:
			}

			conn.SetReadDeadline(time.Now().Add(time.Second * 5))
			ackFramePayLoad, err := frameCodec.Decode(conn)
			if err != nil {
				// 判断是否是超时错误
				if e, ok := err.(net.Error); ok {
					if e.Timeout() {
						continue
					}
				}
				panic(err)
			}

			p, err := packet.Decode(ackFramePayLoad)
			if err != nil {
				fmt.Println("packet decode error", err)
			}
			switch ack := p.(type) {
			case *packet.SubmitAck:
				fmt.Printf("[client %d]: the result of submit ack[%s] is %d\n", clientId, ack.ID, ack.Result)
			case *packet.ConAck:
				fmt.Printf("[client %d]: the result of submit ack[%s] is %d\n", clientId, ack.ID, ack.Result)
			default:
				panic("not submitAck or connAck")
			}
		}
	}()

	// 处理 client 向 server 发送信息
	for {
		// send submit
		counter++
		id := fmt.Sprintf("%08d", counter)
		payload := codename.Generate(rng, 4)
		// todo 这里应该 先获取 connAck 后才进行
		// 构建 submit packet 实例
		s := &packet.Submit{
			ID:      id,
			Payload: []byte(payload),
		}

		framePayload, err := packet.Encode(s)
		if err != nil {
			panic(err)
		}

		fmt.Printf("[client %d]: send submit id = %s, payload=%s, frame length = %d\n",
			clientId, s.ID, s.Payload, len(framePayload)+4)

		err = frameCodec.Encode(conn, framePayload)
		if err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)

		if counter >= 10 {
			// 关闭客户端
			quit <- struct{}{}
			// 等待处理 server ack 结束的信号
			<-done
			fmt.Printf("[client %d]: exit ok\n", clientId)
			return
		}
	}
}
