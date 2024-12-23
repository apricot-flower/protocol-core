package broker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
)

// 断开连接的通道
var disConnectCh chan string
var messageCh chan *eventStatute

//收到客户端消息的通道

func init() {
	disConnectCh = make(chan string, 5)
	messageCh = make(chan *eventStatute, 5)
}

type Server struct {
	clientId                     string                    //标志
	port                         int                       //端口
	blackList                    []string                  //黑名单
	ClientConnectedHandler       ClientConnectedHandler    //有客户端连接进来的回调
	ReceivedClientMessageHandler interface{}               //收到客户端的消息
	ClientDisconnectedHandler    ClientDisconnectedHandler //有客户端断开连接的回调
	listener                     net.Listener
	clients                      sync.Map
	messageHandler               map[string]EventMessageHandler //消息分发
	ReceivedMessageHandler       ReceivedMessageHandler         //缺省的消息分发
	closeCh                      chan struct{}
}

func BuildBazookaServer(port int) *Server {
	return &Server{
		clientId:       "FEFEFEFEFEFE",
		port:           port,
		blackList:      make([]string, 0),
		clients:        sync.Map{},
		messageHandler: make(map[string]EventMessageHandler),
		closeCh:        make(chan struct{}),
	}
}

// Send 发送
func (s *Server) Send(clientIpPort, broker, message string, needResponse bool) (id uint64, err error) {
	statute := newStatute(broker, message, needResponse, s.clientId)
	request, err := statute.encode()
	if err != nil {
		return 0, err
	}
	if conn, ok := s.clients.Load(clientIpPort); ok {
		_, err = conn.(net.Conn).Write(request)
		if err != nil {
			return 0, err
		}
		return statute.eventId, nil
	} else {
		return 0, errors.New("no such connect client:" + clientIpPort)
	}
}

func (s *Server) appendBlack(blacks ...string) {
	s.blackList = append(s.blackList, blacks...)
}

// 添加一个broker
func (s *Server) appendMessageHandler(broker string, handler EventMessageHandler) {
	s.messageHandler[broker] = handler
}

// Close 关闭server
func (s *Server) Close() {
	close(s.closeCh)
	_ = s.listener.Close()
	s.clients.Range(func(k, v interface{}) bool {
		client := v.(*clientConn)
		close(client.closeCh)
		return true
	})
	s.clients.Clear()
}

func (s *Server) Open() error {
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(s.port))
	if err != nil {
		return err
	}

	go s.messageDistribution()

	go func() {
		for disConnect := range disConnectCh {
			//关闭连接
			if conn, ok := s.clients.Load(disConnect); ok {
				_ = conn.(net.Conn).Close()
				s.clients.Delete(disConnect)
			}
			//通知回调
			if s.ClientDisconnectedHandler != nil {
				s.ClientDisconnectedHandler(disConnect)
			}
		}
	}()
	s.listener = listen
	go func() {
	cons:
		for {
			select {
			case <-s.closeCh:
				break cons
			default:
				conn, err := listen.Accept()
				if err != nil {
					fmt.Println("error accepting connection:" + err.Error())
					continue
				}
				ipPort := conn.RemoteAddr().String()
				for _, black := range s.blackList {
					if strings.HasPrefix(ipPort, black) {
						_ = conn.Close()
						continue cons
					}
				}
				go s.clientConn(conn, ipPort)
			}
		}
	}()
	return nil
}

func (s *Server) clientConn(conn net.Conn, ipPort string) {
	if s.ClientConnectedHandler != nil {
		flag := s.ClientConnectedHandler(ipPort)
		if !flag {
			return
		}
	}
	client := &clientConn{
		ipPort:  ipPort,
		conn:    conn,
		reader:  bufio.NewReader(conn),
		decoder: &Decoder{},
		closeCh: make(chan struct{}),
	}
	go client.syncRead()
	s.clients.Store(ipPort, client)
}

// 消息分发
func (s *Server) messageDistribution() {
	for statute := range messageCh {
		needResponse := false
		if statute.needResponse == 1 {
			needResponse = true
		}
		if handler, ok := s.messageHandler[statute.broker]; ok {
			handler(statute.id, statute.brokerData, statute.eventId, needResponse)
		} else {
			if s.ReceivedMessageHandler != nil {
				s.ReceivedMessageHandler(statute.id, statute.broker, statute.brokerData, statute.eventId, needResponse)
			}
		}
	}
}

type clientConn struct {
	ipPort  string
	conn    net.Conn
	reader  *bufio.Reader
	decoder *Decoder
	closeCh chan struct{}
}

func (c *clientConn) syncRead() {
	for {
		select {
		case <-c.closeCh:
			_ = c.conn.Close()
			break
		default:
			frameArray, err := c.decoder.Decode(c.reader)
			fmt.Println("client read err:" + err.Error())
			if err != nil && err == io.EOF {
				//通知删除连接
				disConnectCh <- c.ipPort
				return
			}
			statute := &eventStatute{}
			err = statute.decode(frameArray)
			if err != nil {
				continue
			}
			messageCh <- statute
		}
	}
}
