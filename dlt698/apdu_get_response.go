package dlt698

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var _ FrameRegion = (*ResultNormal)(nil)
var _ FrameRegion = (*GetResult)(nil)
var _ FrameRegion = (*ResultRecord)(nil)

var _ DataInter = (*DAR)(nil)
var _ DataInter = (*RecordRow)(nil)
var _ APDURegion = (*GetResponseNormal)(nil)
var _ APDURegion = (*GetResponseNormalList)(nil)
var _ APDURegion = (*GetResponseRecord)(nil)
var _ APDURegion = (*GetResponseRecordList)(nil)
var _ APDURegion = (*GetResponseNext)(nil)
var _ APDURegion = (*GetResponseMD5)(nil)

const (
	GetResponseNormalIdent     string = "8501"
	GetResponseNormalListIdent string = "8502"
	GetResponseRecordIdent     string = "8503"
	GetResponseRecordListIdent string = "8504"
	GetResponseNextIdent       string = "8505"
	GetResponseMD5Ident        string = "8506"
)

func init() {
	apduMap[GetResponseNormalIdent] = func() APDURegion {
		return new(GetResponseNormal)
	}
	apduMap[GetResponseNormalListIdent] = func() APDURegion {
		return new(GetResponseNormalList)
	}
	apduMap[GetResponseRecordIdent] = func() APDURegion {
		return new(GetResponseRecord)
	}
	apduMap[GetResponseRecordListIdent] = func() APDURegion {
		return new(GetResponseRecordList)
	}
	apduMap[GetResponseNextIdent] = func() APDURegion {
		return new(GetResponseNext)
	}
	apduMap[GetResponseMD5Ident] = func() APDURegion {
		return new(GetResponseMD5)
	}
}

type GetResponseNormal struct {
	ResultNormal *ResultNormal `json:"result_normal"`
}

func (g *GetResponseNormal) decoder(buf *bytes.Reader) error {
	g.ResultNormal = new(ResultNormal)
	return g.ResultNormal.decoder(buf)
}

func (g *GetResponseNormal) encoder() ([]byte, error) {
	return g.ResultNormal.encoder()
}

func (g *GetResponseNormal) APDUType() string {
	return GetResponseNormalIdent
}

func (g *GetResponseNormal) APDUMark() string {
	return "get_response_normal"
}

func (g *GetResponseNormal) hasFollowReport() bool {
	return true
}

func (g *GetResponseNormal) hasTimeTag() bool {
	return true
}

// ResultNormal 对象属性及结果
type ResultNormal struct {
	OAD       []byte     `json:"oad"`
	GetResult *GetResult `json:"get_result"`
}

func (r *ResultNormal) DataType() byte {
	return 0
}

func (r *ResultNormal) Value() interface{} {
	return r
}

func (r *ResultNormal) decoder(buf *bytes.Reader) error {
	r.OAD = make([]byte, 4)
	err := binary.Read(buf, binary.LittleEndian, &r.OAD)
	if err != nil {
		return errors.New("ResultNormal decode err :" + err.Error())
	}
	r.GetResult = new(GetResult)
	err = r.GetResult.decoder(buf)
	return err
}

func (r *ResultNormal) encoder() ([]byte, error) {
	getResultArray, err := r.GetResult.encoder()
	if err != nil {
		return nil, err
	}
	return append(r.OAD, getResultArray...), nil
}

type GetResult struct {
	Data DataInter `json:"data"`
}

func (g *GetResult) decoder(buf *bytes.Reader) error {
	dataType, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if dataType == 0 {
		g.Data = &DAR{}
		err = g.Data.decoder(buf)
		return err
	}
	//获取data类型
	dataType, err = buf.ReadByte()
	if err != nil {
		return err
	}
	g.Data = dataTranslate(dataType)
	if g.Data == nil {
		return errors.New("GetResult type not found！")
	}
	err = g.Data.decoder(buf)
	return err
}

func (g *GetResult) encoder() ([]byte, error) {
	dataArray, err := g.Data.encoder()
	if err != nil {
		return nil, err
	}
	switch g.Data.(type) {
	case *DAR:
		return append([]byte{0x00}, dataArray...), nil
	default:
		return append([]byte{0x01, g.Data.DataType()}, dataArray...), nil
	}
}

type DAR struct {
	Data byte `json:"value"`
}

func (D *DAR) DataType() byte {
	return 0
}

func (D *DAR) Value() interface{} {
	return D.Data
}

func (D *DAR) decoder(buf *bytes.Reader) error {
	return binary.Read(buf, binary.LittleEndian, &D.Data)
}

func (D *DAR) encoder() ([]byte, error) {
	return []byte{D.Data}, nil
}

/*-----------------------*/

type GetResponseNormalList struct {
	ResultNormals []*ResultNormal `json:"result_normals"`
}

func (g *GetResponseNormalList) decoder(buf *bytes.Reader) error {
	var rnLen byte
	err := binary.Read(buf, binary.LittleEndian, &rnLen)
	if err != nil {
		return errors.New("GetResponseNormalList decode err:" + err.Error())
	}
	if rnLen == 0 {
		return nil
	}
	g.ResultNormals = make([]*ResultNormal, rnLen)
	for i := 0; i < int(rnLen); i++ {
		g.ResultNormals[i] = &ResultNormal{}
		if err := g.ResultNormals[i].decoder(buf); err != nil {
			return err
		}
	}
	return nil
}

func (g *GetResponseNormalList) encoder() ([]byte, error) {
	if g.ResultNormals == nil || len(g.ResultNormals) == 0 {
		return nil, errors.New("getResponseNormalList ResultNormals must not be null！")
	}
	encodeArray := []byte{byte(len(g.ResultNormals))}
	for _, rn := range g.ResultNormals {
		rnArray, err := rn.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, rnArray...)
	}
	return encodeArray, nil
}

func (g *GetResponseNormalList) APDUType() string {
	return GetResponseNormalListIdent
}

func (g *GetResponseNormalList) APDUMark() string {
	return "get_response_normal_list"
}

func (g *GetResponseNormalList) hasFollowReport() bool {
	return true
}

func (g *GetResponseNormalList) hasTimeTag() bool {
	return true
}

/*-----------------------------------*/

type GetResponseRecord struct {
	ResultRecord *ResultRecord `json:"result_record"` // 一个记录型对象属性及结果
}

func (g *GetResponseRecord) decoder(buf *bytes.Reader) error {
	g.ResultRecord = &ResultRecord{}
	return g.ResultRecord.decoder(buf)
}

func (g *GetResponseRecord) encoder() ([]byte, error) {
	return g.ResultRecord.encoder()
}

func (g *GetResponseRecord) APDUType() string {
	return GetResponseRecordIdent
}

func (g *GetResponseRecord) APDUMark() string {
	return "get_response_record"
}

func (g *GetResponseRecord) hasFollowReport() bool {
	return true
}

func (g *GetResponseRecord) hasTimeTag() bool {
	return true
}

type ResultRecord struct {
	Oad  []byte    `json:"oad"`        //记录型对象属性描述符
	Rcsd *RCSD     `json:"rcsd"`       //记录的N列属性描述符
	Data DataInter `json:"record_row"` //M 条记录
}

func (r *ResultRecord) DataType() byte {
	return 0
}

func (r *ResultRecord) Value() interface{} {
	return r
}

func (r *ResultRecord) decoder(buf *bytes.Reader) error {
	r.Oad = make([]byte, 4)
	err := binary.Read(buf, binary.LittleEndian, r.Oad)
	if err != nil {
		return errors.New("ResultRecord decode err :" + err.Error())
	}
	r.Rcsd = &RCSD{}
	if err = r.Rcsd.decoder(buf); err != nil {
		return err
	}
	dataType, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if dataType == 0 {
		r.Data = &DAR{}
	} else {
		r.Data = &RecordRow{length: len(r.Rcsd.CSDs)}
	}
	return r.Data.decoder(buf)
}

func (r *ResultRecord) encoder() ([]byte, error) {
	rcsdArray, err := r.Rcsd.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray := append(r.Oad, rcsdArray...)
	dataArray, err := r.Data.encoder()
	if err != nil {
		return nil, err
	}
	switch r.Data.(type) {
	case *DAR:
		encodeArray = append(encodeArray, 0x00)
		return append(encodeArray, dataArray...), nil
	default:
		encodeArray = append(encodeArray, 0x01)
		return append(encodeArray, dataArray...), nil
	}
}

type RecordRow struct {
	length    int
	RecordRow []DataInter `json:"record_row"` //M 条记录
}

func (r *RecordRow) decoder(buf *bytes.Reader) error {
	r.RecordRow = make([]DataInter, r.length)
	_, err := buf.ReadByte()
	if err != nil {
		return err
	}
	for i := 0; i < r.length; i++ {
		dataType, err := buf.ReadByte()
		if err != nil {
			return err
		}
		r.RecordRow[i] = dataTranslate(dataType)
		err = r.RecordRow[i].decoder(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RecordRow) encoder() ([]byte, error) {
	encodeArray := []byte{1}
	for _, row := range r.RecordRow {
		encodeArray = append(encodeArray, row.DataType())
		rowArray, err := row.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, rowArray...)
	}
	return encodeArray, nil
}

func (r *RecordRow) DataType() byte {
	return 0
}

func (r *RecordRow) Value() interface{} {
	return r
}

/*--------------------------------*/

type GetResponseRecordList struct {
	ResultRecords []*ResultRecord `json:"result_records"`
}

func (g *GetResponseRecordList) decoder(buf *bytes.Reader) error {
	var length byte
	err := binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return errors.New("GetResponseRecordList decode err:" + err.Error())
	}
	if length == 0 {
		return nil
	}
	g.ResultRecords = make([]*ResultRecord, length)
	for i := 0; i < int(length); i++ {
		g.ResultRecords[i] = &ResultRecord{}
		if err := g.ResultRecords[i].decoder(buf); err != nil {
			return err
		}
	}
	return nil
}

func (g *GetResponseRecordList) encoder() ([]byte, error) {
	if g.ResultRecords == nil || len(g.ResultRecords) == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := []byte{byte(len(g.ResultRecords))}
	for _, g := range g.ResultRecords {
		recordArray, err := g.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, recordArray...)
	}
	return encodeArray, nil
}

func (g *GetResponseRecordList) APDUType() string {
	return GetResponseRecordListIdent
}

func (g *GetResponseRecordList) APDUMark() string {
	return "get_response_record_list"
}

func (g *GetResponseRecordList) hasFollowReport() bool {
	return true
}

func (g *GetResponseRecordList) hasTimeTag() bool {
	return false
}

/*--------------------*/

type GetResponseNext struct {
	EndFlag     byte        `json:"end_flag"` //末帧标志
	FrameNumber uint16      `json:"frame_number"`
	Dar         *DAR        `json:"dar"`
	Data        []DataInter `json:"data"`
}

func (g *GetResponseNext) decoder(buf *bytes.Reader) error {
	var err error
	if err = binary.Read(buf, binary.BigEndian, &g.EndFlag); err != nil {
		return errors.New("GetResponseNext decode err :" + err.Error())
	}
	if err = binary.Read(buf, binary.BigEndian, &g.FrameNumber); err != nil {
		return errors.New("GetResponseNext decode err :" + err.Error())
	}
	var dataType byte
	err = binary.Read(buf, binary.BigEndian, &dataType)
	if err != nil {
		return errors.New("GetResponseNext decode err :" + err.Error())
	}
	if dataType == 0 {
		g.Dar = &DAR{}
		err = g.Dar.decoder(buf)
		return err
	}
	var dataLen byte
	err = binary.Read(buf, binary.BigEndian, &dataLen)
	if err != nil {
		return err
	}
	g.Data = make([]DataInter, dataLen)
	switch dataType {
	case 1:
		for i := 0; i < int(dataLen); i++ {
			g.Data[i] = &ResultNormal{}
			err = g.Data[i].decoder(buf)
			if err != nil {
				return err
			}
		}
	case 2:
		for i := 0; i < int(dataLen); i++ {
			g.Data[i] = &ResultRecord{}
			err = g.Data[i].decoder(buf)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("GetResponseNext unknown data type")
	}
	return nil
}

func (g *GetResponseNext) encoder() ([]byte, error) {
	var err error
	buf := new(bytes.Buffer)
	if err = binary.Write(buf, binary.BigEndian, &g.EndFlag); err != nil {
		return nil, errors.New("GetResponseNext encode err:" + err.Error())
	}
	if err = binary.Write(buf, binary.BigEndian, &g.FrameNumber); err != nil {
		return nil, errors.New("GetResponseNext encode err:" + err.Error())
	}
	if g.Dar != nil {
		err = binary.Write(buf, binary.BigEndian, []byte{0x00, g.Dar.Data})
		return buf.Bytes(), err
	}
	switch g.Data[0].(type) {
	case *ResultNormal:
		err = binary.Write(buf, binary.BigEndian, []byte{0x01, byte(len(g.Data))})
	case *ResultRecord:
		err = binary.Write(buf, binary.BigEndian, []byte{0x02, byte(len(g.Data))})
	default:
		return nil, errors.New("GetResponseNext unknown data type")
	}
	if err != nil {
		return nil, err
	}
	for _, d := range g.Data {
		dArray, err := d.encoder()
		if err != nil {
			return nil, err
		}
		err = binary.Write(buf, binary.BigEndian, dArray)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (g *GetResponseNext) APDUType() string {
	return GetResponseNextIdent
}

func (g *GetResponseNext) APDUMark() string {
	return "get_response_next"
}

func (g *GetResponseNext) hasFollowReport() bool {
	return true
}

func (g *GetResponseNext) hasTimeTag() bool {
	return true
}

/*----------------------------------------*/

type GetResponseMD5 struct {
	Oad  []byte    `json:"oad"`
	Data DataInter `json:"data"`
}

func (g *GetResponseMD5) decoder(buf *bytes.Reader) error {
	g.Oad = make([]byte, 4)
	err := binary.Read(buf, binary.LittleEndian, &g.Oad)
	if err != nil {
		return err
	}
	dataType, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if dataType == 0 {
		g.Data = &DAR{}
	} else {
		g.Data = &OctetString{}
	}
	return g.Data.decoder(buf)
}

func (g *GetResponseMD5) encoder() ([]byte, error) {
	encodeArray := g.Oad[:]
	switch g.Data.(type) {
	case *DAR:
		encodeArray = append(encodeArray, byte(0x00))
	default:
		encodeArray = append(encodeArray, byte(0x01))
	}
	dataArray, err := g.Data.encoder()
	if err != nil {
		return nil, err
	}
	return append(encodeArray, dataArray...), nil
}

func (g *GetResponseMD5) APDUType() string {
	return GetResponseMD5Ident
}

func (g *GetResponseMD5) APDUMark() string {
	return "get_response_md5"
}

func (g *GetResponseMD5) hasFollowReport() bool {
	return true
}

func (g *GetResponseMD5) hasTimeTag() bool {
	return true
}
