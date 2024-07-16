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
	Maxfile uint32
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	// 系统关闭信号到来进行关闭
	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
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
	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		waitDone.Add(1)
		logger.Info("accept link")
		go func() {
			defer waitDone.Done()
			handler.Handler(ctx, conn)
		}()
	}
	waitDone.Wait()
	return nil
}
