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
var _ APDURegion = (*GetRequestRecordList)(nil)
var _ APDURegion = (*GetRequestNext)(nil)
var _ APDURegion = (*GetRequestMD5)(nil)

const (
	GetRequestNormalIdent     string = "0501"
	GetRequestNormalListIdent string = "0502"
	GetRequestRecordIdent     string = "0503"
	GetRequestRecordListIdent string = "0504"
	GetRequestNextIdent       string = "0505"
	GetRequestMD5Ident        string = "0506"
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
	apduMap[GetRequestRecordListIdent] = func() APDURegion {
		return new(GetRequestRecordList)
	}
	apduMap[GetRequestNextIdent] = func() APDURegion {
		return new(GetRequestNext)
	}
	apduMap[GetRequestMD5Ident] = func() APDURegion {
		return new(GetRequestMD5)
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
	g.GetRecord = new(GetRecord)
	err := g.GetRecord.decoder(buf)
	return err
}

func (g *GetRequestRecord) encoder() ([]byte, error) {
	return g.GetRecord.encoder()
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
	g.OAD = make([]byte, 4)
	err := binary.Read(buf, binary.LittleEndian, &g.OAD)
	if err != nil {
		return err
	}
	g.Rsd = new(RSD)
	err = g.Rsd.decoder(buf)
	if err != nil {
		return err
	}
	g.Rcsd = new(RCSD)
	err = g.Rcsd.decoder(buf)
	return err
}

func (g *GetRecord) encoder() ([]byte, error) {
	rsdArray, err := g.Rsd.encoder()
	if err != nil {
		return nil, err
	}
	rcsdArray, err := g.Rcsd.encoder()
	if err != nil {
		return nil, err
	}
	rsdArray = append(g.OAD, rsdArray...)
	return append(rsdArray, rcsdArray...), nil
}

/*-----------------------------------------*/

type GetRequestRecordList struct {
	GetRecords []*GetRecord `json:"get_records"` // 读取一个记录型对象属性
}

func (g *GetRequestRecordList) decoder(buf *bytes.Reader) error {
	var length byte
	err := binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return errors.New("decode GetRequestRecordList's get_records length err:" + err.Error())
	}
	g.GetRecords = make([]*GetRecord, length)
	for i := 0; i < int(length); i++ {
		g.GetRecords[i] = &GetRecord{}
		if err := g.GetRecords[i].decoder(buf); err != nil {
			return err
		}
	}
	return nil
}

func (g *GetRequestRecordList) encoder() ([]byte, error) {
	if len(g.GetRecords) == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := []byte{byte(len(g.GetRecords))}
	for _, record := range g.GetRecords {
		arr, err := record.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, arr...)
	}
	return encodeArray, nil
}

func (g *GetRequestRecordList) APDUType() string {
	return GetRequestRecordListIdent
}

func (g *GetRequestRecordList) APDUMark() string {
	return "get_request_record_list"
}

func (g GetRequestRecordList) hasFollowReport() bool {
	return false
}

func (g GetRequestRecordList) hasTimeTag() bool {
	return true
}

/*---------------------------*/

type GetRequestNext struct {
	LastId uint16 `json:"last_id"`
}

func (g *GetRequestNext) decoder(buf *bytes.Reader) error {
	return binary.Read(buf, binary.BigEndian, &g.LastId)
}

func (g *GetRequestNext) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, g.LastId)
	return buf.Bytes(), err
}

func (g *GetRequestNext) APDUType() string {
	return GetRequestNextIdent
}

func (g *GetRequestNext) APDUMark() string {
	return "get_request_next"
}

func (g *GetRequestNext) hasFollowReport() bool {
	return false
}

func (g *GetRequestNext) hasTimeTag() bool {
	return true
}

/*---------------------*/

type GetRequestMD5 struct {
	OAD []byte `json:"oad"`
}

func (g *GetRequestMD5) decoder(buf *bytes.Reader) error {
	g.OAD = make([]byte, 4)
	return binary.Read(buf, binary.LittleEndian, &g.OAD)
}

func (g *GetRequestMD5) encoder() ([]byte, error) {
	return g.OAD, nil
}

func (g *GetRequestMD5) APDUType() string {
	return GetRequestMD5Ident
}

func (g *GetRequestMD5) APDUMark() string {
	return "get_request_md5"
}

func (g *GetRequestMD5) hasFollowReport() bool {
	return false
}

func (g *GetRequestMD5) hasTimeTag() bool {
	return true
}
