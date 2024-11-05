package dlt698

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var _ FrameRegion = (*GetRecord)(nil)

var _ APDURegion = (*GetRequestNormal)(nil)
var _ APDURegion = (*GetRequestNormalList)(nil)
var _ APDURegion = (*GetRequestRecord)(nil)

const (
	GetRequestNormalIdent     string = "0501"
	GetRequestNormalListIdent string = "0502"
	GetRequestRecordIdent     string = "0503"
)

func init() {
	apduMap[GetRequestNormalIdent] = func() APDURegion {
		return new(GetRequestNormal)
	}
	apduMap[GetRequestNormalListIdent] = func() APDURegion {
		return new(GetRequestNormalList)
	}
	apduMap[GetRequestRecordIdent] = func() APDURegion {
		return new(GetRequestRecord)
	}
}

type GetRequestNormal struct {
	OAD []byte `json:"oad"` //一个对象属性描述符
}

func (g *GetRequestNormal) decoder(buf *bytes.Reader) error {
	g.OAD = make([]byte, 4)
	if _, err := buf.Read(g.OAD); err != nil {
		return errors.New("decode GetRequestNormal err:" + err.Error())
	}
	return nil
}

func (g *GetRequestNormal) encoder() ([]byte, error) {
	if len(g.OAD) != 4 {
		return nil, errors.New("oad length is not 4")
	}
	return g.OAD, nil
}

func (g *GetRequestNormal) APDUType() string {
	return GetRequestNormalIdent
}

func (g *GetRequestNormal) APDUMark() string {
	return "get_requestNormal"
}

func (g *GetRequestNormal) hasFollowReport() bool {
	return false
}

func (g *GetRequestNormal) hasTimeTag() bool {
	return true
}

/*-------------------GetRequestNormalList---------------------------*/

type GetRequestNormalList struct {
	OADs [][]byte `json:"oad_array"` //若干对象属性描述符
}

func (g *GetRequestNormalList) decoder(buf *bytes.Reader) error {
	var oadLen byte
	if err := binary.Read(buf, binary.BigEndian, &oadLen); err != nil {
		return errors.New("decode GetRequestNormalList's oad_array length err:" + err.Error())
	}
	g.OADs = make([][]byte, oadLen)
	for index := 0; index < int(oadLen); index++ {
		oad := make([]byte, 4)
		if _, err := buf.Read(oad); err != nil {
			return errors.New("decode GetRequestNormalList's oad_array err:" + err.Error())
		}
		g.OADs[index] = oad
	}
	return nil
}

func (g *GetRequestNormalList) encoder() ([]byte, error) {
	encodeArray := []byte{byte(len(g.OADs))}
	for _, oad := range g.OADs {
		encodeArray = append(encodeArray, oad...)
	}
	return encodeArray, nil
}

func (g *GetRequestNormalList) APDUType() string {
	return GetRequestNormalListIdent
}

func (g *GetRequestNormalList) APDUMark() string {
	return "get_request_normal_list"
}

func (g *GetRequestNormalList) hasFollowReport() bool {
	return false
}

func (g *GetRequestNormalList) hasTimeTag() bool {
	return true
}

/*----------------GetRequestRecord-----------*/

type GetRequestRecord struct {
	GetRecord *GetRecord `json:"get_record"` // 读取一个记录型对象属性
}

func (g *GetRequestRecord) decoder(buf *bytes.Reader) error {
	//TODO implement me
	panic("implement me")
}

func (g *GetRequestRecord) encoder() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GetRequestRecord) APDUType() string {
	return GetRequestRecordIdent
}

func (g *GetRequestRecord) APDUMark() string {
	return "get_request_record"
}

func (g *GetRequestRecord) hasFollowReport() bool {
	return false
}

func (g *GetRequestRecord) hasTimeTag() bool {
	return true
}

type GetRecord struct {
	OAD  []byte `json:"oad"`  //对象属性描述符
	Rsd  *RSD   `json:"rsd"`  //记录行选择描述符
	Rcsd *RCSD  `json:"rcsd"` // 记录列选择描述符
}

func (g *GetRecord) decoder(buf *bytes.Reader) error {
	//TODO implement me
	panic("implement me")
}

func (g *GetRecord) encoder() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
