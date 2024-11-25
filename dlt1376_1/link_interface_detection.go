package dlt1376_1

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var _ Afn = (*LinkInterfaceDetection)(nil)

// LinkInterfaceDetection 链路接口检测
type LinkInterfaceDetection struct {
	pn   uint64
	fn   uint64
	data *TimeA1
}

func (l *LinkInterfaceDetection) Decode(buf *bytes.Reader) error {
	afnBytes := make([]byte, 4)
	_, err := buf.Read(afnBytes)
	if err != nil {
		return err
	}
	pn, fn := analyzeUnit(afnBytes)
	if len(pn) != 1 || len(fn) != 1 {
		return errors.New("afn.LinkInterfaceDetection: wrong number of parameters")
	}
	l.fn = fn[0]
	l.pn = pn[0]
	if l.pn != 0 {
		return errors.New("afn.LinkInterfaceDetection: pn must be 0")
	}
	if l.fn == 0x03 {
		//解析心跳中的终端时钟
		timeArr := make([]byte, 6)
		err := binary.Read(buf, binary.LittleEndian, &timeArr)
		if err != nil {
			return err
		}
		if checkData(timeArr...) {
			l.data = &TimeA1{}
			err = l.data.Decode(timeArr...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *LinkInterfaceDetection) Encode() ([]byte, error) {
	encodeArray, err := encodeUnit(l.pn, l.fn)
	if err != nil {
		return nil, err
	}
	if l.fn == 0x03 {
		timerArray, err := l.data.Encode()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, timerArray...)
	}
	return encodeArray, nil
}

func (l *LinkInterfaceDetection) Idents() ([]*Dlt13761Data, error) {
	return []*Dlt13761Data{{
		P:    l.pn,
		F:    l.fn,
		Data: l.data,
	}}, nil
}

func (l *LinkInterfaceDetection) Flag() (byte, string) {
	return LINK_INTERFACE_DETECTION, "链路接口检测"
}

func (l *LinkInterfaceDetection) HasAux() bool {
	return false
}

func (l *LinkInterfaceDetection) Append(data *Dlt13761Data) error {
	l.pn = data.P
	l.fn = data.F
	if l.pn != 0 {
		return errors.New("afn.LinkInterfaceDetection: pn must be 0")
	}
	if l.fn == 0x03 {
		//心跳
		//创建一个timer
		l.data = &TimeA1{}
		err := l.data.Build()
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	Afns[LINK_INTERFACE_DETECTION] = func() Afn {
		return &LinkInterfaceDetection{}
	}
}
