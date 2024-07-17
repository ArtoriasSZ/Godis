package tcp

import (
	"Godis/interface/tcp"
	"Godis/lib/logger"
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

func ListenAndServeWithSystemSignal(cfg *Config, handler tcp.Handler) error {
	// 系统关闭信号到来进行关闭
	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			sig := <-sigChan
			switch sig {
			case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				// 处理信号的逻辑
				closeChan <- struct{}{}
				return
			default:
				// 处理未预期的信号或者忽略它们
			}
		}
	}()
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info("start listen")
	err = ListenAndServe(listener, handler, make(chan struct{}))
	return nil
}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) error {
	go func() {
		// 系统关闭信号
		<-closeChan
		logger.Info("shutting down")
		_ = listener.Close()
		_ = handler.Close()
	}()
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	//deadline, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*10))
	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			// 如果listener关闭，就退出循环，等待所有客户端断开连接，
			break
		}
		// 每一个新的连接就+1
		waitDone.Add(1)
		logger.Info("accept link")
		go func() {
			//
			defer waitDone.Done()
			handler.Handler(ctx, conn)
		}()
	}
	waitDone.Wait()
	return nil
}
