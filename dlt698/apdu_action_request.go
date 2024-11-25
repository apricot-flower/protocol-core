package dlt698

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var _ FrameRegion = (*ActionThenGetRequestNormalListItem)(nil)

var _ APDURegion = (*ActionRequestNormal)(nil)
var _ APDURegion = (*ActionRequestNormalList)(nil)
var _ APDURegion = (*ActionThenGetRequestNormalList)(nil)

const (
	ActionRequestNormalIdent            string = "0701"
	ActionRequestNormalListIdent        string = "0702"
	ActionThenGetRequestNormalListIdent string = "0703"
)

func init() {
	apduMap[ActionRequestNormalIdent] = func() APDURegion {
		return new(ActionRequestNormal)
	}
	apduMap[ActionRequestNormalListIdent] = func() APDURegion {
		return new(ActionRequestNormalList)
	}
	apduMap[ActionThenGetRequestNormalListIdent] = func() APDURegion {
		return new(ActionThenGetRequestNormalList)
	}
}

type ActionRequestNormal struct {
	OMD  *OMD      `json:"omd"`
	Data DataInter `json:"data"`
}

func (a *ActionRequestNormal) decoder(buf *bytes.Reader) error {
	a.OMD = &OMD{}
	if err := a.OMD.decoder(buf); err != nil {
		return err
	}
	var dataType byte
	err := binary.Read(buf, binary.BigEndian, &dataType)
	if err != nil {
		return errors.New("ActionRequestNormal decode err: " + err.Error())
	}
	a.Data = dataTranslate(dataType)
	return a.Data.decoder(buf)
}

func (a *ActionRequestNormal) encoder() ([]byte, error) {
	omdArray, err := a.OMD.encoder()
	if err != nil {
		return nil, err
	}
	if a.Data == nil {
		return append(omdArray, 0x00), nil
	}
	dataArray, err := a.Data.encoder()
	if err != nil {
		return nil, err
	}
	omdArray = append(omdArray, a.Data.DataType())
	return append(omdArray, dataArray...), nil
}

func (a *ActionRequestNormal) APDUType() string {
	return ActionRequestNormalIdent
}

func (a *ActionRequestNormal) APDUMark() string {
	return "action_request_normal"
}

func (a *ActionRequestNormal) hasFollowReport() bool {
	return false
}

func (a *ActionRequestNormal) hasTimeTag() bool {
	return true
}

/*--------------------------------*/

type ActionRequestNormalList struct {
	Data []*ActionRequestNormal `json:"data"`
}

func (a *ActionRequestNormalList) decoder(buf *bytes.Reader) error {
	var length byte
	err := binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return errors.New("ActionRequestNormalList decode err: " + err.Error())
	}
	if length == 0 {
		return nil
	}
	a.Data = make([]*ActionRequestNormal, length)
	for i := 0; i < int(length); i++ {
		a.Data[i] = &ActionRequestNormal{}
		if err := a.Data[i].decoder(buf); err != nil {
			return err
		}
	}
	return nil
}

func (a *ActionRequestNormalList) encoder() ([]byte, error) {
	length := byte(len(a.Data))
	if length == 0 {
		return []byte{length}, nil
	}
	encodeArray := []byte{length}
	for _, d := range a.Data {
		dArray, err := d.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, dArray...)
	}
	return encodeArray, nil
}

func (a *ActionRequestNormalList) APDUType() string {
	return ActionRequestNormalListIdent
}

func (a *ActionRequestNormalList) APDUMark() string {
	return "action_request_normal_list"
}

func (a *ActionRequestNormalList) hasFollowReport() bool {
	return false
}

func (a *ActionRequestNormalList) hasTimeTag() bool {
	return true
}

/*-----------------------------*/

type ActionThenGetRequestNormalList struct {
	Data []*ActionThenGetRequestNormalListItem `json:"data"`
}

func (a *ActionThenGetRequestNormalList) decoder(buf *bytes.Reader) error {
	var length byte
	err := binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return errors.New("ActionThenGetRequestNormalList decode err: " + err.Error())
	}
	if length == 0 {
		return nil
	}
	a.Data = make([]*ActionThenGetRequestNormalListItem, length)
	for i := 0; i < int(length); i++ {
		a.Data[i] = &ActionThenGetRequestNormalListItem{}
		if err := a.Data[i].decoder(buf); err != nil {
			return err
		}
	}
	return nil
}

func (a *ActionThenGetRequestNormalList) encoder() ([]byte, error) {
	length := byte(len(a.Data))
	if length == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := []byte{length}
	for _, x := range a.Data {
		xArray, err := x.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, xArray...)
	}
	return encodeArray, nil
}

func (a *ActionThenGetRequestNormalList) APDUType() string {
	return ActionThenGetRequestNormalListIdent
}

func (a *ActionThenGetRequestNormalList) APDUMark() string {
	return "action_then_get_request_normal_list"
}

func (a *ActionThenGetRequestNormalList) hasFollowReport() bool {
	return false
}

func (a *ActionThenGetRequestNormalList) hasTimeTag() bool {
	return true
}

type ActionThenGetRequestNormalListItem struct {
	Omd   *OMD      `json:"omd"`
	Data  DataInter `json:"data"`
	Oad   []byte    `json:"oad"`
	Daley uint8     `json:"daley"`
}

func (a *ActionThenGetRequestNormalListItem) decoder(buf *bytes.Reader) error {
	a.Omd = &OMD{}
	if err := a.Omd.decoder(buf); err != nil {
		return err
	}
	var dataType byte
	err := binary.Read(buf, binary.BigEndian, &dataType)
	if err != nil {
		return errors.New("ActionThenGetRequestNormalListItem decode err: " + err.Error())
	}
	a.Data = dataTranslate(dataType)
	if a.Data == nil {
		return errors.New("actionThenGetRequestNormalList data is nil")
	}
	if err = a.Data.decoder(buf); err != nil {
		return err
	}
	a.Oad = make([]byte, 4)
	err = binary.Read(buf, binary.LittleEndian, a.Oad)
	if err != nil {
		return err
	}
	return binary.Read(buf, binary.LittleEndian, &a.Daley)
}

func (a *ActionThenGetRequestNormalListItem) encoder() ([]byte, error) {
	omdArray, err := a.Omd.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray := omdArray[:]
	if a.Data == nil {
		encodeArray = append(encodeArray, 0x00)
	} else {
		encodeArray = append(encodeArray, a.Data.DataType())
		dataArray, err := a.Data.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, dataArray...)
	}
	encodeArray = append(encodeArray, a.Oad...)
	return append(encodeArray, a.Daley), nil
}
