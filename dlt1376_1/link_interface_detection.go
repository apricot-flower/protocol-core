package dlt1376_1

import (
	"bytes"
	"errors"
)

const (
	LinkInterfaceDetectionIdent byte = 0x02
)

var _ LinkDataInter = (*LinkInterfaceDetection)(nil)

// LinkInterfaceDetection 链路接口检测
type LinkInterfaceDetection struct {
	F             uint64
	TerminalClock *TimeA1 //终端时钟
}

func (l *LinkInterfaceDetection) decode(buf *bytes.Reader) error {
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
	l.F = fn[0]
	if l.F == 0x03 {
		l.TerminalClock = new(TimeA1)
		err = l.TerminalClock.decode(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LinkInterfaceDetection) encode() ([]byte, error) {
	encodeArray, err := encodeUnit(0, l.F)
	if err != nil {
		return nil, err
	}
	if l.F == 0x03 {
		if l.TerminalClock == nil {
			l.TerminalClock = new(TimeA1)
			err = l.TerminalClock.Build()
			if err != nil {
				return nil, err
			}
		}
		clockArray, err := l.TerminalClock.encode()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, clockArray...)
	}
	return encodeArray, nil
}

func (l *LinkInterfaceDetection) haAux() bool {
	return false
}

func (l *LinkInterfaceDetection) AfnFlag() byte {
	return LinkInterfaceDetectionIdent
}
