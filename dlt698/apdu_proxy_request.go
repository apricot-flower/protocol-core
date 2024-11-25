package dlt698

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
)

var _ APDURegion = (*ProxyGetRequestList)(nil)
var _ APDURegion = (*ProxyGetRequestRecord)(nil)

var _ FrameRegion = (*ProxyGetRequestListItem)(nil)

const (
	ProxyGetRequestListIdent   string = "0901"
	ProxyGetRequestRecordIdent string = "0902"
)

func init() {
	apduMap[ProxyGetRequestListIdent] = func() APDURegion {
		return new(ProxyGetRequestList)
	}
	apduMap[ProxyGetRequestRecordIdent] = func() APDURegion {
		return new(ProxyGetRequestRecord)
	}
}

type ProxyGetRequestList struct {
	TimeOut *LongUnsigned //代理整个请求的超时时间，单位秒，非0值
	Data    []*ProxyGetRequestListItem
}

func (p *ProxyGetRequestList) decoder(buf *bytes.Reader) error {
	p.TimeOut = &LongUnsigned{}
	if err := p.TimeOut.decoder(buf); err != nil {
		return err
	}
	length, err := buf.ReadByte()
	if err != nil {
		return err
	}
	p.Data = make([]*ProxyGetRequestListItem, length)
	for i := 0; i < int(length); i++ {
		p.Data[i] = &ProxyGetRequestListItem{}
		if err := p.Data[i].decoder(buf); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProxyGetRequestList) encoder() ([]byte, error) {
	encodeArray, err := p.TimeOut.encoder()
	if err != nil {
		return nil, err
	}
	length := len(p.Data)
	if length == 0 {
		return append(encodeArray, 0x00), nil
	}
	encodeArray = append(encodeArray, byte(length))
	for _, item := range p.Data {
		itemArray, err := item.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, itemArray...)
	}
	return encodeArray, err
}

func (p *ProxyGetRequestList) APDUType() string {
	return ProxyGetRequestListIdent
}

func (p *ProxyGetRequestList) APDUMark() string {
	return "proxy_get_request_list"
}

func (p *ProxyGetRequestList) hasFollowReport() bool {
	return false
}

func (p *ProxyGetRequestList) hasTimeTag() bool {
	return true
}

type ProxyGetRequestListItem struct {
	Tsa     string        //一个目标服务器地址
	TimeOut *LongUnsigned // 代理一个目标服气器的超时时间
	Oads    [][]byte      // 若干个对象属性描述符
}

func (p *ProxyGetRequestListItem) decoder(buf *bytes.Reader) error {
	var tsaAllLength byte
	if err := binary.Read(buf, binary.BigEndian, &tsaAllLength); err != nil {
		return errors.New("ProxyGetRequestListItem decode err: " + err.Error())
	}
	var tsaLen byte
	if err := binary.Read(buf, binary.BigEndian, &tsaLen); err != nil {
		return errors.New("ProxyGetRequestListItem decode err: " + err.Error())
	}
	if tsaAllLength != tsaLen+2 {
		return errors.New("tsaAllLength != tsaLen + 2")
	}
	if tsaLen > 0 {
		tsaArray := make([]byte, tsaLen+1)
		err := binary.Read(buf, binary.LittleEndian, &tsaArray)
		if err != nil {
			return err
		}
		p.Tsa = hex.EncodeToString(tsaArray)
	}
	p.TimeOut = &LongUnsigned{}
	if err := p.TimeOut.decoder(buf); err != nil {
		return err
	}
	var length byte
	if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
		return errors.New("ProxyGetRequestListItem decode err: " + err.Error())
	}
	p.Oads = make([][]byte, length)
	for i := 0; i < int(length); i++ {
		p.Oads[i] = make([]byte, 4)
		if err := binary.Read(buf, binary.LittleEndian, &p.Oads[i]); err != nil {
			return errors.New("ProxyGetRequestListItem decode err: " + err.Error())
		}
	}
	return nil
}

func (p *ProxyGetRequestListItem) encoder() ([]byte, error) {
	encodeArray, err := hex.DecodeString(p.Tsa)
	if err != nil {
		return nil, err
	}
	encodeArray = append([]byte{byte(len(encodeArray) + 1), byte(len(encodeArray) - 1)}, encodeArray...)
	toArray, err := p.TimeOut.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, toArray...)
	length := len(p.Oads)
	if length == 0 {
		return append(encodeArray, 0x00), nil
	}
	encodeArray = append(encodeArray, byte(length))
	for _, item := range p.Oads {
		encodeArray = append(encodeArray, item...)
	}
	return encodeArray, nil
}

/*-------------------------------------------*/

type ProxyGetRequestRecord struct {
	TimeOut uint16 // 代理请求的超时时间
	Tsa     string //目标服务器地址
	Oad     []byte //对象属性描述符
	Rsd     *RSD   // 记录行选择描述符
	Rcsd    *RCSD  //记录列选择描述符
}

func (p *ProxyGetRequestRecord) decoder(buf *bytes.Reader) error {
	err := binary.Read(buf, binary.LittleEndian, &p.TimeOut)
	if err != nil {
		return errors.New("ProxyGetRequestRecord decode err: " + err.Error())
	}
	tsaAllLength, err := buf.ReadByte()
	if err != nil {
		return errors.New("ProxyGetRequestRecord decode err: " + err.Error())
	}
	tsaLen, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if tsaAllLength != tsaLen+2 {
		return errors.New("tsaAllLength != tsaLen + 2")
	}
	if tsaLen > 0 {
		tsaArray := make([]byte, tsaLen+1)
		err := binary.Read(buf, binary.LittleEndian, &tsaArray)
		if err != nil {
			return err
		}
		p.Tsa = hex.EncodeToString(tsaArray)
	}
	p.Oad = make([]byte, 4)
	if err = binary.Read(buf, binary.LittleEndian, &p.Oad); err != nil {
		return err
	}
	p.Rsd = new(RSD)
	if err = p.Rsd.decoder(buf); err != nil {
		return err
	}
	p.Rcsd = new(RCSD)
	err = p.Rcsd.decoder(buf)
	return err
}

func (p *ProxyGetRequestRecord) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, p.TimeOut)
	if err != nil {
		return nil, err
	}
	encodeArray, err := hex.DecodeString(p.Tsa)
	if err != nil {
		return nil, err
	}
	encodeArray = append([]byte{byte(len(encodeArray) + 1), byte(len(encodeArray) - 1)}, encodeArray...)
	encodeArray = append(buf.Bytes(), encodeArray...)
	encodeArray = append(encodeArray, p.Oad...)
	rsdArray, err := p.Rsd.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, rsdArray...)
	rcsdArray, err := p.Rcsd.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, rcsdArray...)
	return encodeArray, nil
}

func (p *ProxyGetRequestRecord) APDUType() string {
	return ProxyGetRequestRecordIdent
}

func (p *ProxyGetRequestRecord) APDUMark() string {
	return "proxy_get_request_record"
}

func (p *ProxyGetRequestRecord) hasFollowReport() bool {
	return false
}

func (p *ProxyGetRequestRecord) hasTimeTag() bool {
	return true
}
