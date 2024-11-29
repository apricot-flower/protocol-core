package dlt1376_1

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
)

const (
	RelayStationCommandIdent byte = 0x03
)

var _ LinkDataInter = (*RelayStationCommand)(nil)

// RelayStationCommand 中继站命令
type RelayStationCommand struct {
	F   uint64
	Cmd *ControlOfWorkingStatusOfRelayStation
}

func (r *RelayStationCommand) decode(buf *bytes.Reader) error {
	pn, fn, err := analyzeUnit(buf)
	if err != nil {
		return err
	}
	if pn == nil || fn == nil || len(pn) == 0 || len(fn) == 0 {
		return errors.New("decode dlt1376.1 err : pn fn is empty")
	}
	for _, value := range pn {
		if value != 0 {
			return errors.New("decode dlt1376.1 err :pn must = 0")
		}
	}
	r.F = fn[0]
	if r.F == 0x01 {
		r.Cmd = &ControlOfWorkingStatusOfRelayStation{}
		err := r.Cmd.decode(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RelayStationCommand) encode() ([]byte, error) {
	encodeArray, err := encodeUnit(0, r.F)
	if err != nil {
		return nil, err
	}
	if r.F == 0x01 {
		dataArray, err := r.Cmd.encode()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, dataArray...)
	}
	return encodeArray, nil
}

func (r *RelayStationCommand) haAux() bool {
	return false
}

func (r *RelayStationCommand) AfnFlag() byte {
	return RelayStationCommandIdent
}

/*-----------------------子类---------------------------*/

// ControlOfWorkingStatusOfRelayStation 中继站工作状态控制
type ControlOfWorkingStatusOfRelayStation struct {
	MSSwitchControl string // D0～D1 值班机/备份机切换控制：D0=0、D1=0：表示不切换；D0=1、D1=1：表示切换；D0、D1为其他是无效
	DutyForwardFlag string //D2～D3 值班机中继转发允许标志：D2=0、D3=0：表示不允许；D2=1、D3=1：表示允许；D2、D3为其他是无效。
}

func (c *ControlOfWorkingStatusOfRelayStation) decode(buf *bytes.Reader) error {
	var data byte
	if err := binary.Read(buf, binary.LittleEndian, &data); err != nil {
		return err
	}
	bit0 := data & 0x01
	bit1 := (data >> 1) & 0x01
	bit2 := (data >> 2) & 0x01
	bit3 := (data >> 3) & 0x01
	c.MSSwitchControl = strconv.Itoa(int(bit1)) + strconv.Itoa(int(bit0))
	c.DutyForwardFlag = strconv.Itoa(int(bit3)) + strconv.Itoa(int(bit2))
	return nil
}

func (c *ControlOfWorkingStatusOfRelayStation) encode() ([]byte, error) {
	data := "0000" + c.DutyForwardFlag + c.MSSwitchControl
	num, err := strconv.ParseInt(data, 2, 64)
	if err != nil {
		return nil, errors.New("encode RelayStationCommand ControlOfWorkingStatusOfRelayStation err: " + err.Error())
	}
	return []byte{byte(num)}, nil
}
