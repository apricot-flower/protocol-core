package dlt698

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var _ APDURegion = (*SetRequestNormal)(nil)
var _ APDURegion = (*SetRequestNormalList)(nil)
var _ APDURegion = (*SetThenGetRequestNormalList)(nil)

var _ FrameRegion = (*SetThenGetRequestItem)(nil)

const (
	SetRequestNormalIdent            string = "0601"
	SetRequestNormalListIdent        string = "0602"
	SetThenGetRequestNormalListIdent string = "0603"
)

func init() {
	apduMap[SetRequestNormalIdent] = func() APDURegion {
		return new(SetRequestNormal)
	}
	apduMap[SetRequestNormalListIdent] = func() APDURegion {
		return new(SetRequestNormalList)
	}
	apduMap[SetThenGetRequestNormalListIdent] = func() APDURegion {
		return new(SetThenGetRequestNormalList)
	}
}

type SetRequestNormal struct {
	Oad  []byte    `json:"oad"`
	Data DataInter `json:"data"`
}

func (s *SetRequestNormal) decoder(buf *bytes.Reader) error {
	s.Oad = make([]byte, 4)
	err := binary.Read(buf, binary.LittleEndian, &s.Oad)
	if err != nil {
		return errors.New("SetRequestNormal decode err: " + err.Error())
	}
	var dataFlg byte
	err = binary.Read(buf, binary.LittleEndian, &dataFlg)
	if err != nil {
		return errors.New("SetRequestNormal decode err: " + err.Error())
	}
	s.Data = dataTranslate(dataFlg)
	if s.Data == nil {
		return errors.New("SetRequestNormal data not available")
	}
	return s.Data.decoder(buf)
}

func (s *SetRequestNormal) encoder() ([]byte, error) {
	encodeArray := append(s.Oad, s.Data.DataType())
	dataArray, err := s.Data.encoder()
	if err != nil {
		return nil, err
	}
	return append(encodeArray, dataArray...), nil
}

func (s *SetRequestNormal) APDUType() string {
	return SetRequestNormalIdent
}

func (s *SetRequestNormal) APDUMark() string {
	return "set_request_normal"
}

func (s *SetRequestNormal) hasFollowReport() bool {
	return false
}

func (s *SetRequestNormal) hasTimeTag() bool {
	return true
}

/*----------------------------------*/

type SetRequestNormalList struct {
	Data []*SetRequestNormal `json:"data"`
}

func (s *SetRequestNormalList) decoder(buf *bytes.Reader) error {
	var length byte
	err := binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return errors.New("SetRequestNormal decode err: " + err.Error())
	}
	s.Data = make([]*SetRequestNormal, length)
	for i := 0; i < int(length); i++ {
		s.Data[i] = &SetRequestNormal{}
		err = s.Data[i].decoder(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SetRequestNormalList) encoder() ([]byte, error) {
	length := byte(len(s.Data))
	if length == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := []byte{length}
	for _, data := range s.Data {
		dataArray, err := data.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, dataArray...)
	}
	return encodeArray, nil
}

func (s *SetRequestNormalList) APDUType() string {
	return SetRequestNormalListIdent
}

func (s *SetRequestNormalList) APDUMark() string {
	return "set_request_normal_list"
}

func (s *SetRequestNormalList) hasFollowReport() bool {
	return false
}

func (s *SetRequestNormalList) hasTimeTag() bool {
	return true
}

/*------------------------------------------*/

type SetThenGetRequestNormalList struct {
	Data []*SetThenGetRequestItem `json:"data"`
}

func (s *SetThenGetRequestNormalList) decoder(buf *bytes.Reader) error {
	var length byte
	err := binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return errors.New("SetThenGetRequestNormal decode err: " + err.Error())
	}
	if length == 0 {
		return nil
	}
	s.Data = make([]*SetThenGetRequestItem, length)
	for i := 0; i < int(length); i++ {
		s.Data[i] = &SetThenGetRequestItem{}
		err = s.Data[i].decoder(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SetThenGetRequestNormalList) encoder() ([]byte, error) {
	length := byte(len(s.Data))
	if length == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := []byte{length}
	for _, d := range s.Data {
		dArray, err := d.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, dArray...)
	}
	return encodeArray, nil
}

func (s *SetThenGetRequestNormalList) APDUType() string {
	return SetThenGetRequestNormalListIdent
}

func (s *SetThenGetRequestNormalList) APDUMark() string {
	return "set_then_get_request_normal_list"
}

func (s *SetThenGetRequestNormalList) hasFollowReport() bool {
	return false
}

func (s *SetThenGetRequestNormalList) hasTimeTag() bool {
	return true
}

type SetThenGetRequestItem struct {
	SetOad  []byte    `json:"set_oad"`
	Data    DataInter `json:"data"`
	ReadOad []byte    `json:"read_oad"`
	Delay   uint8     `json:"delay"`
}

func (s *SetThenGetRequestItem) decoder(buf *bytes.Reader) error {
	s.SetOad = make([]byte, 4)
	if _, err := buf.Read(s.SetOad); err != nil {
		return errors.New("SetThenGetRequestItem decode err: " + err.Error())
	}
	var dataType byte
	if err := binary.Read(buf, binary.LittleEndian, &dataType); err != nil {
		return errors.New("SetThenGetRequestItem decode err: " + err.Error())
	}
	s.Data = dataTranslate(dataType)
	err := s.Data.decoder(buf)
	if err != nil {
		return err
	}
	s.ReadOad = make([]byte, 4)
	if _, err = buf.Read(s.ReadOad); err != nil {
		return err
	}
	return binary.Read(buf, binary.LittleEndian, &s.Delay)
}

func (s *SetThenGetRequestItem) encoder() ([]byte, error) {
	encodeArray := s.SetOad[:]
	dataArray, err := s.Data.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, s.Data.DataType())
	encodeArray = append(encodeArray, dataArray...)
	encodeArray = append(encodeArray, s.ReadOad...)
	return append(encodeArray, s.Delay), nil
}
