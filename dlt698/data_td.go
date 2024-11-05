package dlt698

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

var _ DataInter = (*DateTime)(nil)
var _ DataInter = (*DateTimes)(nil)
var _ DataInter = (*TI)(nil)
var _ DataInter = (*Date)(nil)
var _ DataInter = (*Time)(nil)

const (
	TIIdent        byte = 84
	DataTimeIdent  byte = 25
	DateIdent      byte = 26
	TimeIdent      byte = 27
	DateTimesIdent byte = 28
)

func init() {
	dataMap[DataTimeIdent] = func() DataInter {
		return &DateTime{}
	}
	dataMap[DateTimesIdent] = func() DataInter {
		return &DateTimes{}
	}
	dataMap[TIIdent] = func() DataInter {
		return &TI{}
	}
	dataMap[DateIdent] = func() DataInter {
		return &Date{}
	}
	dataMap[TimeIdent] = func() DataInter {
		return &Time{}
	}
}

/*--------------------------------------------------*/

type Time struct {
	Hour   uint8 `json:"hour"`
	Minute uint8 `json:"minute"`
	Second uint8 `json:"second"`
}

func (t *Time) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, t); err != nil {
		return errors.New("decode data<Time> err " + err.Error())
	}
	return nil
}

func (t *Time) encoder() ([]byte, error) {
	return []byte{t.Hour, t.Minute, t.Second}, nil
}

func (t *Time) DataType() byte {
	return TimeIdent
}

func (t *Time) Value() interface{} {
	return t
}

/*-------------------------------------------------*/

type Date struct {
	Year       uint16 `json:"year"`
	Month      uint8  `json:"month"`
	DayOfMonth uint8  `json:"dayOfMonth"`
	DayOfWeek  uint8  `json:"dayOfWeek"`
}

func (d *Date) Build() *Date {
	now := time.Now()
	year, month, day := now.Date()
	dayOfWeek := now.Weekday()
	d.Year = uint16(year)
	d.Month = uint8(month)
	d.DayOfMonth = uint8(day)
	d.DayOfWeek = uint8(dayOfWeek)
	return d
}

func (d *Date) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, d); err != nil {
		return errors.New("decode data<Date> err " + err.Error())
	}
	return nil
}

func (d *Date) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, d)
	if err != nil {
		return nil, errors.New("encode data<Date> err " + err.Error())
	}
	return buf.Bytes(), nil
}

func (d *Date) DataType() byte {
	return DateIdent
}

func (d *Date) Value() interface{} {
	return d
}

/*-----------------------------------------------------------------*/

type DateTime struct {
	Year        uint16 `json:"year"`
	Month       uint8  `json:"month"`
	DayOfMonth  uint8  `json:"dayOfMonth"`
	DayOfWeek   uint8  `json:"dayOfWeek"`
	Hour        uint8  `json:"hour"`
	Minute      uint8  `json:"minute"`
	Second      uint8  `json:"second"`
	Millisecond uint16 `json:"millisecond"`
}

func (d *DateTime) Build() *DateTime {
	now := time.Now()
	year, month, day := now.Date()
	hour, minute, second := now.Clock()
	millisecond := now.Nanosecond() / 1e6
	dayOfWeek := now.Weekday()
	d.Year = uint16(year)
	d.Month = uint8(month)
	d.DayOfMonth = uint8(day)
	d.DayOfWeek = uint8(dayOfWeek)
	d.Hour = uint8(hour)
	d.Minute = uint8(minute)
	d.Second = uint8(second)
	d.Millisecond = uint16(millisecond)
	return d
}

func (d *DateTime) decoder(buf *bytes.Reader) error {
	err := binary.Read(buf, binary.BigEndian, d)
	if err != nil {
		return errors.New("decode date_time err :" + err.Error())
	}
	return nil
}

func (d *DateTime) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, d)
	if err != nil {
		return nil, errors.New("encode date_time err :" + err.Error())
	}
	return buf.Bytes(), nil
}

func (d *DateTime) DataType() byte {
	return DataTimeIdent
}

func (d *DateTime) Value() interface{} {
	return d
}

/*------------dateTimes---------------*/

type DateTimes struct {
	Year   uint16 `json:"year"`
	Month  uint8  `json:"month"`
	Day    uint8  `json:"day"`
	Hour   uint8  `json:"hour"`
	Minute uint8  `json:"minute"`
	Second uint8  `json:"second"`
}

func (d *DateTimes) Build() *DateTimes {
	now := time.Now()
	year, month, day := now.Date()
	hour, minute, second := now.Clock()
	d.Year = uint16(year)
	d.Month = uint8(month)
	d.Day = uint8(day)
	d.Hour = uint8(hour)
	d.Minute = uint8(minute)
	d.Second = uint8(second)
	return d
}

func (d *DateTimes) decoder(buf *bytes.Reader) error {
	err := binary.Read(buf, binary.BigEndian, d)
	if err != nil {
		return errors.New("decode date_times err :" + err.Error())
	}
	return nil
}

func (d *DateTimes) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, d)
	if err != nil {
		return nil, errors.New("encode date_times err :" + err.Error())
	}
	return buf.Bytes(), nil
}

func (d *DateTimes) DataType() byte {
	return DateTimesIdent
}

func (d *DateTimes) Value() interface{} {
	return d
}

/*-------------------TI-------------------*/

type TI struct {
	TimeUnit byte   //时间单位
	Interval uint16 //间隔时间
}

func (T *TI) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, &T.TimeUnit); err != nil {
		return errors.New("decode TI's TimeUnit err:" + err.Error())
	}
	if T.TimeUnit > 5 {
		return errors.New("TI's TimeUnit must be < 5")
	}
	err := binary.Read(buf, binary.BigEndian, &T.Interval)
	if err != nil {
		return errors.New("decode TI's Interval err:" + err.Error())
	}
	return nil
}

func (T *TI) encoder() ([]byte, error) {
	if T.TimeUnit > 5 {
		return nil, errors.New("TI's TimeUnit must be < 5")
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, T.TimeUnit)
	err = binary.Write(buf, binary.BigEndian, T.Interval)
	if err != nil {
		return nil, errors.New("encode TI's Interval err:" + err.Error())
	}
	return buf.Bytes(), nil
}

func (T *TI) DataType() byte {
	return TIIdent
}

func (T *TI) Value() interface{} {
	return T
}
