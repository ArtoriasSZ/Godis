package handler

import (
	"Godis/database"
	databaseface "Godis/interface/database"
	"Godis/lib/logger"
	"Godis/lib/sync/atomic"
	"Godis/resp/connection"
	"Godis/resp/parser"
	"Godis/resp/reply"
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"sync"
)

type RespHandler struct {
	activeConn sync.Map
	db         databaseface.Database
	closing    atomic.Boolean
}

func NewRespHandler() *RespHandler {
	var db = database.NewEchodatabase()
	//todo: 实现Database
	return &RespHandler{
		db: db,
	}

}

func (r *RespHandler) closeClientConnection(client *connection.ClientConnection) {
	client.Close()
	r.db.AfterClientClose()
	r.activeConn.Delete(client)
}

func (r *RespHandler) Handler(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		_ = conn.Close()
	}
	client := connection.NewClientConnection(conn)
	r.activeConn.Store(client, struct{}{})
	ch := parser.ParseStream(conn)
	for payLoad := range ch {
		// error
		if payLoad.Err != nil {
			if payLoad.Err == io.EOF ||
				errors.Is(payLoad.Err, io.ErrUnexpectedEOF) ||
				strings.Contains(payLoad.Err.Error(), "use of closed network connection") {
				r.closeClientConnection(client)
				logger.Info("1client closed the connection:" + client.RemoteAddr().String())
				return
			}
			// 协议错误
			errReply := reply.NewErrReply(payLoad.Err.Error())
			errReplyBytes, _ := errReply.ToBytes()
			err := client.Writer(errReplyBytes)
			if err != nil {
				logger.Info("2client closed the connection:" + client.RemoteAddr().String())
				r.closeClientConnection(client)
				return
			}
			continue
		}
		// exec
		if payLoad.Date == nil {
			continue
		}
		multiBulkReply, ok := payLoad.Date.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("protocol error: expect multi bulk reply")
			continue
		}
		result := r.db.Exec(client, multiBulkReply.Args)
		if result != nil {
			// 多行的需要处理缓冲区溢出的错误了
			bytes, err := result.ToBytes()
			if err != nil {
				logger.Error("encode reply error:" + err.Error())
			}
			_ = client.Writer(bytes)
		} else {
			// 常量及错误类型无需处理error
			bytes, _ := reply.UnknownErrorReply{}.ToBytes()
			_ = client.Writer(bytes)
		}
	}
}

func (r *RespHandler) Close() error {
	logger.Info("handler shut down")
	r.closing.Set(true)
	r.activeConn.Range(func(key, value interface{}) bool {
		clientConnection := key.(*connection.ClientConnection)
		clientConnection.Close()
		return true
	})
	r.db.Close()
	return nil
}
