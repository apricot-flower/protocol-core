package dlt1376_1

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

var _ Dlt13761DataInter = (*RelayStationCommandDownF1)(nil)
var _ Dlt13761DataInter = (*RelayStationCommandUpF1)(nil)
var _ Dlt13761DataInter = (*WorkStatusSwitchRecordResponse)(nil)

// RelayStationCommandDownF1 中继站命令下行F1
type RelayStationCommandDownF1 struct {
	SwitchControl byte // 值班机/备用机切换控制
	AllowFlag     byte //值班机中继转发允许标志
}

func (r *RelayStationCommandDownF1) Decode(buf ...byte) error {
	value := buf[0]
	r.SwitchControl = value & 0b11
	r.AllowFlag = (value & 0b00001100) >> 2
	return nil
}

func (r *RelayStationCommandDownF1) Encode() ([]byte, error) {
	if r.SwitchControl > 3 {
		return nil, fmt.Errorf("SwitchControl out of range")
	}
	if r.AllowFlag > 3 {
		return nil, fmt.Errorf("AllowFlag out of range")
	}
	switchControlStr := fmt.Sprintf("%02b", r.SwitchControl&0b11)
	allowFlagStr := fmt.Sprintf("%02b", r.AllowFlag&0b11)
	bitsStr := "0000" + allowFlagStr + switchControlStr
	value, err := strconv.ParseInt(bitsStr, 2, 8)
	if err != nil {
		return nil, err
	}
	return []byte{byte(value)}, err
}

// RelayStationCommandUpF1 中继站命令上行F1
type RelayStationCommandUpF1 struct {
	ADeviceStatus    byte //A机状态
	ADeviceType      byte //0-A机为备用机，1-A机为值班机
	ADeviceRelayFlag byte //0-A机禁止中继转发，1-A机允许中继转发
	BDeviceStatus    byte //B机状态
	BDeviceType      byte //0-B机为备份机，1-B机为值班机
	BDeviceRelayFlag byte //0-B机禁止中继转发，1-B机允许中继转发
}

func (r *RelayStationCommandUpF1) Decode(buf ...byte) error {
	value := buf[0]
	//解析A机
	r.ADeviceStatus = value & 0b11
	r.ADeviceType = (value >> 2) & 0b00000001
	r.ADeviceRelayFlag = (value >> 3) & 0b00000001
	//解析B机
	r.BDeviceStatus = (value >> 4) & 0b00000011
	r.BDeviceType = (value >> 6) & 0b00000001
	r.BDeviceRelayFlag = (value >> 7) & 0b00000001
	return nil
}

func (r *RelayStationCommandUpF1) Encode() ([]byte, error) {
	if r.ADeviceStatus > 3 {
		return nil, errors.New("ADeviceStatus out of range, must <= 3")
	}
	if r.ADeviceType > 1 {
		return nil, errors.New("ADeviceType out of range, must <= 1")
	}
	if r.ADeviceRelayFlag > 1 {
		return nil, errors.New("ADeviceRelayFlag out of range, must <= 1")
	}
	if r.BDeviceStatus > 3 {
		return nil, errors.New("BDeviceStatus out of range, must <= 3")
	}
	if r.BDeviceType > 1 {
		return nil, errors.New("BDeviceType out of range, must <= 1")
	}
	if r.BDeviceRelayFlag > 1 {
		return nil, errors.New("BDeviceRelayFlag out of range, must <= 1")
	}
	bitsStr := strconv.Itoa(int(r.BDeviceRelayFlag)) + strconv.Itoa(int(r.BDeviceType)) + fmt.Sprintf("%02b", r.BDeviceStatus&0b11) + strconv.Itoa(int(r.ADeviceRelayFlag)) + strconv.Itoa(int(r.ADeviceType)) + fmt.Sprintf("%02b", r.ADeviceStatus&0b11)
	value, err := strconv.ParseInt(bitsStr, 2, 8)
	if err != nil {
		return nil, err
	}
	return []byte{byte(value)}, err
}

// SwitchRecord 切换记录
type SwitchRecord struct {
	Time             *TimeA15                 //切换时间
	BeforeWorkStatus *RelayStationCommandUpF1 //切换前工作状态
	AfterWorkStatus  *RelayStationCommandUpF1 //切换后工作状态
}

// WorkStatusSwitchRecordResponse 中继站工作状态切换记录应答
type WorkStatusSwitchRecordResponse struct {
	Data map[int]*SwitchRecord
}

func (w *WorkStatusSwitchRecordResponse) Decode(buf ...byte) error {
	buff := bytes.NewReader(buf)
	w.Data = make(map[int]*SwitchRecord)
	for index := 1; index < 11; index++ {
		//time
		timeArr := make([]byte, 5)
		if err := binary.Read(buff, binary.LittleEndian, &timeArr); err != nil {
			return err
		}
		status1, err := buff.ReadByte()
		if err != nil {
			return err
		}
		status2, err := buff.ReadByte()
		if err != nil {
			return err
		}
		if !checkData(timeArr...) && !checkData(status1, status2) {
			w.Data[index] = nil
			continue
		}
		w.Data[index] = &SwitchRecord{}
		if checkData(timeArr...) {
			err = w.Data[index].Time.Decode(timeArr...)
			if err != nil {
				return err
			}
		}
		if checkData(status1) {
			err = w.Data[index].BeforeWorkStatus.Decode(status1)
			if err != nil {
				return err
			}
		}
		if checkData(status2) {
			err = w.Data[index].AfterWorkStatus.Decode(status2)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *WorkStatusSwitchRecordResponse) Encode() ([]byte, error) {

	return nil, nil
}

var _ Afn = (*RelayStationCommand)(nil)

// RelayStationCommand 中继站命令
type RelayStationCommand struct {
	dir  string
	data []*Dlt13761Data
}

func (r *RelayStationCommand) Decode(buf *bytes.Reader) error {
	for buf.Len() != 0 {
		afnBytes := make([]byte, 4)
		_, err := buf.Read(afnBytes)
		if err != nil {
			return err
		}
		pn, fn := analyzeUnit(afnBytes)
		err = r.expand(pn, fn, buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RelayStationCommand) expand(pn []uint64, fn []uint64, buf *bytes.Reader) error {
	for _, p := range pn {
		if p > 0 {
			return errors.New("afn.RelayStationCommand pn must be 0")
		}
		var err error
		for _, f := range fn {
			if f == 1 {
				//中继站工作状态控制
				err = r.statusControl(buf)
			} else if f == 2 {
				err = r.workStatusResponse(buf)
			} else if f == 3 {
				err = r.workStatusSwitchRecordResponse(buf)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 中继站工作状态切换记录应答
func (r *RelayStationCommand) workStatusSwitchRecordResponse(buf *bytes.Reader) error {
	if r.dir == "0" {
		//下行

	}
	return nil
}

// 中继站工作状态应答
func (r *RelayStationCommand) workStatusResponse(buf *bytes.Reader) error {
	if r.dir == "1" {
		value, err := buf.ReadByte()
		if err != nil {
			return err
		}
		if !checkData(value) {
			return nil
		}
		data := &RelayStationCommandUpF1{}
		err = data.Decode(value)
		if err != nil {
			return err
		}
		r.data = append(r.data, &Dlt13761Data{P: 0, F: 2, Data: data})
	}
	return nil
}

// 中继站工作状态控制
func (r *RelayStationCommand) statusControl(buf *bytes.Reader) error {
	var data Dlt13761DataInter
	value, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if !checkData(value) {
		return nil
	}
	if r.dir == "0" {
		//下行
		data = &RelayStationCommandDownF1{}
	} else {
		data = &RelayStationCommandUpF1{}
	}
	err = data.Decode(value)
	if err == nil {
		r.data = append(r.data, &Dlt13761Data{P: 0, F: 1, Data: data})
	}
	return nil
}

func (r *RelayStationCommand) Encode() ([]byte, error) {
	if r.dir == "0" {
		//下行
		//需要进行排序
		sortByF(r.data)
		fArr := make([]*uint64, len(r.data))
		fMap := make(map[uint64]*Dlt13761Data, len(r.data))
		for index, d := range r.data {
			fArr[index] = &d.F
			fMap[d.F] = d
		}
		//封装FN

	}
	return nil, nil
}

func (r *RelayStationCommand) Idents() ([]*Dlt13761Data, error) {
	return r.data, nil
}

func (r *RelayStationCommand) Flag() (byte, string) {
	return RELAY_STATION_COMMAND, "中继站命令"
}

func (r *RelayStationCommand) HasAux() bool {
	return false
}

func (r *RelayStationCommand) Append(data *Dlt13761Data) error {
	r.data = append(r.data, data)
	return nil
}

// Direction 确认方向
func (r *RelayStationCommand) Direction(dir string) {
	r.dir = dir
}

func init() {
	Afns[RELAY_STATION_COMMAND] = func() Afn {
		return &RelayStationCommand{}
	}
}
