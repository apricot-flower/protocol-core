package dlt645

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
)

const (
	StartChar byte = 0x68
	EndChar   byte = 0x16
	SCRAMB    byte = 0x33
)

// Dlt645statute dlt645-2007规约主体
type Dlt645statute struct {
	startChar1 byte     //帧起始符
	Address    string   `json:"address"` //地址域
	startChar2 byte     //帧起始符
	Control    byte     `json:"control"` //控制域
	Length     byte     `json:"length"`  //长度域
	Ident      string   `json:"ident"`   //数据标识
	Data       Analyser `json:"data"`    //数据域
	cs         byte     //校验码
	endChar    byte     //结束符
}

func (d *Dlt645statute) Decode(frame []byte) error {
	var err error
	buf := bytes.NewReader(frame)
	if err = binary.Read(buf, binary.LittleEndian, &d.startChar1); err != nil {
		return err
	}
	if d.startChar1 != StartChar {
		return errors.New("invalid start char")
	}
	//处理地址域
	addressArray := make([]byte, 6)
	if err = binary.Read(buf, binary.LittleEndian, &addressArray); err != nil {
		return err
	}
	d.overTurn(addressArray)
	d.Address = hex.EncodeToString(addressArray)
	if err = binary.Read(buf, binary.LittleEndian, &d.startChar2); err != nil {
		return err
	}
	if d.startChar2 != StartChar {
		return errors.New("invalid start char")
	}
	if err = binary.Read(buf, binary.LittleEndian, &d.Control); err != nil {
		return err
	}
	//解析length
	if err = binary.Read(buf, binary.LittleEndian, &d.Length); err != nil {
		return err
	}
	if d.Length != 0 {
		dataArray := make([]byte, d.Length)
		if err = binary.Read(buf, binary.LittleEndian, &dataArray); err != nil {
			return err
		}
		for i, value := range dataArray {
			dataArray[i] = value - SCRAMB
		}
		err = d.analyze(bytes.NewReader(dataArray))
		if err != nil {
			return err
		}
	}
	//处理cs
	if err = binary.Read(buf, binary.LittleEndian, &d.cs); err != nil {
		return err
	}
	if d.cs != d.Cs(frame[:len(frame)-2]) {
		return errors.New("cs error")
	}
	if err = binary.Read(buf, binary.LittleEndian, &d.endChar); err != nil {
		return err
	}
	if d.endChar != EndChar {
		return errors.New("invalid end char")
	}
	return nil
}

func (d *Dlt645statute) analyze(buf *bytes.Reader) error {
	switch d.Control {
	case MasterReadMeterData:
		return d.masterReadMeterData(buf)
	case MasterReadMeterNormalResponseNoNext, MasterReadMeterNormalResponseNext:
		return d.normalResponse(buf, false)
	case MasterReadMeterDataErrorResponse, MasterReadNextError, MasterRequestMeterSetFloatDataErrorResponse:
		return d.analyzeBytes(buf, 1)
	case MasterReadNext:
		return d.analyzeByte2(buf)
	case MasterReadNextNormalResponseFloatNext, MasterReadNextNormalResponseFloatNoNext:
		return d.normalResponse(buf, true)
	case MasterRequestMeterSetData:
		//主站向从站请求设置数据
		return d.masterRequestMeterSetData(buf)
	case MasterReadMeterAddressNormalResponse, MasterWriteMeterAddress, MasterWriteMeterAddressNormalResponse:
		return d.analyzeBytes(buf, 6)
	default:
		return errors.New("invalid control")
	}
}

func (d *Dlt645statute) masterRequestMeterSetData(buf *bytes.Reader) error {
	var err error
	d.Ident, err = IdentHandle(buf)
	if err != nil {
		return err
	}
	length := d.Length - 4
	dataArray := make([]byte, length)
	if err = binary.Read(buf, binary.LittleEndian, &dataArray); err != nil {
		return err
	}
	d.Data = &WriteAnalyzer{}
	err = d.Data.Decode(dataArray)
	return err
}

func (d *Dlt645statute) analyzeByte2(buf *bytes.Reader) error {
	if d.Length != 5 {
		return errors.New("invalid length")
	}
	var err error
	d.Ident, err = IdentHandle(buf)
	if err != nil {
		return err
	}
	var data byte
	if err := binary.Read(buf, binary.LittleEndian, &data); err != nil {
		return err
	}
	d.Data = &ByteArrayAnalyzer{Value: []byte{data}}
	return nil
}

func (d *Dlt645statute) analyzeBytes(buf *bytes.Reader, length byte) error {
	if d.Length != length {
		return errors.New("invalid length")
	}
	data := make([]byte, length)
	if err := binary.Read(buf, binary.LittleEndian, &data); err != nil {
		return err
	}
	d.Data = &ByteArrayAnalyzer{Value: data}
	return nil
}

func (d *Dlt645statute) normalResponse(buf *bytes.Reader, hasSeq bool) error {
	var err error
	d.Ident, err = IdentHandle(buf)
	if err != nil {
		return err
	}
	length := d.Length - 4
	data := make([]byte, length)
	if err = binary.Read(buf, binary.LittleEndian, &data); err != nil {
		return err
	}
	analyzer := Translate(d.Ident)
	if analyzer == nil {
		d.Data = nil
		return nil
	}
	if hasSeq {
		err = analyzer.Decode(data[:len(data)-1])
		if err != nil {
			return err
		}
		if an, ok := analyzer.(*FloatAnalyser); ok {
			an.Mark = []byte{data[len(data)-1]}
		}
	}
	err = analyzer.Decode(data)
	if err != nil {
		return err
	}
	d.Data = analyzer
	return nil
}

func (d *Dlt645statute) masterReadMeterData(buf *bytes.Reader) error {
	var err error
	d.Ident, err = IdentHandle(buf)
	if err != nil {
		return err
	}
	if d.Length == 4 {
		return nil
	} else if d.Length == 5 {
		//存在读给定块数的负荷记录
		var data byte
		if err = binary.Read(buf, binary.LittleEndian, &data); err != nil {
			return err
		}
		d.Data = &ByteArrayAnalyzer{Value: []byte{data}}
	} else if d.Length == 10 {
		data := make([]byte, 6)
		if err = binary.Read(buf, binary.LittleEndian, &data); err != nil {
			return err
		}
		d.Data = &ByteArrayAnalyzer{Value: data}
	} else {
		return errors.New("请求读电能表数据报文中数据域长度错误！")
	}
	return nil

}

func (d *Dlt645statute) Encode() ([]byte, error) {
	encodeArray := []byte{StartChar}
	if addressArray, err := hex.DecodeString(d.Address); err != nil {
		return nil, err
	} else {
		d.overTurn(addressArray)
		encodeArray = append(encodeArray, addressArray...)
	}
	encodeArray = append(encodeArray, StartChar, d.Control)
	dataArray := make([]byte, 0)
	if d.Ident != "" {
		identArray, err := hex.DecodeString(d.Ident)
		if err != nil {
			return nil, err
		}
		d.overTurn(identArray)
		dataArray = append(dataArray, identArray...)
	}
	if d.Data != nil {
		ddArray, err := d.Data.Encode()
		if err != nil {
			return nil, err
		}
		dataArray = append(dataArray, ddArray...)
	}
	d.Length = byte(len(dataArray))
	encodeArray = append(encodeArray, d.Length)
	for i, value := range dataArray {
		dataArray[i] = value + SCRAMB
	}
	encodeArray = append(encodeArray, dataArray...)
	encodeArray = append(encodeArray, d.Cs(encodeArray), EndChar)
	return encodeArray, nil
}

func (d *Dlt645statute) Cs(frame []byte) byte {
	var sum uint32 // 使用 uint32 来避免溢出
	for _, b := range frame {
		sum += uint32(b)
	}
	return byte(sum % 256)
}

func (d *Dlt645statute) overTurn(arr []byte) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}
