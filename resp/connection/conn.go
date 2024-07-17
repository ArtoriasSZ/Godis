package connection

import (
	"Godis/lib/sync/wait"
	"net"
	"sync"
	"time"
)

type ClientConnection struct {
	conn     net.Conn
	waiting  wait.Wait
	mu       sync.Mutex
	selectDB int
}

func NewClientConnection(conn net.Conn) *ClientConnection {
	return &ClientConnection{
		conn: conn,
	}
}

func (c *ClientConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *ClientConnection) Close() {
	c.waiting.WaitWithTimeout(10 * time.Second)
	_ = c.conn.Close()
}

func (c *ClientConnection) Writer(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}
	c.mu.Lock()
	// 防止在写的过程中被关闭
	c.waiting.Add(1)
	defer func() {
		c.waiting.Done()
		c.mu.Unlock()
	}()
	_, err := c.conn.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConnection) GetDBIndex() int {
	return c.selectDB
}

func (c *ClientConnection) SelectDB(dbNum int) error {
	c.selectDB = dbNum
	return nil
}
