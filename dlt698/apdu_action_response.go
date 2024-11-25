package dlt698

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var _ APDURegion = (*ActionResponseNormal)(nil)
var _ APDURegion = (*ActionResponseNormalList)(nil)
var _ APDURegion = (*ActionThenGetResponseNormalList)(nil)

var _ FrameRegion = (*ActionThenGetResponseNormalListItem)(nil)

const (
	ActionResponseNormalIdent            string = "8701"
	ActionResponseNormalListIdent        string = "8702"
	ActionThenGetResponseNormalListIdent string = "8703"
)

func init() {
	apduMap[ActionResponseNormalIdent] = func() APDURegion {
		return new(ActionResponseNormal)
	}
	apduMap[ActionResponseNormalListIdent] = func() APDURegion {
		return new(ActionResponseNormalList)
	}
	apduMap[ActionThenGetResponseNormalListIdent] = func() APDURegion {
		return new(ActionThenGetResponseNormalList)
	}
}

type ActionResponseNormal struct {
	Omd  *OMD      `json:"omd"`  // 一个对象方法描述符
	DAR  *DAR      `json:"dar"`  // 操作执行结果
	Data DataInter `json:"data"` // 操作返回数据
}

func (a *ActionResponseNormal) decoder(buf *bytes.Reader) error {
	a.Omd = &OMD{}
	if err := a.Omd.decoder(buf); err != nil {
		return err
	}
	a.DAR = &DAR{}
	if err := a.DAR.decoder(buf); err != nil {
		return err
	}
	var hasData byte
	err := binary.Read(buf, binary.LittleEndian, &hasData)
	if err != nil {
		return errors.New("ActionResponseNormal decode err: " + err.Error())
	}
	if hasData == 0 {
		return nil
	}
	//获取data类型
	var dataType byte
	err = binary.Read(buf, binary.LittleEndian, &dataType)
	if err != nil {
		return errors.New("ActionResponseNormal decode err: " + err.Error())
	}
	a.Data = dataTranslate(dataType)
	if a.Data == nil {
		return errors.New("GetResult data type not found")
	}
	err = a.Data.decoder(buf)
	return err
}

func (a *ActionResponseNormal) encoder() ([]byte, error) {
	encodeArray, err := a.Omd.encoder()
	if err != nil {
		return nil, err
	}
	darArray, err := a.DAR.encoder()
	if err != nil {
		return nil, err
	}
	if a.Data == nil {
		return append(encodeArray, 0x00), nil
	}
	encodeArray = append(encodeArray, darArray...)
	encodeArray = append(encodeArray, 0x01, a.Data.DataType())
	dataArray, err := a.Data.encoder()
	if err != nil {
		return nil, err
	}
	return append(encodeArray, dataArray...), nil
}

func (a *ActionResponseNormal) APDUType() string {
	return ActionResponseNormalIdent
}

func (a *ActionResponseNormal) APDUMark() string {
	return "action_response_normal"
}

func (a *ActionResponseNormal) hasFollowReport() bool {
	return true
}

func (a *ActionResponseNormal) hasTimeTag() bool {
	return true
}

/*---------------------------------*/

type ActionResponseNormalList struct {
	Data []*ActionResponseNormal `json:"data"`
}

func (a *ActionResponseNormalList) decoder(buf *bytes.Reader) error {
	var length byte
	err := binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return errors.New("ActionResponseNormalList decoder err : " + err.Error())
	}
	a.Data = make([]*ActionResponseNormal, length)
	for i := 0; i < int(length); i++ {
		a.Data[i] = new(ActionResponseNormal)
		err = a.Data[i].decoder(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *ActionResponseNormalList) encoder() ([]byte, error) {
	if a.Data == nil || len(a.Data) == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := []byte{byte(len(a.Data))}
	for _, data := range a.Data {
		dataArray, err := data.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, dataArray...)
	}
	return encodeArray, nil
}

func (a *ActionResponseNormalList) APDUType() string {
	return ActionResponseNormalListIdent
}

func (a *ActionResponseNormalList) APDUMark() string {
	return "action_response_normal_list"
}

func (a *ActionResponseNormalList) hasFollowReport() bool {
	return true
}

func (a *ActionResponseNormalList) hasTimeTag() bool {
	return true
}

/*----------------------------------------*/

type ActionThenGetResponseNormalList struct {
	Data []*ActionThenGetResponseNormalListItem `json:"data"`
}

func (a *ActionThenGetResponseNormalList) decoder(buf *bytes.Reader) error {
	length, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if length == 0 {
		return nil
	}
	a.Data = make([]*ActionThenGetResponseNormalListItem, length)
	for i := 0; i < int(length); i++ {
		a.Data[i] = &ActionThenGetResponseNormalListItem{}
		err = a.Data[i].decoder(buf)
		if err != nil {
			return err
		}
	}
	return err
}

func (a *ActionThenGetResponseNormalList) encoder() ([]byte, error) {
	if a.Data == nil || len(a.Data) == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := []byte{byte(len(a.Data))}
	for _, data := range a.Data {
		dataArray, err := data.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, dataArray...)
	}
	return encodeArray, nil
}

func (a *ActionThenGetResponseNormalList) APDUType() string {
	return ActionThenGetResponseNormalListIdent
}

func (a *ActionThenGetResponseNormalList) APDUMark() string {
	return "action_then_get_response_normal_list"
}

func (a *ActionThenGetResponseNormalList) hasFollowReport() bool {
	return true
}

func (a *ActionThenGetResponseNormalList) hasTimeTag() bool {
	return true
}

type ActionThenGetResponseNormalListItem struct {
	Omd          *OMD          `json:"omd"`           //一个设置的对象方法描述符
	Dar          *DAR          `json:"dar"`           //操作执行结果
	Data         DataInter     `json:"data"`          //操作返回数据
	ResultNormal *ResultNormal `json:"result_normal"` //一个对象属性及结果
}

func (a *ActionThenGetResponseNormalListItem) decoder(buf *bytes.Reader) error {
	a.Omd = &OMD{}
	if err := a.Omd.decoder(buf); err != nil {
		return err
	}
	a.Dar = &DAR{}
	if err := a.Dar.decoder(buf); err != nil {
		return err
	}
	hasData, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if hasData == 1 {
		//获取data类型
		dataType, err := buf.ReadByte()
		if err != nil {
			return err
		}
		a.Data = dataTranslate(dataType)
		if a.Data == nil {
			return errors.New("GetResult data type not found！")
		}
		err = a.Data.decoder(buf)
		if err != nil {
			return err
		}
	}
	a.ResultNormal = &ResultNormal{}
	err = a.ResultNormal.decoder(buf)
	return err
}

func (a *ActionThenGetResponseNormalListItem) encoder() ([]byte, error) {
	encodeArray, err := a.Omd.encoder()
	if err != nil {
		return nil, err
	}
	darArray, err := a.Dar.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, darArray...)
	if a.Data == nil {
		encodeArray = append(encodeArray, 0x00)
	} else {
		encodeArray = append(encodeArray, 0x01, a.Data.DataType())
		dataArray, err := a.Data.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, dataArray...)
	}
	resultNormalArray, err := a.ResultNormal.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, resultNormalArray...)
	return encodeArray, nil
}
