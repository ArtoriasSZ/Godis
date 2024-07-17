package main

import (
	"Godis/config"
	"Godis/lib/logger"
	"Godis/resp/handler"
	"Godis/tcp"
	"fmt"
	"io"
	"os"
)

// 配置文件
const configFile = "redis.conf"

// 默认配置
var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

// 配置文件是否存在
func fileExist(fileName string) bool {
	info, err := os.Stat(fileName)
	return err == nil && !info.IsDir()
}

// MockReader implements io.Reader for testing purposes.
type MockReader struct {
	data []byte
}

func (r *MockReader) Read(p []byte) (n int, err error) {
	copy(p, r.data)
	return len(r.data), io.EOF
}
func main() {
	//goland:noinspection SpellCheckingInspection
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "Godis",
		Ext:        "log",
		TimeFormat: "2003-06-16",
	})
	if fileExist(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}
	err := tcp.ListenAndServeWithSystemSignal(
		&tcp.Config{
			Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
		}, handler.NewRespHandler(),
	)
	if err != nil {
		logger.Error(err)
	}

}
