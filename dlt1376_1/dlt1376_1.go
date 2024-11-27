package dlt1376_1

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	DLT13761_START_CHAR byte = 0x68
	DLT13761_END_CHAR   byte = 0x16
)

type Dlt13761statute struct {
	length  uint16        //长度域
	Control *ControlFiled `json:"control"` //控制域
	Address *AddressFiled `json:"address"` //地址域
	Data    *LinkUserData `json:"data"`    //链路用户数据
	cs      byte
}

func (d *Dlt13761statute) Decode(frame []byte) error {
	var err error
	buf := bytes.NewReader(frame)
	var startChar uint8
	if err = binary.Read(buf, binary.LittleEndian, &startChar); err != nil {
		return errors.New("decode dlt1376.1 startChar err:" + err.Error())
	}
	if startChar != DLT13761_START_CHAR {
		return errors.New("decode dlt1376.1 frame err: invalid start char")
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
	if err = d.lengthHandle(length1); err != nil {
		return err
	}
	if err = binary.Read(buf, binary.LittleEndian, &startChar); err != nil {
		return err
	}
	if startChar != DLT13761_START_CHAR {
		return errors.New("decode dlt1376.1 startChar err: invalid start char")
	}
	d.Control = &ControlFiled{}
	if err = d.Control.Decode(buf); err != nil {
		return err
	}
	d.Address = &AddressFiled{}
	if err = d.Address.Decode(buf); err != nil {
		return err
	}
	d.Data = &LinkUserData{}
	if err = d.Data.Decode(buf, d.length-6, d.Control.DIR == "1" && d.Control.FCBorACD == "1", d.Control.DIR, d.Control.FCBorACD); err != nil {
		return err
	}
	csArray := frame[6 : len(frame)-2]
	checkCs := d.checkCs(csArray)
	if err = binary.Read(buf, binary.LittleEndian, &d.cs); err != nil {
		return err
	}
	if checkCs != d.cs {
		return errors.New("invalid cs")
	}
	var endChar byte
	if err = binary.Read(buf, binary.LittleEndian, &endChar); err != nil {
		return err
	}
	if endChar != DLT13761_END_CHAR {
		return errors.New("invalid end char")
	}
	return nil
}

func (d *Dlt13761statute) Encode() ([]byte, error) {
	controlArray, err := d.Control.Encode()
	if err != nil {
		return nil, err
	}
	addressArray, err := d.Address.Encode()
	if err != nil {
		return nil, err
	}
	dataArray, err := d.Data.Encode(d.Control.DIR)
	if err != nil {
		return nil, err
	}
	d.length = uint16(len(controlArray) + len(addressArray) + len(dataArray))
	//处理长度域
	d.length = d.length & 0x3FFF // 0x3FFF 是 14 位全 1 的掩码
	// 将提取的值转换为 14 位的二进制字符串
	binaryStr := fmt.Sprintf("%014b", d.length)
	// 在二进制字符串的末尾追加 "10"
	binaryStr += "10"
	// 将新的二进制字符串转换回 uint16
	newValue, err := strconv.ParseUint(binaryStr, 2, 16)
	if err != nil {
		return nil, err
	}
	d.length = uint16(newValue)
	encodeArray := []byte{DLT13761_START_CHAR, byte(d.length & 0xFF), byte(d.length >> 8), byte(d.length & 0xFF), byte(d.length >> 8), DLT13761_START_CHAR}
	userDataArray := append(controlArray, addressArray...)
	userDataArray = append(userDataArray, dataArray...)
	csArray := append(controlArray, addressArray...)
	d.cs = d.checkCs(append(csArray, dataArray...))
	userDataArray = append(userDataArray, d.cs)
	encodeArray = append(encodeArray, userDataArray...)
	return append(encodeArray, DLT13761_END_CHAR), nil

}

func (d *Dlt13761statute) checkCs(buf []byte) byte {
	var sum uint8
	for _, b := range buf {
		sum += b
	}
	return sum
}

func (d *Dlt13761statute) lengthHandle(length uint16) error {
	// 检查bit0是否为0
	bit0 := length & 0x0001
	if bit0 != 0 {
		return errors.New("invalid length bit0")
	}
	// 检查bit1是否为1
	bit1 := length & 0x0002
	if bit1 == 0 {
		return errors.New("invalid length bit1")
	}
	mask := uint16(0x3FFF) << 2
	maskedValue := length & mask
	d.length = maskedValue >> 2
	return nil
}

type ControlFiled struct {
	DIR      string `json:"dir"`        //传输方向位
	PRM      string `json:"prm"`        //启动标志位
	FCBorACD string `json:"fcb_or_acd"` // 帧计数位FCB 要求访问位ACD
	FCV      string `json:"fcv"`        //帧计数有效位
	FuncCode string `json:"func_code"`  //功能码
}

func (c *ControlFiled) Decode(buf *bytes.Reader) error {
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

func (c *ControlFiled) Encode() ([]byte, error) {
	bits := c.DIR + c.PRM + c.FCBorACD + c.FCV + c.FuncCode
	control, err := strconv.ParseUint(bits, 2, 8)
	if err != nil {
		return nil, err
	}
	return []byte{byte(control)}, nil
}

type AddressFiled struct {
	Area        []byte `json:"area"`         //行政区划码A1
	Terminal    []byte `json:"terminal"`     //终端地址A2
	AddressType string `json:"address_type"` //终端组地址标志,D0=0表示终端地址A2为单地址；D0=1表示终端地址A2为组地址
	MSA         string `json:"MSA"`          //A3的D1～D7组成0～127个主站地址MSA
}

func (a *AddressFiled) Decode(buf *bytes.Reader) error {
	a.Area = make([]byte, 2)
	if err := binary.Read(buf, binary.LittleEndian, &a.Area); err != nil {
		return err
	}
	a.Terminal = make([]byte, 2)
	if err := binary.Read(buf, binary.LittleEndian, &a.Terminal); err != nil {
		return err
	}
	addressFlag, err := buf.ReadByte()
	if err != nil {
		return err
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

func (a *AddressFiled) Encode() ([]byte, error) {
	encodeArray := append(a.Area, a.Terminal...)
	AddressFlag := a.MSA + a.AddressType
	value, err := strconv.ParseUint(AddressFlag, 2, 8)
	if err != nil {
		return nil, err
	}
	return append(encodeArray, byte(value)), nil
}

func (a *AddressFiled) Build(address string) error {
	if len(address) != 8 {
		return errors.New("invalid address")
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

// LinkUserData 链路用户数据
type LinkUserData struct {
	AFN      byte `json:"afn"`       //应用层功能码
	Seq      *SEQ `json:"seq"`       //帧序列域
	DataUnit Afn  `json:"data_unit"` //数据单元
	Aux      *AUX `json:"aux"`       //附加信息域
}

func (l *LinkUserData) Decode(buf *bytes.Reader, dataLength uint16, hasTp bool, dir string, acd string) error {
	var err error
	if err = binary.Read(buf, binary.LittleEndian, &l.AFN); err != nil {
		return errors.New("decode 1376.1 AFN err:" + err.Error())
	}
	l.Seq = &SEQ{}
	if err = l.Seq.Decode(buf); err != nil {
		return err
	}
	if l.Seq.TpV == "1" {
		dataLength = dataLength - 6
	}
	if hasTp {
		dataLength = dataLength - 2
	}
	dataArray := make([]byte, dataLength-2)
	if err := binary.Read(buf, binary.LittleEndian, dataArray); err != nil {
		return err
	}
	if dataLength < 0 {
		return errors.New("invalid data length < 0")
	}
	if dataLength == 0 {
		return nil
	}
	dataBuf := bytes.NewReader(dataArray)
	l.DataUnit = Translate(l.AFN)
	if l.DataUnit == nil {
		return errors.New("invalid data unit")
	}
	//中继命令绑定方向
	if du, ok := l.DataUnit.(*RelayStationCommand); ok {
		du.Direction(dir)
	}
	err = l.DataUnit.Decode(dataBuf)
	if err != nil {
		return err
	}
	//先判断是否有附加信息域
	if l.DataUnit.HasAux() {
		//存在附加信息域
		l.Aux = &AUX{}
		err = l.Aux.Decode(dataBuf, hasTp, dir, acd)
	}
	return err
}

func (l *LinkUserData) Encode(dir string) ([]byte, error) {
	if l.DataUnit == nil {
		return nil, errors.New("invalid data unit, must data != nil")
	}
	l.AFN, _ = l.DataUnit.Flag()
	encodeArray := []byte{l.AFN}
	seqArray, err := l.Seq.Encode()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, seqArray...)
	//解析dataUnit
	unitArray, err := l.DataUnit.Encode()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, unitArray...)
	if l.DataUnit.HasAux() {
		if l.Aux == nil {
			return nil, errors.New("invalid data unit, must has aux")
		}
		auxArray, err := l.Aux.Encode(dir)
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

func (S *SEQ) Decode(buf *bytes.Reader) error {
	seq, err := buf.ReadByte()
	if err != nil {
		return err
	}
	S.TpV = strconv.Itoa(int((seq >> 7) & 1))
	S.FIR = strconv.Itoa(int((seq >> 6) & 1))
	S.FIN = strconv.Itoa(int((seq >> 5) & 1))
	S.CON = strconv.Itoa(int((seq >> 4) & 1))
	S.PSEQorRSEQ = fmt.Sprintf("%d%d%d%d", (seq>>3)&1, (seq>>2)&1, (seq>>1)&1, seq&1)
	return nil
}

func (S *SEQ) Encode() ([]byte, error) {
	bits := S.TpV + S.FIR + S.FIN + S.CON + S.PSEQorRSEQ
	seq, err := strconv.ParseUint(bits, 2, 8)
	if err != nil {
		return nil, err
	}
	return []byte{byte(seq)}, nil
}

type AUX struct {
	PW  []byte `json:"pw"`  //消息认证码字段PW用于重要下行报文，由16字节组成。没约定就传16个00
	EC1 byte   `json:"ec1"` //事件计数器EC用于ACD位置“1”的上行响应报文中，EC由2字节组成，分别为重要事件计数器EC1和一般事件计数器EC2。计数范围0～255，循环加1递增
	EC2 byte   `json:"ec2"` //事件计数器EC用于ACD位置“1”的上行响应报文中，EC由2字节组成，分别为重要事件计数器EC1和一般事件计数器EC2。计数范围0～255，循环加1递增
	TP  *Tp    `json:"TP"`  //时间标签域
}

func (a *AUX) Decode(buf *bytes.Reader, hasTp bool, dir string, acd string) error {
	if dir == "0" {
		//下行 解析PW
		a.PW = make([]byte, 16)
		if err := binary.Read(buf, binary.LittleEndian, a.PW); err != nil {
			return err
		}
	} else {
		if acd == "1" {
			//上行 解析EC
			if err := binary.Read(buf, binary.LittleEndian, &a.EC1); err != nil {
				return err
			}
			if err := binary.Read(buf, binary.LittleEndian, &a.EC2); err != nil {
				return err
			}
		}
	}
	if hasTp {
		//解析时间标签域
		a.TP = &Tp{}
		if err := a.TP.Decode(buf); err != nil {
			return err
		}
	}
	return nil
}

func (a *AUX) Encode(dir string) ([]byte, error) {
	var encodeArray []byte
	if dir == "0" {
		//下行 编码消息认证码字段PW
		encodeArray = a.PW
	} else {
		//TODO 上行 编码事件计数器EC（上行）
	}
	//解析时间标签
	if a.TP != nil {
		tpArray, err := a.TP.Encode()
		if err != nil {
			return nil, err
		}
		return append(encodeArray, tpArray...), nil
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

func (t *Tp) Encode() ([]byte, error) {
	return []byte{t.PFC, t.Second, t.Minute, t.Hour, t.Day, t.Delayed}, nil
}

func (t *Tp) Decode(buf *bytes.Reader) error {
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
		return err
	}
	minute := now.Minute()
	t.Minute, err = t.montage(minute)
	if err != nil {
		return err
	}
	hour := now.Hour()
	t.Hour, err = t.montage(hour)
	if err != nil {
		return err
	}
	day := now.Day()
	t.Day, err = t.montage(day)
	return err
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
