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

var _ DataInter = (*VisibleString)(nil)
var _ DataInter = (*OctetString)(nil)
var _ DataInter = (*BitString)(nil)
var _ DataInter = (*UTF8String)(nil)

const (
	BitStringIdent     byte = 4
	VisibleStringIdent byte = 10
	OctetStringIdent   byte = 9
	UTF8StringIdent    byte = 12
)

func init() {
	dataMap[BitStringIdent] = func() DataInter {
		return &BitString{}
	}
	dataMap[VisibleStringIdent] = func() DataInter {
		return &VisibleString{}
	}
	dataMap[OctetStringIdent] = func() DataInter {
		return &OctetString{}
	}
	dataMap[UTF8StringIdent] = func() DataInter {
		return &UTF8String{}
	}
}

type UTF8String struct {
	Data string `json:"data"`
}

func (U *UTF8String) decoder(buf *bytes.Reader) error {
	var strLen byte
	if err := binary.Read(buf, binary.BigEndian, &strLen); err != nil {
		return errors.New("decode data<UTF8String> len err: " + err.Error())
	}
	arr := make([]byte, int(strLen))
	if err := binary.Read(buf, binary.BigEndian, arr); err != nil {
		U.Data = string(arr)
	}
	return nil
}

func (U *UTF8String) encoder() ([]byte, error) {
	if U.Data == "" {
		return []byte{0x00}, nil
	}
	arr, err := hex.DecodeString(U.Data)
	if err != nil {
		return nil, errors.New("decode data<UTF8String> err: " + err.Error())
	}
	return append([]byte{byte(len(U.Data))}, arr...), nil
}

func (U *UTF8String) DataType() byte {
	return UTF8StringIdent
}

func (U *UTF8String) Value() interface{} {
	return U.Data
}

/*-----------------------------------------------------------------*/

type BitString struct {
	Size byte
	Data string `json:"data"`
}

func (b *BitString) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, &b.Size); err != nil {
		return errors.New("decode data<BitString> length err: " + err.Error())
	}
	arr := make([]byte, int(b.Size/8))
	err := binary.Read(buf, binary.BigEndian, arr)
	if err != nil {
		return errors.New("decode data<BitString> err: " + err.Error())
	}
	for _, value := range arr {
		bs := fmt.Sprintf("%08b", value)
		b.Data += bs
	}
	return nil
}

func (b *BitString) encoder() ([]byte, error) {
	if b.Size == 0 {
		return []byte{0x00}, nil
	}
	b.adjustStringLength()
	arr, err := b.binaryStringToBytes()
	if err != nil {
		return nil, err
	}
	return append([]byte{byte(len(b.Data))}, arr...), nil
}

func (b *BitString) binaryStringToBytes() ([]byte, error) {
	var arr []byte
	length := len(b.Data)
	if length%8 != 0 {
		padding := 8 - (length % 8)
		b.Data = strings.Repeat("0", padding) + b.Data
	}
	for i := 0; i < len(b.Data); i += 8 {
		chunk := b.Data[i : i+8]
		b, err := strconv.ParseUint(chunk, 2, 8)
		if err != nil {
			return nil, err
		}
		arr = append(arr, byte(b))
	}

	return arr, nil
}

func (b *BitString) adjustStringLength() {
	if byte(len(b.Data)) > b.Size {
		// 如果字符串比目标长度长，就截断它
		b.Data = b.Data[:b.Size]
	} else if byte(len(b.Data)) < b.Size {
		// 如果字符串比目标长度短，就在后面添加 "0" 直到达到目标长度
		padding := b.Size - byte(len(b.Data))
		b.Data += strings.Repeat("0", int(padding))
	}
}

func (b *BitString) DataType() byte {
	return BitStringIdent
}

func (b *BitString) Value() interface{} {
	return b.Data
}

/*-----------------------------------------*/

type VisibleString struct {
	Data string `json:"data"`
}

func (v *VisibleString) decoder(buf *bytes.Reader) error {
	var length byte
	if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
		return errors.New("decode VisibleString's length err: " + err.Error())
	}
	if length == 0 {
		return nil
	}
	asciiArray := make([]byte, length)
	if _, err := buf.Read(asciiArray); err != nil {
		return errors.New("decode VisibleString's ascii array err: " + err.Error())
	}
	v.Data = v.asciiSliceToString(asciiArray)
	return nil
}

func (v *VisibleString) DecodeByLen(buf *bytes.Reader, length int) error {
	asciiArray := make([]byte, length)
	if _, err := buf.Read(asciiArray); err != nil {
		return errors.New("decode VisibleString's ascii array err: " + err.Error())
	}
	v.Data = v.asciiSliceToString(asciiArray)
	return nil
}

func (v *VisibleString) asciiSliceToString(array []byte) string {
	var sb strings.Builder
	for _, code := range array {
		sb.WriteRune(rune(code))
	}
	return sb.String()
}

func (v *VisibleString) encoder() ([]byte, error) {
	if v.Data == "" {
		return []byte{0x00}, nil
	}
	result := v.stringToAscii(v.Data)
	return append([]byte{byte(len(result))}, result...), nil
}

func (v *VisibleString) stringToAscii(s string) []byte {
	var asciiValues []byte
	for _, char := range s {
		asciiValues = append(asciiValues, byte(char))
	}
	return asciiValues
}

func (v *VisibleString) DataType() byte {
	return VisibleStringIdent
}

func (v *VisibleString) Value() interface{} {
	return v.Data
}

/*-----------------------------OctetString--------------------------*/

type OctetString struct {
	Data string `json:"data"`
}

func (o *OctetString) decoder(buf *bytes.Reader) error {
	var length byte
	if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
		return errors.New("decode OctetString's length err: " + err.Error())
	}
	if length == 0 {
		return nil
	}
	array := make([]byte, length)
	if err := binary.Read(buf, binary.LittleEndian, &array); err != nil {
		return errors.New("decode OctetString's array err: " + err.Error())
	}
	o.Data = hex.EncodeToString(array)
	return nil
}

func (o *OctetString) encoder() ([]byte, error) {
	if o.Data == "" {
		return []byte{0x00}, nil
	}
	if len(o.Data)%2 != 0 {
		return nil, errors.New("encode OctetString's data length should be even")
	}
	array, err := hex.DecodeString(o.Data)
	if err != nil {
		return nil, errors.New("encode OctetString's array err: " + err.Error())
	}
	return append([]byte{byte(len(array))}, array...), nil
}

func (o *OctetString) DataType() byte {
	return OctetStringIdent
}

func (o *OctetString) Value() interface{} {
	return o.Data
}
