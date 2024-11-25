package dlt1376_1

import (
	"bytes"
	"errors"
)

var _ Afn = (*ResetCommand)(nil)

// ResetCommand 复位命令
type ResetCommand struct {
	pn uint64
	fn uint64
}

func (r *ResetCommand) Decode(buf *bytes.Reader) error {
	afnBytes := make([]byte, 4)
	_, err := buf.Read(afnBytes)
	if err != nil {
		return err
	}
	pn, fn := analyzeUnit(afnBytes)
	if len(pn) != 1 || len(fn) != 1 {
		return errors.New("afn.comFirmDeny: wrong number of parameters")
	}
	r.fn = fn[0]
	r.pn = pn[0]
	return nil
}

func (r *ResetCommand) Encode() ([]byte, error) {
	return encodeUnit(r.pn, r.fn)
}

func (r *ResetCommand) Idents() ([]*Dlt13761Data, error) {
	var ret []*Dlt13761Data
	data := &Dlt13761Data{P: r.pn, F: r.fn, Data: nil}
	ret = append(ret, data)
	return ret, nil
}

func (r *ResetCommand) Flag() (byte, string) {
	return RESET_COMMAND, "复位命令"
}

func (r *ResetCommand) HasAux() bool {
	return true
}

func (r *ResetCommand) Append(data *Dlt13761Data) error {
	r.pn = data.P
	r.fn = data.F
	if r.pn != 0 {
		return errors.New("afn.ResetCommand: pn must be 0")
	}
	return nil
}

func init() {
	Afns[RESET_COMMAND] = func() Afn {
		return &ResetCommand{}
	}
}
