package dlt1376_1

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	StartChar byte = 0x68
	EndChar   byte = 0x16
)

type LinkDataInter interface {
	decode(buf *bytes.Reader) error
	encode() ([]byte, error)
	haAux() bool
	AfnFlag() byte
}

// StatuteDlt13761 dlt1376.1 报文对象
type StatuteDlt13761 struct {
	startChar byte
	length    uint16
	Control   *ControlFiled     //控制域
	Address   *AddressFiled     //地址域
	Data      *ApplicationLayer //数据
	cs        byte              //校验位
	endChar   byte              //结束字符
}

// DecodeByStr 解码
func (s *StatuteDlt13761) DecodeByStr(frame string) error {
	frameStr := strings.ReplaceAll(frame, " ", "")
	frameArray, err := hex.DecodeString(frameStr)
	if err != nil {
		return errors.New("Error decoding dlt1376.1 frame hex string: " + err.Error())
	}
	return s.DecodeByArr(frameArray)
}

// DecodeByArr 解码
func (s *StatuteDlt13761) DecodeByArr(frame []byte) error {
	buf := bytes.NewReader(frame)
	err := binary.Read(buf, binary.LittleEndian, &s.startChar)
	if err != nil {
		return errors.New("Error encoding dlt1376.1 frame starting char1: " + err.Error())
	}
	if s.startChar != StartChar {
		return errors.New("error encoding dlt1376.1 frame starting char: start char1 != 0x68")
	}
	var length1, length2 uint16
	if err = binary.Read(buf, binary.LittleEndian, &length1); err != nil {
		return errors.New("decode dlt1376.1 length err:" + err.Error())
	}
	if err = binary.Read(buf, binary.LittleEndian, &length2); err != nil {
		return errors.New("decode dlt1376.1 length err:" + err.Error())
	}
	if length1 != length2 {
		return errors.New("decode dlt1376.1 length err: invalid length")
	}
	if err = s.lengthHandle(length1); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &s.startChar); err != nil {
		return errors.New("Error encoding dlt1376.1 frame starting char2: " + err.Error())
	}
	if s.startChar != StartChar {
		return errors.New("error encoding dlt1376.1 frame starting char: start char2 != 0x68")
	}
	//解析控制域
	s.Control = &ControlFiled{}
	if err = s.Control.decode(buf); err != nil {
		return err
	}
	//解析地址域
	s.Address = &AddressFiled{}
	if err = s.Address.decode(buf); err != nil {
		return err
	}
	//解析链路用户数据
	s.Data = &ApplicationLayer{}
	err = s.Data.decode(buf, s.length-6, s.Control.DIR == "1" && s.Control.FCBorACD == "1", s.Control.DIR)
	if err != nil {
		return err
	}
	csArray := frame[6 : len(frame)-2]
	checkCs := s.checkCs(csArray)
	if err = binary.Read(buf, binary.LittleEndian, &s.cs); err != nil {
		return err
	}
	if checkCs != s.cs {
		return errors.New("dlt1376.1 invalid cs")
	}
	var endChar byte
	if err = binary.Read(buf, binary.LittleEndian, &endChar); err != nil {
		return err
	}
	if endChar != EndChar {
		return errors.New("dlt1376.1 invalid end char")
	}
	return nil
}

func (s *StatuteDlt13761) Encode() ([]byte, error) {
	controlArray, err := s.Control.encode()
	if err != nil {
		return nil, err
	}
	addressArray, err := s.Address.encode()
	if err != nil {
		return nil, err
	}
	dataArray, err := s.Data.encode(s.Control.DIR, s.Control.DIR == "1" && s.Control.FCBorACD == "1")
	if err != nil {
		return nil, err
	}
	s.length = uint16(len(controlArray) + len(addressArray) + len(dataArray))
	//处理长度域
	s.length = s.length & 0x3FFF // 0x3FFF 是 14 位全 1 的掩码
	// 将提取的值转换为 14 位的二进制字符串
	binaryStr := fmt.Sprintf("%014b", s.length)
	// 在二进制字符串的末尾追加 "10"
	binaryStr += "10"
	// 将新的二进制字符串转换回 uint16
	newValue, err := strconv.ParseUint(binaryStr, 2, 16)
	if err != nil {
		return nil, err
	}
	s.length = uint16(newValue)
	encodeArray := []byte{StartChar, byte(s.length & 0xFF), byte(s.length >> 8), byte(s.length & 0xFF), byte(s.length >> 8), StartChar}
	userDataArray := append(controlArray, addressArray...)
	userDataArray = append(userDataArray, dataArray...)
	csArray := append(controlArray, addressArray...)
	s.cs = s.checkCs(append(csArray, dataArray...))
	userDataArray = append(userDataArray, s.cs)
	encodeArray = append(encodeArray, userDataArray...)
	return append(encodeArray, EndChar), nil
}

func (s *StatuteDlt13761) lengthHandle(length uint16) error {
	// 检查bit0是否为0
	bit0 := length & 0x0001
	if bit0 != 0 {
		return errors.New("dlt1376.1 invalid length bit0")
	}
	// 检查bit1是否为1
	bit1 := length & 0x0002
	if bit1 == 0 {
		return errors.New("dlt1376.1 invalid length bit1")
	}
	mask := uint16(0x3FFF) << 2
	maskedValue := length & mask
	s.length = maskedValue >> 2
	return nil
}

func (s *StatuteDlt13761) checkCs(buf []byte) byte {
	var sum uint8
	for _, b := range buf {
		sum += b
	}
	return sum
}

type ControlFiled struct {
	DIR      string `json:"dir"`        //传输方向位
	PRM      string `json:"prm"`        //启动标志位
	FCBorACD string `json:"fcb_or_acd"` // 帧计数位FCB 要求访问位ACD
	FCV      string `json:"fcv"`        //帧计数有效位
	FuncCode string `json:"func_code"`  //功能码
}

func (c *ControlFiled) decode(buf *bytes.Reader) error {
	control, err := buf.ReadByte()
	if err != nil {
		return errors.New("decode dlt1376.1 control err:" + err.Error())
	}
	c.DIR = strconv.Itoa(int((control >> 7) & 1))
	c.PRM = strconv.Itoa(int((control >> 6) & 1))
	c.FCBorACD = strconv.Itoa(int((control >> 5) & 1))
	c.FCV = strconv.Itoa(int((control >> 4) & 1))
	c.FuncCode = fmt.Sprintf("%d%d%d%d", (control>>3)&1, (control>>2)&1, (control>>1)&1, control&1)
	return nil
}

func (c *ControlFiled) encode() ([]byte, error) {
	bits := c.DIR + c.PRM + c.FCBorACD + c.FCV + c.FuncCode
	control, err := strconv.ParseUint(bits, 2, 8)
	if err != nil {
		return nil, errors.New("encode 1376.1 control err:" + err.Error())
	}
	return []byte{byte(control)}, nil
}

type AddressFiled struct {
	Area        []byte `json:"area"`         //行政区划码A1
	Terminal    []byte `json:"terminal"`     //终端地址A2
	AddressType string `json:"address_type"` //终端组地址标志,D0=0表示终端地址A2为单地址；D0=1表示终端地址A2为组地址
	MSA         string `json:"MSA"`          //A3的D1～D7组成0～127个主站地址MSA
}

func (a *AddressFiled) decode(buf *bytes.Reader) error {
	a.Area = make([]byte, 2)
	if err := binary.Read(buf, binary.LittleEndian, &a.Area); err != nil {
		return errors.New("Error decoding dlt1376.1 address: " + err.Error())
	}
	a.Terminal = make([]byte, 2)
	if err := binary.Read(buf, binary.LittleEndian, &a.Terminal); err != nil {
		return errors.New("Error decoding dlt1376.1 address: " + err.Error())
	}
	addressFlag, err := buf.ReadByte()
	if err != nil {
		return errors.New("Error decoding dlt1376.1 address: " + err.Error())
	}
	a.AddressType = strconv.Itoa(int(addressFlag & 1))
	a.MSA = fmt.Sprintf("%07b", (addressFlag>>1)&0b01111111)
	return nil
}

func (a *AddressFiled) overTurn(arr []byte) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func (a *AddressFiled) encode() ([]byte, error) {
	encodeArray := append(a.Area, a.Terminal...)
	AddressFlag := a.MSA + a.AddressType
	value, err := strconv.ParseUint(AddressFlag, 2, 8)
	if err != nil {
		return nil, errors.New("encode dlt1376.1 err:" + err.Error())
	}
	return append(encodeArray, byte(value)), nil
}

func (a *AddressFiled) Build(address string) error {
	if len(address) != 8 {
		return errors.New("build dlt1376.1 address: address length = 8")
	}
	addressArray, err := hex.DecodeString(address)
	if err != nil {
		return err
	}
	a.Area = []byte{addressArray[1], addressArray[0]}
	a.Terminal = []byte{addressArray[3], addressArray[2]}
	return nil
}

func (a *AddressFiled) GetValue() interface{} {
	area := a.Area[:]
	terminal := a.Terminal[:]
	a.overTurn(area)
	a.overTurn(terminal)
	arr := append(area, terminal...)
	return hex.EncodeToString(arr)
}

// ApplicationLayer 应用层
type ApplicationLayer struct {
	Afn  byte          //功能码
	Seq  *SEQ          //帧序列号
	Data LinkDataInter //数据
	Aux  *AUX          //附加信息域
}

func (a *ApplicationLayer) decode(buf *bytes.Reader, dataLength uint16, hasEc bool, dir string) error {
	err := binary.Read(buf, binary.LittleEndian, &a.Afn)
	if err != nil {
		return errors.New("Error decoding dlt1376.1 application layer: " + err.Error())
	}
	a.Seq = &SEQ{}
	err = a.Seq.decode(buf)
	if err != nil {
		return err
	}
	a.Data = translate(a.Afn)
	if a.Data == nil {
		return errors.New("dlt1376.1 application layer data type is empty")
	}
	err = a.Data.decode(buf)
	if err != nil {
		return err
	}
	if a.Data.haAux() {
		//存在附加信息域
		a.Aux = &AUX{}
		err = a.Aux.decode(buf, hasEc, dir, a.Seq.TpV, a.Afn)
	}
	return err
}

func (a *ApplicationLayer) encode(dir string, hasEc bool) ([]byte, error) {
	encodeArray := []byte{a.Data.AfnFlag()}
	seqArray, err := a.Seq.encode()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, seqArray...)
	dataArray, err := a.Data.encode()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, dataArray...)
	if a.Aux != nil {
		auxArray, err := a.Aux.encode(dir, hasEc, a.Seq.TpV, a.Data.AfnFlag())
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, auxArray...)
	}

	return encodeArray, nil
}

// SEQ 帧序列域SEQ
type SEQ struct {
	TpV        string `json:"tpv"` //帧时间标签有效位
	FIR        string `json:"fir"`
	FIN        string `json:"fin"`
	CON        string `json:"CON"`
	PSEQorRSEQ string `json:"PSEQorRSEQ"`
}

func (S *SEQ) decode(buf *bytes.Reader) error {
	seq, err := buf.ReadByte()
	if err != nil {
		return errors.New("Error decoding dlt1376.1 SEQ: " + err.Error())
	}
	S.TpV = strconv.Itoa(int((seq >> 7) & 1))
	S.FIR = strconv.Itoa(int((seq >> 6) & 1))
	S.FIN = strconv.Itoa(int((seq >> 5) & 1))
	S.CON = strconv.Itoa(int((seq >> 4) & 1))
	S.PSEQorRSEQ = fmt.Sprintf("%d%d%d%d", (seq>>3)&1, (seq>>2)&1, (seq>>1)&1, seq&1)
	return nil
}

func (S *SEQ) encode() ([]byte, error) {
	bits := S.TpV + S.FIR + S.FIN + S.CON + S.PSEQorRSEQ
	seq, err := strconv.ParseUint(bits, 2, 8)
	if err != nil {
		return nil, errors.New("encode 1376.1 SEQ err:" + err.Error())
	}
	return []byte{byte(seq)}, nil
}

type AUX struct {
	PW  []byte `json:"pw"`  //消息认证码字段PW用于重要下行报文，由16字节组成。没约定就传16个00
	EC1 byte   `json:"ec1"` //事件计数器EC用于ACD位置“1”的上行响应报文中，EC由2字节组成，分别为重要事件计数器EC1和一般事件计数器EC2。计数范围0～255，循环加1递增
	EC2 byte   `json:"ec2"` //事件计数器EC用于ACD位置“1”的上行响应报文中，EC由2字节组成，分别为重要事件计数器EC1和一般事件计数器EC2。计数范围0～255，循环加1递增
	TP  *Tp    `json:"TP"`  //时间标签域
}

func (a *AUX) decode(buf *bytes.Reader, hasEc bool, dir string, tpv string, afn byte) error {
	if dir == "0" && afn == ResetIdent {
		//下行 解析PW
		a.PW = make([]byte, 16)
		if err := binary.Read(buf, binary.LittleEndian, &a.PW); err != nil {
			return err
		}
	} else {
		if hasEc {
			//上行 解析EC
			if err := binary.Read(buf, binary.LittleEndian, &a.EC1); err != nil {
				return err
			}
			if err := binary.Read(buf, binary.LittleEndian, &a.EC2); err != nil {
				return err
			}
		}
	}
	if tpv == "1" {
		//解析时间标签域
		a.TP = &Tp{}
		if err := a.TP.decode(buf); err != nil {
			return err
		}
	}
	return nil
}

func (a *AUX) encode(dir string, hasEc bool, TpV string, afn byte) ([]byte, error) {
	var encodeArray []byte
	if dir == "0" {
		//下行
		if afn == ResetIdent {
			encodeArray = append(encodeArray, a.PW...)
		}
	} else {
		//上行

	}
	if TpV == "1" {
		if a.TP != nil {
			tpArray, err := a.TP.encode()
			if err != nil {
				return nil, err
			}
			return append(encodeArray, tpArray...), nil
		}
	}
	return encodeArray, nil
}

// Tp 时间标签
type Tp struct {
	PFC     byte `json:"pfc"` //启动帧帧序号计数器PFC
	Second  byte `json:"second"`
	Minute  byte `json:"minute"`
	Hour    byte `json:"hour"`
	Day     byte `json:"day"`
	Delayed byte `json:"delayed"` //允许发送传输延时时间
}

func (t *Tp) encode() ([]byte, error) {
	return []byte{t.PFC, t.Second, t.Minute, t.Hour, t.Day, t.Delayed}, nil
}

func (t *Tp) decode(buf *bytes.Reader) error {
	err := binary.Read(buf, binary.LittleEndian, t)
	if err != nil {
		return err
	}
	t.Second = t.extractDigits(t.Second)
	t.Minute = t.extractDigits(t.Minute)
	t.Hour = t.extractDigits(t.Hour)
	t.Day = t.extractDigits(t.Day)
	return err
}

func (t *Tp) extractDigits(b byte) byte {
	tens := (b >> 4) & 0xF
	units := b & 0xF
	result := tens*10 + units
	return result
}

func (t *Tp) BuildByNow() error {
	var err error
	now := time.Now()
	second := now.Second()
	t.Second, err = t.montage(second)
	if err != nil {
		return errors.New("build TP err:" + err.Error())
	}
	minute := now.Minute()
	t.Minute, err = t.montage(minute)
	if err != nil {
		return errors.New("build TP err:" + err.Error())
	}
	hour := now.Hour()
	t.Hour, err = t.montage(hour)
	if err != nil {
		return errors.New("build TP err:" + err.Error())
	}
	day := now.Day()
	t.Day, err = t.montage(day)
	if err != nil {
		return errors.New("build TP err:" + err.Error())
	}
	return nil
}

func (t *Tp) montage(data int) (byte, error) {
	value := byte(data)
	tens := value / 10  // 十位
	units := value % 10 // 个位
	binStr := fmt.Sprintf("%04b", tens) + fmt.Sprintf("%04b", units)
	result, err := strconv.ParseUint(binStr, 2, 8)
	if err != nil {
		return 0, err
	}
	return byte(result), nil
}
