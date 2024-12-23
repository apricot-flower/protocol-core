package broker

import "net"

// ConnectedHandler 连接成功之后的回调
type ConnectedHandler func(localAddr net.Addr, remoteAddr net.Addr)

// DisconnectHandler 断开连接之后或者无法连接的回调
type DisconnectHandler func(serverHost string, serverPort int, err error)

// EventMessageHandler 事件分发回调
type EventMessageHandler func(clientId, message string, messageId uint64, needResponse bool)

// ReceivedMessageHandler 缺省的事件分发回调
type ReceivedMessageHandler func(clientId, broker, message string, messageId uint64, needResponse bool)

// ClientConnectedHandler 有客户端连接进来的回调,返回true就代表这个连接是有效的，false就是无效的
type ClientConnectedHandler func(ipAndPort string) bool

// ClientDisconnectedHandler 有客户端断开连接的回调
type ClientDisconnectedHandler func(ipAndPort string)
