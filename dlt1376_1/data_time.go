package dlt1376_1

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

type DataInter interface {
	decode(buf *bytes.Reader) error
	encode() ([]byte, error)
}

type TimeInter interface {
	DataInter
	Build() error
}

var _ TimeInter = (*TimeA1)(nil)

type TimeA1 struct {
	Second byte //秒
	Minute byte //分
	Hour   byte //时
	Day    byte //日
	Week   byte //星期
	Month  byte //月
	Year   byte //年
}

func (t *TimeA1) decode(buf *bytes.Reader) error {
	err := binary.Read(buf, binary.LittleEndian, &t.Second)
	if err != nil {
		return errors.New("decode A.1 time type Second err:" + err.Error())
	}
	low, high := t.splitByte(t.Second)
	t.Second = high*10 + low
	err = binary.Read(buf, binary.LittleEndian, &t.Minute)
	if err != nil {
		return errors.New("decode A.1 time type Minute err:" + err.Error())
	}
	low, high = t.splitByte(t.Minute)
	t.Minute = high*10 + low
	err = binary.Read(buf, binary.LittleEndian, &t.Hour)
	if err != nil {
		return errors.New("decode A.1 time type Hour err:" + err.Error())
	}
	low, high = t.splitByte(t.Hour)
	t.Hour = high*10 + low
	err = binary.Read(buf, binary.LittleEndian, &t.Day)
	if err != nil {
		return errors.New("decode A.1 time type Day err:" + err.Error())
	}
	low, high = t.splitByte(t.Day)
	t.Day = high*10 + low
	var weekMonth byte
	err = binary.Read(buf, binary.LittleEndian, &weekMonth)
	if err != nil {
		return errors.New("decode A.1 time type weekMonth err:" + err.Error())
	}
	bits0to3 := weekMonth & 0x0F
	// 提取 bit4
	bit4 := (weekMonth >> 4) & 0x01
	// 提取 bit7~bit5
	bits7to5 := (weekMonth >> 5) & 0x07
	t.Week = bits7to5
	t.Month = bit4*10 + bits0to3
	err = binary.Read(buf, binary.LittleEndian, &t.Year)
	if err != nil {
		return errors.New("decode A.1 time type year err:" + err.Error())
	}
	low, high = t.splitByte(t.Year)
	t.Year = high*10 + low
	return nil
}

func (t *TimeA1) encode() ([]byte, error) {
	encodeArray := []byte{
		((t.Second / 10) << 4) | t.Second%10,
		((t.Minute / 10) << 4) | t.Minute%10,
		((t.Hour / 10) << 4) | t.Hour%10,
		((t.Day / 10) << 4) | t.Day%10}
	// 提取 a 的低4位
	lowBitsA := (t.Month % 10) & 0x0F // 0x0F 即 15，在二进制为 1111
	// 提取 b 的 bit0
	bit0B := (t.Month / 10) & 0x01 // 0x01 在二进制为 00000001
	// 提取 c 的低3位
	lowBitsC := t.Week & 0x07 // 0x07 在二进制为 00000111
	// 将 a 的低4位左移5位，以腾出位置给 b 的 bit0 和 c 的低3位
	shiftedLowBitsA := lowBitsA << 5
	// 将 b 的 bit0 左移3位，以腾出位置给 c 的低3位
	shiftedBit0B := bit0B << 3
	// 组合所有部分
	combined := shiftedLowBitsA | shiftedBit0B | lowBitsC
	encodeArray = append(encodeArray, byte(combined))
	return append(encodeArray, ((t.Year/10)<<4)|t.Year%10), nil
}

func (t *TimeA1) Build() error {
	now := time.Now()
	// 提取年、月、日、时、分、秒
	year, month, day := now.Date()
	hour, minute, sec := now.Clock()
	weekday := now.Weekday()
	t.Year = byte(year)
	t.Week = byte(weekday)
	t.Month = byte(month)
	t.Day = byte(day)
	t.Hour = byte(hour)
	t.Minute = byte(minute)
	t.Second = byte(sec)
	return nil
}

func (t *TimeA1) splitByte(b byte) (low, high byte) {
	// 0xF 即 15 在二进制表示为 1111，用于提取低4位
	low = b & 0x0F // 提取低4位
	// 右移4位以获得高4位，然后同样使用 0xF 来确保只保留最后4位
	high = (b >> 4) & 0x0F
	return low, high
}
