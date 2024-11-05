package dlt698

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var csTabs []int

func init() {
	csTabs = []int{0x00, 0x1189, 0x2312, 0x329b, 0x4624,
		0x57ad, 0x6536, 0x74bf, 0x8c48, 0x9dc1, 0xaf5a, 0xbed3, 0xca6c,
		0xdbe5, 0xe97e, 0xf8f7, 0x1081, 0x108, 0x3393, 0x221a, 0x56a5,
		0x472c, 0x75b7, 0x643e, 0x9cc9, 0x8d40, 0xbfdb, 0xae52, 0xdaed,
		0xcb64, 0xf9ff, 0xe876, 0x2102, 0x308b, 0x210, 0x1399, 0x6726,
		0x76af, 0x4434, 0x55bd, 0xad4a, 0xbcc3, 0x8e58, 0x9fd1, 0xeb6e,
		0xfae7, 0xc87c, 0xd9f5, 0x3183, 0x200a, 0x1291, 0x318, 0x77a7,
		0x662e, 0x54b5, 0x453c, 0xbdcb, 0xac42, 0x9ed9, 0x8f50, 0xfbef,
		0xea66, 0xd8fd, 0xc974, 0x4204, 0x538d, 0x6116, 0x709f, 0x420,
		0x15a9, 0x2732, 0x36bb, 0xce4c, 0xdfc5, 0xed5e, 0xfcd7, 0x8868,
		0x99e1, 0xab7a, 0xbaf3, 0x5285, 0x430c, 0x7197, 0x601e, 0x14a1,
		0x528, 0x37b3, 0x263a, 0xdecd, 0xcf44, 0xfddf, 0xec56, 0x98e9,
		0x8960, 0xbbfb, 0xaa72, 0x6306, 0x728f, 0x4014, 0x519d, 0x2522,
		0x34ab, 0x630, 0x17b9, 0xef4e, 0xfec7, 0xcc5c, 0xddd5, 0xa96a,
		0xb8e3, 0x8a78, 0x9bf1, 0x7387, 0x620e, 0x5095, 0x411c, 0x35a3,
		0x242a, 0x16b1, 0x738, 0xffcf, 0xee46, 0xdcdd, 0xcd54, 0xb9eb,
		0xa862, 0x9af9, 0x8b70, 0x8408, 0x9581, 0xa71a, 0xb693, 0xc22c,
		0xd3a5, 0xe13e, 0xf0b7, 0x840, 0x19c9, 0x2b52, 0x3adb, 0x4e64,
		0x5fed, 0x6d76, 0x7cff, 0x9489, 0x8500, 0xb79b, 0xa612, 0xd2ad,
		0xc324, 0xf1bf, 0xe036, 0x18c1, 0x948, 0x3bd3, 0x2a5a, 0x5ee5,
		0x4f6c, 0x7df7, 0x6c7e, 0xa50a, 0xb483, 0x8618, 0x9791, 0xe32e,
		0xf2a7, 0xc03c, 0xd1b5, 0x2942, 0x38cb, 0xa50, 0x1bd9, 0x6f66,
		0x7eef, 0x4c74, 0x5dfd, 0xb58b, 0xa402, 0x9699, 0x8710, 0xf3af,
		0xe226, 0xd0bd, 0xc134, 0x39c3, 0x284a, 0x1ad1, 0xb58, 0x7fe7,
		0x6e6e, 0x5cf5, 0x4d7c, 0xc60c, 0xd785, 0xe51e, 0xf497, 0x8028,
		0x91a1, 0xa33a, 0xb2b3, 0x4a44, 0x5bcd, 0x6956, 0x78df, 0xc60,
		0x1de9, 0x2f72, 0x3efb, 0xd68d, 0xc704, 0xf59f, 0xe416, 0x90a9,
		0x8120, 0xb3bb, 0xa232, 0x5ac5, 0x4b4c, 0x79d7, 0x685e, 0x1ce1,
		0xd68, 0x3ff3, 0x2e7a, 0xe70e, 0xf687, 0xc41c, 0xd595, 0xa12a,
		0xb0a3, 0x8238, 0x93b1, 0x6b46, 0x7acf, 0x4854, 0x59dd, 0x2d62,
		0x3ceb, 0xe70, 0x1ff9, 0xf78f, 0xe606, 0xd49d, 0xc514, 0xb1ab,
		0xa022, 0x92b9, 0x8330, 0x7bc7, 0x6a4e, 0x58d5, 0x495c, 0x3de3,
		0x2c6a, 0x1ef1, 0xf78}
}

const (
	StartChar byte = 0x68
	EndChar   byte = 0x16
	ScCode    byte = 0x33 //加扰码
)

var _ FrameRegion = (*ControlRegion)(nil)
var _ FrameRegion = (*AddressRegion)(nil)
var _ FrameRegion = (*APDU)(nil)
var _ FrameRegion = (*TimeTag)(nil)
var _ FrameRegion = (*FollowReport)(nil)

// FrameRegion 数据域
type FrameRegion interface {
	// Decoder 解码
	decoder(buf *bytes.Reader) error
	// Encoder 编码
	encoder() ([]byte, error)
}

type APDURegion interface {
	FrameRegion
	APDUType() string //获取类型
	APDUMark() string //获取类型
	hasFollowReport() bool
	hasTimeTag() bool
}

type DataInter interface {
	FrameRegion
	DataType() byte
	Value() interface{}
}

// ProtocolDlt645Model dlt698.45 模型
type ProtocolDlt645Model struct {
	Length  uint16         `json:"length"`         //长度域
	Control *ControlRegion `json:"control_region"` //控制域
	Address *AddressRegion `json:"address_region"` //地址域
	Data    *APDU          `json:"data_region"`    //数据域
}

// DecodeByStr 根据字符串解析
func (p *ProtocolDlt645Model) DecodeByStr(str string) error {
	arr, err := hex.DecodeString(strings.ReplaceAll(str, " ", ""))
	if err != nil {
		return err
	}
	return p.DecodeByBytes(arr)
}

// DecodeByBytes 根据字节数组解析
func (p *ProtocolDlt645Model) DecodeByBytes(frame []byte) error {
	buf := bytes.NewReader(frame)
	//解析起始字符
	var startChar byte
	if err := binary.Read(buf, binary.BigEndian, &startChar); err != nil {
		return errors.New("decode startChar err:" + err.Error())
	}
	if startChar != StartChar {
		return errors.New("decode startChar err: start char err, != 0x68")
	}
	//解析长度域
	if err := binary.Read(buf, binary.LittleEndian, &p.Length); err != nil {
		return errors.New("decode frame length err:" + err.Error())
	}
	if (p.Length & 0x6000) != 0 {
		return errors.New("frame length's bit14 and bit15 != 0！")
	}
	//解析控制域
	p.Control = &ControlRegion{}
	if err := p.Control.decoder(buf); err != nil {
		return errors.New("decode control err:" + err.Error())
	}
	//解析地址域
	p.Address = &AddressRegion{}
	if err := p.Address.decoder(buf); err != nil {
		return errors.New("decode address err:" + err.Error())
	}
	//校验帧头校验
	hcsLen := 2 + 1 + 1 + p.Address.AddressLength + 1
	checkHcs := p.Cs(frame[1 : hcsLen+1])
	hcs := make([]byte, 2)
	if err := binary.Read(buf, binary.BigEndian, &hcs); err != nil {
		return errors.New("decode HCS err:" + err.Error())
	}
	if checkHcs[0] != hcs[0] || checkHcs[1] != hcs[1] {
		return errors.New("decode HCS err！")
	}
	//长度域 - 控制域长度 - 地址域长度 - hcs长度 - fcs长度
	dataLen := p.Length - 2 - 1 - 1 - uint16(p.Address.AddressLength) - 1 - 2 - 2
	apduArray := make([]byte, dataLen)
	if err := binary.Read(buf, binary.BigEndian, &apduArray); err != nil {
		return errors.New("decode APDU err:" + err.Error())
	}
	//计算帧校验
	checkFcs := p.Cs(frame[1 : len(frame)-3])
	fcs := make([]byte, 2)
	if err := binary.Read(buf, binary.LittleEndian, &fcs); err != nil {
		return errors.New("decode FCS err:" + err.Error())
	}
	if checkFcs[0] != fcs[0] || checkFcs[1] != fcs[1] {
		return errors.New("decode FCS err！")
	}
	var endChar byte
	if err := binary.Read(buf, binary.BigEndian, &endChar); err != nil {
		return errors.New("decode endChar err:" + err.Error())
	}
	if endChar != EndChar {
		return errors.New("decode endChar err: endChar err, != 0x16")
	}
	//解析数据域
	if p.Control.Sc == "1" {
		for index, value := range apduArray {
			apduArray[index] = value - ScCode
		}
	}
	p.Data = &APDU{}
	err := p.Data.decoder(bytes.NewReader(apduArray))
	return err
}

func (p *ProtocolDlt645Model) Encoder() ([]byte, error) {
	controlArray, err := p.Control.encoder()
	if err != nil {
		return nil, err
	}
	addressArray, err := p.Address.encoder()
	if err != nil {
		return nil, err
	}
	dataArray, err := p.Data.encoder()
	if err != nil {
		return nil, err
	}
	if p.Control.Sc == "1" {
		for index, value := range dataArray {
			dataArray[index] = value + ScCode
		}
	}
	p.Length = uint16(len(controlArray) + len(addressArray) + len(dataArray) + 4 + 2)
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, p.Length)
	if err != nil {
		return nil, errors.New("encode length err:" + err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, controlArray)
	if err != nil {
		return nil, errors.New("encode control err:" + err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, addressArray)
	if err != nil {
		return nil, errors.New("encode address err:" + err.Error())
	}
	//编码HCS
	hcs := p.Cs(append(buf.Bytes(), 0x00))
	err = binary.Write(buf, binary.LittleEndian, hcs)
	if err != nil {
		return nil, errors.New("encode HCS err:" + err.Error())
	}
	err = binary.Write(buf, binary.LittleEndian, dataArray)
	if err != nil {
		return nil, errors.New("encode apdu err:" + err.Error())
	}
	//编码FCS
	fcs := p.Cs(append(buf.Bytes(), 0x00))
	err = binary.Write(buf, binary.LittleEndian, fcs)
	if err != nil {
		return nil, errors.New("encode FCS err:" + err.Error())
	}
	encodeArray := buf.Bytes()
	encodeArray = append(encodeArray, EndChar)
	encodeArray = append([]byte{StartChar}, encodeArray...)
	return encodeArray, nil
}

func (p *ProtocolDlt645Model) Cs(data []byte) []byte {
	data = append(data, 0x00, 0x00)
	fcs := p.cs16(data)
	fcs ^= 0xffff
	data[len(data)-2] = byte(fcs & 0x00ff)
	data[len(data)-2+1] = byte((fcs >> 8) & 0x00ff)
	c := fcs & 0x00ff
	c1 := (fcs >> 8) & 0x00ff
	return []byte{byte(c), byte(c1)}
}

func (p *ProtocolDlt645Model) cs16(data []byte) int {
	length := len(data) - 2
	index := 0
	CS16 := 0xffff
	for {
		length--
		if length == 0 {
			break
		}
		CS16 = (CS16 >> 8) ^ csTabs[(CS16^int(data[index]))&0xff]
		index++
	}
	return CS16
}

// ControlRegion 控制域
type ControlRegion struct {
	Dir     string `json:"dir"`     //传输方向位： bit7=0 表示此帧是由客户机发出的；bit7=1 表示此帧是由服务器发出的
	Prm     string `json:"prm"`     //启动标志位：bit6=0 表示此帧是由服务器发起的；bit6=1 表示此帧是由客户机发起的
	Framing string `json:"framing"` //分帧标志
	Sc      string `json:"sc"`      //扰码标志
	Func    string `json:"func"`    //功能码
}

func (c *ControlRegion) decoder(buf *bytes.Reader) error {
	var control byte
	if err := binary.Read(buf, binary.BigEndian, &control); err != nil {
		return errors.New("decode control err:" + err.Error())
	}
	c.Dir = strconv.Itoa(int((control >> 7) & 1))
	c.Prm = strconv.Itoa(int((control >> 6) & 1))
	c.Framing = strconv.Itoa(int((control >> 5) & 1))
	c.Sc = strconv.Itoa(int(control >> 3 & 1))
	bits := (control >> 0) & 0b111
	c.Func = fmt.Sprintf("%d%d%d", (bits>>2)&1, (bits>>1)&1, bits&1)
	return nil
}

func (c *ControlRegion) encoder() ([]byte, error) {
	bits := c.Dir + c.Prm + c.Framing + "0" + c.Sc + c.Func
	control, err := strconv.ParseUint(bits, 2, 8)
	if err != nil {
		return nil, errors.New("encode control region err:" + err.Error())
	}
	return []byte{byte(control)}, nil
}

type AddressRegion struct {
	AddressType   uint8  `json:"address_type"`   //地址类型
	LogicAddress  uint8  `json:"logic_address"`  //逻辑地址
	AddressLength uint8  `json:"address_length"` //地址长度 索引0
	Address       string `json:"address"`        //地址
	CA            byte   `json:"ca"`             //客户机地址
}

func (a *AddressRegion) decoder(buf *bytes.Reader) error {
	var addressFeatures byte
	if err := binary.Read(buf, binary.BigEndian, &addressFeatures); err != nil {
		return errors.New("decode address' address_features err:" + err.Error())
	}
	a.AddressType = (addressFeatures >> 6) & 0b11
	if a.AddressType != 0 && a.AddressType != 1 && a.AddressType != 2 && a.AddressType != 3 {
		return errors.New("decode address' address_type err: not in (0,1,2,3)")
	}
	//解析逻辑地址
	a.LogicAddress = (addressFeatures >> 4) & 0b11
	a.AddressLength = ((addressFeatures >> 0) & 0b1111) + 1
	addressArray := make([]byte, a.AddressLength)
	if err := binary.Read(buf, binary.BigEndian, &addressArray); err != nil {
		return errors.New("decode address' address err:" + err.Error())
	}
	//翻转数组
	for i, j := 0, len(addressArray)-1; i < j; i, j = i+1, j-1 {
		addressArray[i], addressArray[j] = addressArray[j], addressArray[i]
	}
	a.Address = hex.EncodeToString(addressArray)
	if err := binary.Read(buf, binary.BigEndian, &a.CA); err != nil {
		return errors.New("decode address' CA err:" + err.Error())
	}
	return nil
}

func (a *AddressRegion) encoder() ([]byte, error) {
	addressTypeStr := fmt.Sprintf("%d%d", (a.AddressType>>1)&1, a.AddressType&1)
	logicAddressStr := fmt.Sprintf("%d%d", (a.LogicAddress>>1)&1, a.LogicAddress&1)
	length := len(a.Address)/2 - 1
	addressLengthStr := strconv.FormatInt(int64(length&0b1111), 2)
	adType, err := strconv.ParseUint(addressTypeStr+logicAddressStr+addressLengthStr, 2, 8)
	if err != nil {
		return nil, err
	}
	addressArray, err := hex.DecodeString(a.Address)
	if err != nil {
		return nil, err
	}
	for i, j := 0, len(addressArray)-1; i < j; i, j = i+1, j-1 {
		addressArray[i], addressArray[j] = addressArray[j], addressArray[i]
	}
	return append(append([]byte{byte(adType)}, addressArray...), []byte{a.CA}...), nil
}

// APDU 链路用户数据
type APDU struct {
	Pid          byte          `json:"pid"`
	Data         APDURegion    `json:"data"`         //链路用户数据
	TimeTag      *TimeTag      `json:"timeTag"`      //时间标签
	FollowReport *FollowReport `json:"followReport"` //附加信息域
}

func (a *APDU) decoder(buf *bytes.Reader) error {
	//先读取一个字节
	var apdu1 byte
	if err := binary.Read(buf, binary.BigEndian, &apdu1); err != nil {
		return errors.New("decode APDU type err:" + err.Error())
	}
	a.Data = translate(apdu1)
	if a.Data == nil {
		var apdu2 byte
		if err := binary.Read(buf, binary.BigEndian, &apdu2); err != nil {
			return errors.New("decode APDU type err:" + err.Error())
		}
		a.Data = translate(apdu1, apdu2)
		if a.Data == nil {
			return errors.New("decode APDU type err, not such type")
		}
	}
	if err := binary.Read(buf, binary.BigEndian, &a.Pid); err != nil {
		return errors.New("decode APDU pid err:" + err.Error())
	}
	if err := a.Data.decoder(buf); err != nil {
		return err
	}
	if a.Data.hasFollowReport() {
		//存在附加信息域
	}
	if a.Data.hasTimeTag() {
		//存在时间标签
		var hasTimeTag byte
		if err := binary.Read(buf, binary.BigEndian, &hasTimeTag); err != nil {
			return errors.New("decode timeTag exist flag err:" + err.Error())
		}
		if hasTimeTag == 0x01 {
			a.TimeTag = &TimeTag{}
			if err := a.TimeTag.decoder(buf); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *APDU) encoder() ([]byte, error) {
	if a.Data == nil {
		return nil, errors.New("apdu data == nil")
	}
	encodeArray, err := hex.DecodeString(a.Data.APDUType())
	if err != nil {
		return nil, errors.New("encode APDU type err:" + err.Error())
	}
	encodeArray = append(encodeArray, a.Pid)
	dataArray, err := a.Data.encoder()
	if err != nil {
		return nil, errors.New("encode APDU data err:" + err.Error())
	}
	encodeArray = append(encodeArray, dataArray...)
	//解析附加信息域
	if a.Data.hasFollowReport() {

	}
	//解析时间标签
	if a.Data.hasTimeTag() {
		if a.TimeTag != nil {
			encodeArray = append(encodeArray, 0x01)
			timeTagArray, err := a.TimeTag.encoder()
			if err != nil {
				return nil, err
			} else {
				encodeArray = append(encodeArray, timeTagArray...)
			}
		} else {
			encodeArray = append(encodeArray, 0x00)
		}
	}

	return encodeArray, nil
}

/*----------------------时间标签-----------------------------*/

type TimeTag struct {
	SendTime *DateTimes `xml:"send_time"` //发送时标
	Ti       *TI        `xml:"ti"`        //允许传输延时时间
}

func (t *TimeTag) decoder(buf *bytes.Reader) error {
	t.SendTime = &DateTimes{}
	if err := t.SendTime.decoder(buf); err != nil {
		return err
	}
	t.Ti = &TI{}
	return t.Ti.decoder(buf)
}

func (t *TimeTag) encoder() ([]byte, error) {
	sendTimeArray, err := t.SendTime.encoder()
	if err != nil {
		return nil, err
	}
	tiArray, err := t.Ti.encoder()
	if err != nil {
		return nil, err
	}
	return append(sendTimeArray, tiArray...), nil
}

type FollowReport struct {
}

func (f FollowReport) decoder(buf *bytes.Reader) error {
	return nil
}

func (f FollowReport) encoder() ([]byte, error) {
	return nil, nil
}
