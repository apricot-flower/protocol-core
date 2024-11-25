package dlt698

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var _ APDURegion = (*ReportNotificationList)(nil)
var _ APDURegion = (*ReportNotificationRecordList)(nil)
var _ APDURegion = (*ReportNotificationTransData)(nil)

const (
	ReportNotificationListIdent       string = "8801"
	ReportNotificationRecordListIdent string = "8802"
	ReportNotificationTransDataIdent  string = "8803"
)

func init() {
	apduMap[ReportNotificationListIdent] = func() APDURegion {
		return new(ReportNotificationList)
	}
	apduMap[ReportNotificationRecordListIdent] = func() APDURegion {
		return new(ReportNotificationRecordList)
	}
	apduMap[ReportNotificationTransDataIdent] = func() APDURegion {
		return new(ReportNotificationTransData)
	}
}

type ReportNotificationList struct {
	Data []*ResultNormal
}

func (r *ReportNotificationList) decoder(buf *bytes.Reader) error {
	var length byte
	err := binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return errors.New("ReportNotificationList decode err: " + err.Error())
	}
	if length == 0 {
		return nil
	}
	r.Data = make([]*ResultNormal, length)
	for i := 0; i < int(length); i++ {
		r.Data[i] = &ResultNormal{}
		err = r.Data[i].decoder(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReportNotificationList) encoder() ([]byte, error) {
	length := len(r.Data)
	if length == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := []byte{byte(length)}
	for _, result := range r.Data {
		resultArr, err := result.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, resultArr...)
	}
	return encodeArray, nil
}

func (r *ReportNotificationList) APDUType() string {
	return ReportNotificationListIdent
}

func (r *ReportNotificationList) APDUMark() string {
	return "report_notification_list"
}

func (r *ReportNotificationList) hasFollowReport() bool {
	return true
}

func (r *ReportNotificationList) hasTimeTag() bool {
	return true
}

/*-----------------------------------------------*/

type ReportNotificationRecordList struct {
	Data []*ResultRecord
}

func (r *ReportNotificationRecordList) decoder(buf *bytes.Reader) error {
	var length byte
	err := binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return errors.New("ReportNotificationRecordList decode err: " + err.Error())
	}
	r.Data = make([]*ResultRecord, length)
	for i := 0; i < int(length); i++ {
		r.Data[i] = &ResultRecord{}
		err = r.Data[i].decoder(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReportNotificationRecordList) encoder() ([]byte, error) {
	length := len(r.Data)
	if length == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := []byte{byte(length)}
	for _, record := range r.Data {
		recordArray, err := record.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, recordArray...)
	}
	return encodeArray, nil
}

func (r *ReportNotificationRecordList) APDUType() string {
	return ReportNotificationRecordListIdent
}

func (r *ReportNotificationRecordList) APDUMark() string {
	return "report_notification_record_list"
}

func (r *ReportNotificationRecordList) hasFollowReport() bool {
	return true
}

func (r *ReportNotificationRecordList) hasTimeTag() bool {
	return true
}

/*------------------------------------*/

type ReportNotificationTransData struct {
	Oad  []byte
	Data *OctetString
}

func (r *ReportNotificationTransData) decoder(buf *bytes.Reader) error {
	r.Oad = make([]byte, 4)
	if err := binary.Read(buf, binary.LittleEndian, &r.Oad); err != nil {
		return errors.New("ReportNotificationTransData decode err: " + err.Error())
	}
	hasData, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if hasData == 0 {
		return nil
	}
	r.Data = &OctetString{}
	return r.Data.decoder(buf)
}

func (r *ReportNotificationTransData) encoder() ([]byte, error) {
	if r.Data == nil {
		return append(r.Oad, 0x00), nil
	}
	dataArray, err := r.Data.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray := append(r.Oad, 0x01)
	return append(encodeArray, dataArray...), err
}

func (r *ReportNotificationTransData) APDUType() string {
	return ReportNotificationTransDataIdent
}

func (r *ReportNotificationTransData) APDUMark() string {
	return "report_notification_trans_data"
}

func (r *ReportNotificationTransData) hasFollowReport() bool {
	return true
}

func (r *ReportNotificationTransData) hasTimeTag() bool {
	return true
}
