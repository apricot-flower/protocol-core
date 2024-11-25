package dlt1376_1

import (
	"bytes"
	"errors"
)

var _ Afn = (*ComFirmDeny)(nil)

// ComFirmDeny 确认否认
type ComFirmDeny struct {
	pn uint64
	fn uint64
}

func (c *ComFirmDeny) Decode(buf *bytes.Reader) error {
	afnBytes := make([]byte, 4)
	_, err := buf.Read(afnBytes)
	if err != nil {
		return err
	}
	pn, fn := analyzeUnit(afnBytes)
	if len(pn) != 1 || len(fn) != 1 {
		return errors.New("afn.comFirmDeny: wrong number of parameters")
	}
	c.fn = fn[0]
	c.pn = pn[0]
	if c.pn != 0 {
		return errors.New("afn.comFirmDeny: pn must be 0")
	}
	return nil
}

func (c *ComFirmDeny) Encode() ([]byte, error) {
	return encodeUnit(c.pn, c.fn)
}

func (c *ComFirmDeny) Idents() ([]*Dlt13761Data, error) {
	var ret []*Dlt13761Data
	data := &Dlt13761Data{P: c.pn, F: c.fn, Data: nil}
	ret = append(ret, data)
	return ret, nil
}

func (c *ComFirmDeny) Flag() (byte, string) {
	return COMFIRM_DENY, "确认/否认"
}

func (c *ComFirmDeny) HasAux() bool {
	return false
}

func (c *ComFirmDeny) Append(data *Dlt13761Data) error {
	c.pn = data.P
	c.fn = data.F
	if c.pn != 0 {
		return errors.New("afn.comFirmDeny: pn must be 0")
	}
	return nil
}

func init() {
	Afns[COMFIRM_DENY] = func() Afn {
		return &ComFirmDeny{}
	}
}
