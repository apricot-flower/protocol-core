package dlt1376_1

import (
	"bytes"
	"errors"
)

var _ LinkDataInter = (*ConfirmDeny)(nil)

const ConfirmDenyIdent = 0x00

// ConfirmDeny 确认否认
type ConfirmDeny struct {
	F uint64
}

func (c *ConfirmDeny) decode(buf *bytes.Reader) error {
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
	c.F = fn[0]
	return nil
}

func (c *ConfirmDeny) encode() ([]byte, error) {
	if c.F == 1 || c.F == 2 {
		return encodeUnit(0, c.F)
	} else if c.F == 3 {
		return nil, nil
	} else if c.F == 4 {
		return nil, nil
	} else {
		return nil, errors.New("encode dlt1376.1 err : f is out of range, must in (1,2,3,4)")
	}

}

func (c *ConfirmDeny) haAux() bool {
	return true
}

func (c *ConfirmDeny) AfnFlag() byte {
	return ConfirmDenyIdent
}
