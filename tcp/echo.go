package tcp

import (
	"Godis/lib/logger"
	"Godis/lib/sync/atomic"
	"Godis/lib/sync/wait"
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"
)

type EchoClient struct {
	Coon    net.Conn
	Waiting wait.Wait
}

func (echoClient *EchoClient) Close() error {
	timeout := echoClient.Waiting.WaitWithTimeout(10 * time.Second)
	if timeout == true {
		logger.Warn("echo client close timeout")
		return echoClient.Coon.Close()
	} else {
		return echoClient.Coon.Close()
	}

}

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (handler *EchoHandler) Handler(ctx context.Context, conn net.Conn) {
	if handler.closing.Get() {
		_ = conn.Close()
	}
	echoClient := &EchoClient{
		Coon: conn,
	}
	handler.activeConn.Store(echoClient, struct{}{})
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("read EOF, Connction closed")
				handler.activeConn.Delete(echoClient)
			} else {
				logger.Warn("read error: ", err)
			}
			return
		}
		echoClient.Waiting.Add(1)
		b := []byte(msg)
		_, _ = conn.Write(b)
		echoClient.Waiting.Done()
	}
}

func (handler *EchoHandler) Close() error {
	logger.Info("echo handler closed")
	handler.closing.Set(true)
	handler.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})
	return nil
}
