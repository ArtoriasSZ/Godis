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
## code

