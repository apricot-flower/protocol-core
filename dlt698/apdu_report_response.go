package dlt698

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var _ APDURegion = (*ReportResponseList)(nil)
var _ APDURegion = (*ReportResponseRecordList)(nil)
var _ APDURegion = (*ReportResponseTransData)(nil)

const (
	ReportResponseListIdent       string = "0801"
	ReportResponseRecordListIdent string = "0802"
	ReportResponseTransDataIdent  string = "0803"
)

func init() {
	apduMap[ReportResponseListIdent] = func() APDURegion {
		return new(ReportResponseList)
	}
	apduMap[ReportResponseRecordListIdent] = func() APDURegion {
		return new(ReportResponseRecordList)
	}
	apduMap[ReportResponseTransDataIdent] = func() APDURegion {
		return new(ReportResponseTransData)
	}
}

type ReportResponseList struct {
	Data [][]byte
}

func (r *ReportResponseList) decoder(buf *bytes.Reader) error {
	var length byte
	err := binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return errors.New("ReportResponseList decode err: " + err.Error())
	}
	if length == 0 {
		return nil
	}
	r.Data = make([][]byte, length)
	for i := 0; i < int(length); i++ {
		r.Data[i] = make([]byte, 4)
		if err := binary.Read(buf, binary.LittleEndian, r.Data[i]); err != nil {
			return err
		}
	}
	return err
}

func (r *ReportResponseList) encoder() ([]byte, error) {
	length := len(r.Data)
	if length == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := []byte{byte(length)}
	for _, o := range r.Data {
		encodeArray = append(encodeArray, o...)
	}
	return encodeArray, nil
}

func (r *ReportResponseList) APDUType() string {
	return ReportResponseListIdent
}

func (r *ReportResponseList) APDUMark() string {
	return "report_response_list"
}

func (r *ReportResponseList) hasFollowReport() bool {
	return false
}

func (r *ReportResponseList) hasTimeTag() bool {
	return true
}

/*-----------------------------------*/

type ReportResponseRecordList struct {
	ReportResponseList
}

func (r *ReportResponseRecordList) APDUType() string {
	return ReportResponseRecordListIdent
}

func (r *ReportResponseRecordList) APDUMark() string {
	return "report_response_record_list"
}

/*-------------------------------------*/

// ReportResponseTransData 响应上报透明数据
type ReportResponseTransData struct {
}

func (r *ReportResponseTransData) decoder(_ *bytes.Reader) error {
	return nil
}

func (r *ReportResponseTransData) encoder() ([]byte, error) {
	return nil, nil
}

func (r *ReportResponseTransData) APDUType() string {
	return ReportResponseTransDataIdent
}

func (r *ReportResponseTransData) APDUMark() string {
	return "report_response_trans_data"
}

func (r *ReportResponseTransData) hasFollowReport() bool {
	return false
}

func (r *ReportResponseTransData) hasTimeTag() bool {
	return true
}
