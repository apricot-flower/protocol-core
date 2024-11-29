package dlt1376_1

import (
	"bytes"
	"errors"
)

var _ LinkDataInter = (*Reset)(nil)

const ResetIdent = 0x01

// Reset 复位
type Reset struct {
	F uint64
}

func (r *Reset) decode(buf *bytes.Reader) error {
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
	return nil
}

func (r *Reset) encode() ([]byte, error) {
	return encodeUnit(0, r.F)
}

func (r *Reset) haAux() bool {
	return true
}

func (r *Reset) AfnFlag() byte {
	return ResetIdent
}
