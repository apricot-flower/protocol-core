package dlt698

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var _ DataInter = (*Null)(nil)
var _ DataInter = (*Array)(nil)
var _ DataInter = (*Structure)(nil)
var _ DataInter = (*Boolean)(nil)
var _ DataInter = (*DoubleLong)(nil)
var _ DataInter = (*DoubleLongUnsigned)(nil)
var _ DataInter = (*Integer)(nil)
var _ DataInter = (*Long)(nil)
var _ DataInter = (*Unsigned)(nil)
var _ DataInter = (*LongUnsigned)(nil)
var _ DataInter = (*Enum)(nil)
var _ DataInter = (*Double32)(nil)
var _ DataInter = (*Double64)(nil)

const (
	NullIdent               byte = 0
	ArrayIdent              byte = 1
	StructureIdent          byte = 2
	BooleanIdent            byte = 3
	DoubleLongIdent         byte = 5
	DoubleLongUnsignedIdent byte = 6
	IntegerIdent            byte = 15
	LongIdent               byte = 16
	UnsignedIdent           byte = 17
	LongUnsignedIdent       byte = 18
	EnumIdent               byte = 22
	Double32Ident           byte = 23
	Double64Ident           byte = 24
)

func init() {
	dataMap[NullIdent] = func() DataInter {
		return &Null{}
	}
	dataMap[ArrayIdent] = func() DataInter {
		return &Array{}
	}
	dataMap[StructureIdent] = func() DataInter {
		return &Structure{}
	}
	dataMap[BooleanIdent] = func() DataInter {
		return &Boolean{}
	}
	dataMap[DoubleLongIdent] = func() DataInter {
		return &DoubleLong{}
	}
	dataMap[DoubleLongUnsignedIdent] = func() DataInter {
		return &DoubleLongUnsigned{}
	}
	dataMap[IntegerIdent] = func() DataInter {
		return &Integer{}
	}
	dataMap[LongIdent] = func() DataInter {
		return &Long{}
	}
	dataMap[UnsignedIdent] = func() DataInter {
		return &Unsigned{}
	}
	dataMap[LongUnsignedIdent] = func() DataInter {
		return &LongUnsigned{}
	}
	dataMap[EnumIdent] = func() DataInter {
		return &Enum{}
	}
	dataMap[Double32Ident] = func() DataInter {
		return &Double32{}
	}
	dataMap[Double64Ident] = func() DataInter {
		return &Double64{}
	}
}

/*-------------------------------------------------*/

type Double64 struct {
	Data float64 `json:"data"`
}

func (d *Double64) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, &d.Data); err != nil {
		return errors.New("decode data<Double64> err: " + err.Error())
	}
	return nil
}

func (d *Double64) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, d.Data)
	if err != nil {
		return nil, errors.New("encode data<Double64> err: " + err.Error())
	}
	return buf.Bytes(), nil
}

func (d *Double64) DataType() byte {
	return Double64Ident
}

func (d *Double64) Value() interface{} {
	return d.Data
}

/*----------------------------------------------------------*/

type Double32 struct {
	Data float32 `json:"data"`
}

func (d *Double32) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, &d.Data); err != nil {
		return errors.New("decode data<Double32> err: " + err.Error())
	}
	return nil
}

func (d *Double32) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, d.Data)
	if err != nil {
		return nil, errors.New("encode data<Double32> err: " + err.Error())
	}
	return buf.Bytes(), nil
}

func (d *Double32) DataType() byte {
	return Double32Ident
}

func (d *Double32) Value() interface{} {
	return d.Data
}

/*------------------------------------------*/

type Enum struct {
	Data uint8 `json:"data"`
}

func (e *Enum) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, &e.Data); err != nil {
		return errors.New("decode data<Enum> err: " + err.Error())
	}
	return nil
}

func (e *Enum) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, e.Data)
	if err != nil {
		return nil, errors.New("encode data<Enum> err: " + err.Error())
	}
	return buf.Bytes(), nil
}

func (e *Enum) DataType() byte {
	return EnumIdent
}

func (e *Enum) Value() interface{} {
	return e.Data
}

/*------------------------------------------------*/

type LongUnsigned struct {
	Data uint16 `json:"data"`
}

func (l *LongUnsigned) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, &l.Data); err != nil {
		return errors.New("decode data<LongUnsigned> err: " + err.Error())
	}
	return nil
}

func (l *LongUnsigned) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, l.Data)
	if err != nil {
		return nil, errors.New("encode data<LongUnsigned> err: " + err.Error())
	}
	return buf.Bytes(), nil
}

func (l *LongUnsigned) DataType() byte {
	return LongUnsignedIdent
}

func (l *LongUnsigned) Value() interface{} {
	return l.Data
}

/*-------------------------------------------*/

type Unsigned struct {
	Data uint8 `json:"data"`
}

func (u *Unsigned) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.LittleEndian, &u.Data); err != nil {
		return errors.New("decode data<Unsigned> err: " + err.Error())
	}
	return nil
}

func (u *Unsigned) encoder() ([]byte, error) {
	return []byte{u.Data}, nil
}

func (u *Unsigned) DataType() byte {
	return UnsignedIdent
}

func (u *Unsigned) Value() interface{} {
	return u.Data
}

/*-----------------------------------------------------------*/

type Long struct {
	Data int8 `json:"data"`
}

func (l *Long) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.LittleEndian, &l.Data); err != nil {
		return errors.New("decode data<Long> err" + err.Error())
	}
	return nil
}

func (l *Long) encoder() ([]byte, error) {
	return []byte{byte(l.Data)}, nil
}

func (l *Long) DataType() byte {
	return LongIdent
}

func (l *Long) Value() interface{} {
	return l.Data
}

/*------------------------------------------------------------------*/

type Integer struct {
	Data int8 `json:"data"`
}

func (i *Integer) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, &i.Data); err != nil {
		return errors.New("decode data<Integer> err:" + err.Error())
	}
	return nil
}

func (i *Integer) encoder() ([]byte, error) {
	return []byte{byte(i.Data)}, nil
}

func (i *Integer) DataType() byte {
	return IntegerIdent
}

func (i *Integer) Value() interface{} {
	return i.Data
}

/*-------------------------------------------------*/

type DoubleLongUnsigned struct {
	Data uint32 `json:"data"`
}

func (d *DoubleLongUnsigned) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, &d.Data); err != nil {
		return errors.New("decode data<DoubleLongUnsigned> err:" + err.Error())
	}
	return nil
}

func (d *DoubleLongUnsigned) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, d.Data)
	if err != nil {
		return nil, errors.New("encode data<DoubleLongUnsigned> err:" + err.Error())
	}
	return buf.Bytes(), nil
}

func (d *DoubleLongUnsigned) DataType() byte {
	return DoubleLongUnsignedIdent
}

func (d *DoubleLongUnsigned) Value() interface{} {
	return d.Data
}

/*------------------------------------------------------*/

type DoubleLong struct {
	Data int32 `json:"data"`
}

func (d *DoubleLong) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, &d.Data); err != nil {
		return errors.New("decode data<DoubleLong> err:" + err.Error())
	}
	return nil
}

func (d *DoubleLong) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, d.Data)
	if err != nil {
		return nil, errors.New("encode data<DoubleLong> err:" + err.Error())
	}
	return buf.Bytes(), nil
}

func (d *DoubleLong) DataType() byte {
	return DoubleLongIdent
}

func (d *DoubleLong) Value() interface{} {
	return d.Data
}

/*--------------------------*/

type Boolean struct {
	Data uint8 `json:"data"`
}

func (b *Boolean) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, &b.Data); err != nil {
		return errors.New("decode data<bool> err: " + err.Error())
	}
	return nil
}

func (b *Boolean) encoder() ([]byte, error) {
	return []byte{b.Data}, nil
}

func (b *Boolean) DataType() byte {
	return BooleanIdent
}

func (b *Boolean) Value() interface{} {
	return b.Data
}

/*----------------------------------------------------------*/

type Structure struct {
	DataArray []DataInter `json:"data_Array"`
}

func (s *Structure) decoder(buf *bytes.Reader) error {
	var arrayLen byte
	if err := binary.Read(buf, binary.LittleEndian, &arrayLen); err != nil {
		return errors.New("decode data<Array> length err:" + err.Error())
	}
	s.DataArray = make([]DataInter, arrayLen)
	for i := 0; i < int(arrayLen); i++ {
		dataType, err := buf.ReadByte()
		if err != nil {
			return err
		}
		s.DataArray[i] = dataTranslate(dataType)
		if s.DataArray[i] == nil {
			return errors.New("does not has this data type")
		}
		err = s.DataArray[i].decoder(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Structure) encoder() ([]byte, error) {
	length := len(s.DataArray)
	if length == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := make([]byte, 0)
	for _, datatype := range s.DataArray {
		dataArray, err := datatype.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, datatype.DataType())
		encodeArray = append(encodeArray, dataArray...)
	}
	return encodeArray, nil
}

func (s *Structure) DataType() byte {
	return StructureIdent
}

func (s *Structure) Value() interface{} {
	return s.DataArray
}

/*----------------------------------------------------------*/

type Array struct {
	DataArray []DataInter `json:"data_Array"`
}

func (a *Array) decoder(buf *bytes.Reader) error {
	var arrayLen byte
	if err := binary.Read(buf, binary.LittleEndian, &arrayLen); err != nil {
		return errors.New("decode data<Array> length err:" + err.Error())
	}
	a.DataArray = make([]DataInter, arrayLen)
	for i := 0; i < int(arrayLen); i++ {
		dataType, err := buf.ReadByte()
		if err != nil {
			return err
		}
		a.DataArray[i] = dataTranslate(dataType)
		if a.DataArray[i] == nil {
			return errors.New("does not has this data type")
		}
		err = a.DataArray[i].decoder(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Array) encoder() ([]byte, error) {
	length := len(a.DataArray)
	if length == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := make([]byte, 0)
	for _, datatype := range a.DataArray {
		dataArray, err := datatype.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, datatype.DataType())
		encodeArray = append(encodeArray, dataArray...)
	}
	return encodeArray, nil
}

func (a *Array) DataType() byte {
	return ArrayIdent
}

func (a *Array) Value() interface{} {
	return a.DataArray
}

/*----------------------------------------*/

type Null struct {
}

func (n *Null) decoder(_ *bytes.Reader) error {
	return nil
}

func (n *Null) encoder() ([]byte, error) {
	return nil, nil
}

func (n *Null) DataType() byte {
	return NullIdent
}

func (n *Null) Value() interface{} {
	return nil
}
