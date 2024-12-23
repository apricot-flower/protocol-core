package broker

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

const (
	startChar byte = 0x68
	endChar   byte = 0x16
)

type eventStatute struct {
	start1       byte
	eventId      uint64 //报文id
	idLength     uint16 //id长度
	id           string //服务端id
	start2       byte
	length       uint32 //长度域
	brokerLength uint16 //事件驱动的id
	broker       string //事件标志
	brokerData   string //信息
	needResponse byte   //是否需要回复，0-不需要，1-需要
	end          byte
}

// 解码
func (E *eventStatute) decode(data []byte) error {
	var err error
	E.start1 = startChar
	E.end = endChar
	reader := bytes.NewReader(data)
	err = binary.Read(reader, binary.BigEndian, &E.eventId)
	if err != nil {
		return err
	}
	//拆idLength
	err = binary.Read(reader, binary.BigEndian, &E.idLength)
	if err != nil {
		return err
	}
	idArray := make([]byte, E.idLength)
	err = binary.Read(reader, binary.BigEndian, &idArray)
	if err != nil {
		return err
	}
	E.id = string(idArray)
	//解析headCs
	headCsArrayLength := 8 + 2 + E.idLength
	headCsArray := data[0:headCsArrayLength]
	headCs := E.cs(headCsArray)
	var frameHeadCs byte
	err = binary.Read(reader, binary.BigEndian, &frameHeadCs)
	if err != nil {
		return err
	}
	if frameHeadCs != headCs {
		return errors.New("invalid frameHeadCs")
	}
	err = binary.Read(reader, binary.BigEndian, &E.start2)
	if err != nil {
		return err
	}
	if E.start2 != startChar {
		return errors.New("invalid start2")
	}
	//解析长度域
	err = binary.Read(reader, binary.BigEndian, &E.length)
	if err != nil {
		return err
	}
	//解析broker长度
	err = binary.Read(reader, binary.BigEndian, &E.brokerLength)
	if err != nil {
		return err
	}
	brokerArray := make([]byte, E.brokerLength)
	err = binary.Read(reader, binary.BigEndian, &brokerArray)
	if err != nil {
		return err
	}
	E.broker = string(brokerArray)
	if E.length-uint32(E.brokerLength) > 0 {
		dataArray := make([]byte, E.length-uint32(E.brokerLength))
		err = binary.Read(reader, binary.BigEndian, &dataArray)
		if err != nil {
			return err
		}
		E.brokerData = string(dataArray)
	}
	err = binary.Read(reader, binary.BigEndian, &E.needResponse)
	if err != nil {
		return err
	}
	//解析cs
	var cs byte
	err = binary.Read(reader, binary.BigEndian, &cs)
	if err != nil {
		return err
	}
	if E.cs(data[0:len(data)-1]) != cs {
		return errors.New("error cs number")
	}
	return nil
}

func (E *eventStatute) cs(data []byte) byte {
	var sum uint8
	for _, b := range data {
		sum += b
	}
	return sum
}

// 编码
func (E *eventStatute) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	var err error
	err = binary.Write(buf, binary.BigEndian, E.eventId)
	if err != nil {
		return nil, err
	}
	idArray := []byte(E.id)
	err = binary.Write(buf, binary.BigEndian, uint16(len(idArray)))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, idArray)
	if err != nil {
		return nil, err
	}
	//拼接一个cs
	headCs := E.cs(buf.Bytes())
	err = binary.Write(buf, binary.BigEndian, headCs)
	if err != nil {
		return nil, err
	}
	//拼接一个start
	err = binary.Write(buf, binary.BigEndian, startChar)
	E.length = 0
	var dataArray []byte
	//解析broker和brokerLength
	if E.broker != "" {
		brokerArray := []byte(E.broker)
		dataArray = append(dataArray, brokerArray...)
		E.brokerLength = uint16(len(brokerArray))
		E.length += uint32(E.brokerLength)
	}
	if E.brokerData != "" {
		messageArray := []byte(E.brokerData)
		dataArray = append(dataArray, messageArray...)
		E.length += uint32(len(messageArray))
	}
	dataArray = append(dataArray, E.needResponse)
	err = binary.Write(buf, binary.BigEndian, E.length)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, E.brokerLength)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, dataArray)
	if err != nil {
		return nil, err
	}
	//计算cs
	frame := buf.Bytes()
	cs := E.cs(frame)
	frame = append(frame, cs)
	return append(append([]byte{startChar}, frame...), endChar), nil
}

func newStatute(broker string, message string, needResponse bool, clientId string) *eventStatute {
	needResponseFlag := 0
	if needResponse {
		needResponseFlag = 1
	}
	statute := &eventStatute{
		start1:       startChar,
		end:          endChar,
		eventId:      uint64(time.Now().UnixNano()),
		start2:       startChar,
		id:           clientId,
		broker:       broker,
		brokerData:   message,
		needResponse: byte(needResponseFlag),
	}
	return statute
}
