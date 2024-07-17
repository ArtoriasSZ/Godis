# Godis
Implementing Redis in Go
## Godis Serialization Protocol Specification
### 类似redis 
https://redis.io/topics/protocol

#### 五种reply类型
- 正常回复：'+'开头，"\r\n"结尾
```asciidoc
+OK\r\n
```
- 错误恢复：'-'开头，"\r\n"结尾
```asciidoc
-ERR unknown command 'foobar'\r\n
```
- 整数回复：':'开头，"\r\n"结尾
```asciidoc
:1000\r\n
```
- 字符串回复："$"开头，后面是长度，"\r\n"结尾，后面是数据，"\r\n"结尾
```asciidoc
$6\r\nfoobar\r\n
```
- 字符串数组回复："*"开头，后面是长度，"\r\n"结尾，后面是多个字符串，每个批量回复之间用"\r\n"分隔
```asciidoc
*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n
```

#### 空回复定义
#### 错误回复定义

## 代码逻辑
### tcp/server.go
- 实现了tcp服务器，监听端口，接收客户端请求，调用实现handler接口的结构体异步处理请求
- 开辟协程监听系统关闭信号，当系统关闭信号到达通过管道发送给g0协程进行关闭,
- 在死循环中阻塞接受新的tcp连接，得到一个新连接就给waitGroup++，然后调用handler接口处理该连接。

### resp/parser.go
- 实现了resp协议的解析，将客户端发送的字符串解析成resp协议的reply类型