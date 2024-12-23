# bazooka
事件驱动式TCP交互框架式消息总线

# client
#### 新建客户端
```go
client := broker.BuildBazookaClient("127.0.0.1", 8314, "202306")
```
#### 添加钩子
###### 连接成功的钩子
```go
client.ConnectedHandler = connectedHandler
```
```go
func connectedHandler(localAddr net.Addr, remoteAddr net.Addr) {
	fmt.Println("Connected")
}
```
###### 断开连接、连接不成功的钩子
```go
client.DisconnectHandler = disconnectHandler
```
```go
func disconnectHandler(serverHost string, serverPort int, err error) {
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println("Disconnect")
	fmt.Println(err.Error())
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!")
}

```

###### 添加分发函数
```go
client.AddMessageHandler("/jojo", func(clientId, message string, messageId uint64, needResponse bool) {
		fmt.Println("------------")
		fmt.Println(clientId)
		fmt.Println(message)
		fmt.Println(needResponse)
		fmt.Println(messageId)
		fmt.Println("------------")
	})
```

###### 异步启动
```go
client.SyncConnect()
```

###### 发送
```go
messageId, err := client.SendMessage(broker, 消息, 消息接收方是否需要回复)
```


# server
###### 创建
```go
	server := broker.BuildBazookaServer(8314)
```
###### 添加黑名单
```go
appendBlack(string)
```
###### 添加一个消息监控器
```go
appendMessageHandler(broker string, handler EventMessageHandler)
```
###### 关闭
```go
Close()
```

###### 打开
```go
_ = server.Open()
```