package broker

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

type Client struct {
	connectFlag            bool //是否连接成功
	linkConn               net.Conn
	lock                   sync.Mutex //锁
	clientId               string     //客户机id
	host                   string
	port                   int
	HeartBeatInterval      time.Duration                  //心跳周期，单位秒
	KeepAliveInterval      time.Duration                  //监测时间，单位秒
	ReConnectedInterval    time.Duration                  //断线重连时间，单位秒
	ConnectTime            time.Duration                  //连接超时时间，单位秒
	ConnectedHandler       ConnectedHandler               //连接成功之后的回调
	DisconnectHandler      DisconnectHandler              //断开连接之后或者无法连接的回调
	ReceivedMessageHandler ReceivedMessageHandler         //收到消息后的回调
	brokers                map[string]EventMessageHandler //事件-报文
	decoder                *Decoder                       //解码器
}

// 连接
func (c *Client) connect() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.connectFlag = false
	c.connectFlag = false
	dialer := &net.Dialer{
		Timeout:   c.ConnectTime,       // 设置连接超时时间
		KeepAlive: c.KeepAliveInterval, // 设置TCP Keep-alive的时间
	}
	conn, err := dialer.Dial("tcp", c.host+":"+strconv.Itoa(c.port))
	if err != nil {
		if c.DisconnectHandler != nil {
			c.DisconnectHandler(c.host, c.port, err)
		}
		return err
	}
	c.linkConn = conn
	c.connectFlag = true
	if c.ConnectedHandler != nil {
		go c.ConnectedHandler(conn.LocalAddr(), conn.RemoteAddr())
	}
	go c.read()
	return nil
}

// SyncConnect 异步连接
func (c *Client) SyncConnect() {
	go func() {
		for {
			err := c.connect()
			if err != nil {
				time.Sleep(c.ReConnectedInterval)
				continue
			}
			return
		}
	}()
}

func (c *Client) read() {
	reader := bufio.NewReader(c.linkConn)
	for {
		frameArray, err := c.decoder.Decode(reader)
		if err != nil {
			//处理错误
			c.errorHandler(err)
			return
		}
		statute := &eventStatute{}
		err = statute.decode(frameArray)
		if err != nil {
			continue
		}
		needResponse := false
		if statute.needResponse == 1 {
			needResponse = true
		}
		if handle, ok := c.brokers[statute.broker]; ok {
			go handle(statute.id, statute.brokerData, statute.eventId, needResponse)
		} else {
			go c.ReceivedMessageHandler(statute.id, statute.broker, statute.brokerData, statute.eventId, needResponse)
		}
	}
}

func (c *Client) AddMessageHandler(broker string, handler EventMessageHandler) {
	c.brokers[broker] = handler
}

func (c *Client) SendMessage(broker, message string, needResponse bool) (id uint64, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.connectFlag {
		return 0, errors.New("not connected client")
	}
	statute := newStatute(broker, message, needResponse, c.clientId)
	request, err := statute.encode()
	if err != nil {
		return 0, err
	}
	_, err = c.linkConn.Write(request)
	if err != nil {
		go c.errorHandler(err)
		return 0, err
	} else {
		return statute.eventId, nil
	}
}

// 处理异常
func (c *Client) errorHandler(err error) {
	if err == io.EOF {
		//断线重连
		c.SyncConnect()
	}
}

// BuildBazookaClient 创建一个基础的客户端
func BuildBazookaClient(host string, port int, clientId string) *Client {
	client := &Client{
		connectFlag:         false,
		host:                host,
		port:                port,
		clientId:            clientId,
		decoder:             &Decoder{},
		HeartBeatInterval:   60 * time.Second,
		KeepAliveInterval:   30 * time.Second,
		ReConnectedInterval: 15 * time.Second,
		ConnectTime:         20 * time.Second,
		brokers:             make(map[string]EventMessageHandler),
	}
	return client
}
