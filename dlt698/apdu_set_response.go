package dlt698

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	SetResponseNormalIdent            string = "8601"
	SetResponseNormalListIdent        string = "8602"
	SetThenGetResponseNormalListIdent string = "8603"
)

func init() {
	apduMap[SetResponseNormalIdent] = func() APDURegion {
		return new(SetResponseNormal)
	}
	apduMap[SetResponseNormalListIdent] = func() APDURegion {
		return new(SetResponseNormalList)
	}
	apduMap[SetThenGetResponseNormalListIdent] = func() APDURegion {
		return new(SetThenGetResponseNormalList)
	}
}

var _ APDURegion = (*SetResponseNormal)(nil)
var _ APDURegion = (*SetResponseNormalList)(nil)
var _ APDURegion = (*SetThenGetResponseNormalList)(nil)

var _ FrameRegion = (*SetThenGetResponseNormalListItem)(nil)

type SetResponseNormal struct {
	Oad []byte `json:"oad"`
	Dar *DAR   `json:"dar"`
}

func (s *SetResponseNormal) decoder(buf *bytes.Reader) error {
	s.Oad = make([]byte, 4)
	err := binary.Read(buf, binary.LittleEndian, &s.Oad)
	if err != nil {
		return errors.New("SetResponseNormal decode err: " + err.Error())
	}
	s.Dar = &DAR{}
	return s.Dar.decoder(buf)
}

func (s *SetResponseNormal) encoder() ([]byte, error) {
	return append(s.Oad, s.Dar.Data), nil
}

func (s *SetResponseNormal) APDUType() string {
	return SetResponseNormalIdent
}

func (s *SetResponseNormal) APDUMark() string {
	return "set_response_normal"
}

func (s *SetResponseNormal) hasFollowReport() bool {
	return true
}

func (s *SetResponseNormal) hasTimeTag() bool {
	return true
}

/*----------------------------------*/

type SetResponseNormalList struct {
	Data []*SetResponseNormal `json:"data"`
}

func (s *SetResponseNormalList) decoder(buf *bytes.Reader) error {
	var length byte
	err := binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return errors.New("SetResponseNormal decode err: " + err.Error())
	}
	if length == 0 {
		return nil
	}
	s.Data = make([]*SetResponseNormal, length)
	for i := 0; i < int(length); i++ {
		s.Data[i] = new(SetResponseNormal)
		err = s.Data[i].decoder(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SetResponseNormalList) encoder() ([]byte, error) {
	if len(s.Data) == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := []byte{byte(len(s.Data))}
	for _, d := range s.Data {
		dArray, err := d.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, dArray...)
	}
	return encodeArray, nil
}

func (s *SetResponseNormalList) APDUType() string {
	return SetResponseNormalListIdent
}

func (s *SetResponseNormalList) APDUMark() string {
	return "set_response_normal_list"
}

func (s *SetResponseNormalList) hasFollowReport() bool {
	return true
}

func (s *SetResponseNormalList) hasTimeTag() bool {
	return true
}

/*-------------------------------------------------*/

type SetThenGetResponseNormalList struct {
	Data []*SetThenGetResponseNormalListItem `json:"data"`
}

func (s *SetThenGetResponseNormalList) decoder(buf *bytes.Reader) error {
	var length byte
	err := binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return errors.New("SetThenGetResponseNormal decode err: " + err.Error())
	}
	if length == 0 {
		return nil
	}
	s.Data = make([]*SetThenGetResponseNormalListItem, length)
	for i := 0; i < int(length); i++ {
		s.Data[i] = &SetThenGetResponseNormalListItem{}
		err = s.Data[i].decoder(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SetThenGetResponseNormalList) encoder() ([]byte, error) {
	length := byte(len(s.Data))
	if length == 0 {
		return []byte{length}, nil
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

func (s *SetThenGetResponseNormalList) APDUType() string {
	return SetThenGetResponseNormalListIdent
}

func (s *SetThenGetResponseNormalList) APDUMark() string {
	return "set_then_get_response_normal_list"
}

func (s *SetThenGetResponseNormalList) hasFollowReport() bool {
	return true
}

func (s *SetThenGetResponseNormalList) hasTimeTag() bool {
	return true
}

type SetThenGetResponseNormalListItem struct {
	Oad          []byte        `json:"oad"`
	Dar          *DAR          `json:"dar"`
	ResultNormal *ResultNormal `json:"result_normal"`
}

func (s *SetThenGetResponseNormalListItem) decoder(buf *bytes.Reader) error {
	s.Oad = make([]byte, 4)
	err := binary.Read(buf, binary.LittleEndian, &s.Oad)
	if err != nil {
		return errors.New("SetThenGetResponseNormalListItem decode err: " + err.Error())
	}
	s.Dar = new(DAR)
	err = s.Dar.decoder(buf)
	if err != nil {
		return err
	}
	s.ResultNormal = new(ResultNormal)
	err = s.ResultNormal.decoder(buf)
	return err
}

func (s *SetThenGetResponseNormalListItem) encoder() ([]byte, error) {
	encodeArray := append(s.Oad, s.Dar.Data)
	dataArray, err := s.ResultNormal.encoder()
	if err != nil {
		return nil, err
	}
	return append(encodeArray, dataArray...), nil
}
