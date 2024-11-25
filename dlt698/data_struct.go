package dlt698

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
)

var _ DataInter = (*OI)(nil)
var _ DataInter = (*OAD)(nil)
var _ DataInter = (*ROAD)(nil)
var _ DataInter = (*OMD)(nil)
var _ DataInter = (*TSA)(nil)
var _ DataInter = (*MAC)(nil)
var _ DataInter = (*Region)(nil)
var _ DataInter = (*ScalerUnit)(nil)
var _ DataInter = (*RSD)(nil)
var _ DataInter = (*CSD)(nil)
var _ DataInter = (*MS)(nil)
var _ DataInter = (*SID)(nil)
var _ DataInter = (*SIDMAC)(nil)
var _ DataInter = (*COMDCB)(nil)
var _ DataInter = (*RCSD)(nil)

const (
	SingleAddress       int = 0x00 //单地址
	DistributionAddress int = 0x01 //通配地址
	GroupAddress        int = 0x02 //主地址
	BroadcastAddress    int = 0x03 //广播地址

	OIIdent         byte = 80
	OADIdent        byte = 81
	ROADIdent       byte = 82
	OMDIdent        byte = 83
	TSAIdent        byte = 85
	MACIdent        byte = 86
	RNIdent         byte = 87
	RegionIdent     byte = 88
	ScalerUnitIdent byte = 89
	RSDIdent        byte = 90
	CSDIdent        byte = 91
	MSIdent         byte = 92
	SIDIdent        byte = 93
	SIDMACIdent     byte = 94
	COMDCBIdent     byte = 95
	RCSDIdent       byte = 96
)

func init() {
	dataMap[OIIdent] = func() DataInter {
		return &OI{}
	}
	dataMap[OADIdent] = func() DataInter {
		return &OAD{}
	}
	dataMap[ROADIdent] = func() DataInter {
		return &ROAD{}
	}
	dataMap[OMDIdent] = func() DataInter {
		return &OMD{}
	}
	dataMap[TSAIdent] = func() DataInter {
		return &TSA{}
	}
	dataMap[MACIdent] = func() DataInter {
		return &MAC{}
	}
	dataMap[RNIdent] = func() DataInter {
		return &RN{}
	}
	dataMap[RegionIdent] = func() DataInter {
		return &Region{}
	}
	dataMap[ScalerUnitIdent] = func() DataInter {
		return &ScalerUnit{}
	}
	dataMap[RSDIdent] = func() DataInter {
		return &RSD{}
	}
	dataMap[CSDIdent] = func() DataInter {
		return &CSD{}
	}
	dataMap[MSIdent] = func() DataInter {
		return &MS{}
	}
	dataMap[SIDIdent] = func() DataInter {
		return &SID{}
	}
	dataMap[SIDMACIdent] = func() DataInter {
		return &SIDMAC{}
	}
	dataMap[COMDCBIdent] = func() DataInter {
		return &COMDCB{}
	}
	dataMap[RCSDIdent] = func() DataInter {
		return &RCSD{}
	}
}

type OI struct {
	Data uint16 `json:"value"`
}

func (O *OI) decoder(buf *bytes.Reader) error {
	return binary.Read(buf, binary.BigEndian, &O.Data)
}

func (O *OI) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, O.Value)
	return buf.Bytes(), err
}

func (O *OI) DataType() byte {
	return OIIdent
}

func (O *OI) Value() interface{} {
	return O.Data
}

/*-----------------------------------*/

type OAD struct {
	Data []byte `json:"value"`
}

func (O *OAD) decoder(buf *bytes.Reader) error {
	O.Data = make([]byte, 4)
	return binary.Read(buf, binary.LittleEndian, &O.Data)
}

func (O *OAD) encoder() ([]byte, error) {
	return O.Data, nil
}

func (O *OAD) DataType() byte {
	return OADIdent
}

func (O *OAD) Value() interface{} {
	return O.Data
}

/*-------------------------*/

type ROAD struct {
	Oad  *OAD   `json:"oad"`
	Oads []*OAD `json:"oads"`
}

func (R *ROAD) decoder(buf *bytes.Reader) error {
	R.Oad = &OAD{}
	err := binary.Read(buf, binary.LittleEndian, &R.Oad)
	if err != nil {
		return errors.New("decode ROAD err: " + err.Error())
	}
	var oadLen byte
	err = binary.Read(buf, binary.LittleEndian, &oadLen)
	if err != nil {
		return errors.New("decode ROAD err: " + err.Error())
	}
	R.Oads = make([]*OAD, oadLen)
	for i := 0; i < int(oadLen); i++ {
		R.Oads[i] = &OAD{}
		err = R.Oads[i].decoder(buf)
		if err != nil {
			return errors.New("decode ROAD err: " + err.Error())
		}
	}
	return nil
}

func (R *ROAD) encoder() ([]byte, error) {
	encodeArray, err := R.Oad.encoder()
	if err != nil {
		return nil, errors.New("encode ROAD err: " + err.Error())
	}
	if len(R.Oads) == 0 {
		return append(encodeArray, 0x00), nil
	}
	encodeArray = append(encodeArray, byte(len(R.Oads)))
	for _, oad := range R.Oads {
		oadArray, err := oad.encoder()
		if err != nil {
			return nil, errors.New("encode ROAD err: " + err.Error())
		}
		encodeArray = append(encodeArray, oadArray...)
	}
	return encodeArray, nil
}

func (R *ROAD) DataType() byte {
	return ROADIdent
}

func (R *ROAD) Value() interface{} {
	return R
}

/*----------------------------------*/

type OMD struct {
	Oi       *OI   `json:"oi"`
	FuncMark uint8 `json:"func_mark"`
	Mode     uint8 `json:"mode"`
}

func (O *OMD) decoder(buf *bytes.Reader) error {
	O.Oi = &OI{}
	err := O.Oi.decoder(buf)
	if err != nil {
		return errors.New("decode OMD err: " + err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &O.FuncMark)
	if err != nil {
		return errors.New("decode OMD err: " + err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &O.Mode)
	return err
}

func (O *OMD) encoder() ([]byte, error) {
	oiArray, err := O.Oi.encoder()
	if err != nil {
		return nil, errors.New("encode OMD err: " + err.Error())
	}
	return append(oiArray, O.FuncMark, O.Mode), nil
}

func (O *OMD) DataType() byte {
	return OMDIdent
}

func (O *OMD) Value() interface{} {
	return O
}

/*----------------------------------*/

type TSA struct {
	AddressType   int    `json:"address_type"`   //地址类型，0-单地址， 1-统配地址， 2-组地址， 3-广播地址
	LogicAddress  int    `json:"logic_address"`  //逻辑地址
	AddressLength int    `json:"address_length"` //地址长度 从0开始要加1
	Address       string `json:"address"`        //地址
	CA            byte   `json:"ca"`             //客户机地址
}

func (T *TSA) decoder(buf *bytes.Reader) error {
	at, err := buf.ReadByte()
	if err != nil {
		return err
	}
	T.AddressType = int((at >> 6) & 0b11)
	if T.AddressType != SingleAddress && T.AddressType != DistributionAddress && T.AddressType != GroupAddress && T.AddressType != BroadcastAddress {
		return errors.New("TSA address type error：" + strconv.Itoa(T.AddressType))
	}
	//获取逻辑地址
	T.LogicAddress = int((at >> 4) & 0b11)
	T.AddressLength = int((at>>0)&0b1111) + 1
	addressArray := make([]byte, T.AddressLength)
	_, err = buf.Read(addressArray)
	if err != nil {
		return errors.New("decode TSA err:" + err.Error())
	}
	for i, j := 0, len(addressArray)-1; i < j; i, j = i+1, j-1 {
		addressArray[i], addressArray[j] = addressArray[j], addressArray[i]
	}
	T.Address = hex.EncodeToString(addressArray)
	T.CA, err = buf.ReadByte()
	if err != nil {
		return errors.New("decode TSA err:" + err.Error())
	}
	return nil
}

func (T *TSA) encoder() ([]byte, error) {
	addressTypeStr := fmt.Sprintf("%d%d", (T.AddressType>>1)&1, T.AddressType&1)
	logicAddressStr := fmt.Sprintf("%d%d", (T.LogicAddress>>1)&1, T.LogicAddress&1)
	length := len(T.Address)/2 - 1
	addressLengthStr := strconv.FormatInt(int64(length&0b1111), 2)
	adType, err := strconv.ParseUint(addressTypeStr+logicAddressStr+addressLengthStr, 2, 8)
	if err != nil {
		return nil, errors.New("encode TSA err:" + err.Error())
	}
	addressArray, err := hex.DecodeString(T.Address)
	if err != nil {
		return nil, errors.New("encode TSA err:" + err.Error())
	}
	for i, j := 0, len(addressArray)-1; i < j; i, j = i+1, j-1 {
		addressArray[i], addressArray[j] = addressArray[j], addressArray[i]
	}
	return append(append([]byte{byte(adType)}, addressArray...), []byte{T.CA}...), nil
}

func (T *TSA) DataType() byte {
	return TSAIdent
}

func (T *TSA) Value() interface{} {
	return T
}

/*------------------------*/

type MAC struct {
	OctetString
}

func (M *MAC) DataType() byte {
	return MACIdent
}

func (M *MAC) Value() interface{} {
	return M.Data
}

/*----------------------*/

type RN struct {
	OctetString
}

func (R *RN) DataType() byte {
	return RNIdent
}

func (R *RN) Value() interface{} {
	return R.Data
}

/*---------------------------*/

type Region struct {
	Unit  byte      `json:"unit"`
	Start DataInter `json:"start"`
	End   DataInter `json:"end"`
}

func (r *Region) decoder(buf *bytes.Reader) error {
	err := binary.Read(buf, binary.LittleEndian, &r.Unit)
	if err != nil {
		return errors.New("decode Region err:" + err.Error())
	}
	var dataType byte
	err = binary.Read(buf, binary.LittleEndian, &dataType)
	if err != nil {
		return errors.New("decode Region err:" + err.Error())
	}
	r.Start = dataTranslate(dataType)
	err = r.Start.decoder(buf)
	if err != nil {
		return errors.New("decode Region err:" + err.Error())
	}
	err = binary.Read(buf, binary.LittleEndian, &dataType)
	if err != nil {
		return errors.New("decode Region err:" + err.Error())
	}
	r.End = dataTranslate(dataType)
	err = r.End.decoder(buf)
	if err != nil {
		return errors.New("decode Region err:" + err.Error())
	}
	return nil
}

func (r *Region) encoder() ([]byte, error) {
	encodeArray := []byte{r.Unit}
	encodeArray = append(encodeArray, r.Start.DataType())
	startArray, err := r.Start.encoder()
	if err != nil {
		return nil, errors.New("encode Region err:" + err.Error())
	}
	encodeArray = append(encodeArray, startArray...)
	encodeArray = append(encodeArray, r.End.DataType())
	endArray, err := r.End.encoder()
	if err != nil {
		return nil, errors.New("encode Region err:" + err.Error())
	}
	return append(encodeArray, endArray...), nil
}

func (r *Region) DataType() byte {
	return RegionIdent
}

func (r *Region) Value() interface{} {
	return r
}

/*--------------------------------------*/

type ScalerUnit struct {
	Conver *Integer `json:"conver"`
	Unit   byte     `json:"unsigned"`
}

func (s *ScalerUnit) decoder(buf *bytes.Reader) error {
	s.Conver = &Integer{}
	err := s.Conver.decoder(buf)
	if err != nil {
		return err
	}
	s.Unit, err = buf.ReadByte()
	return err
}

func (s *ScalerUnit) encoder() ([]byte, error) {
	converArray, err := s.Conver.encoder()
	if err != nil {
		return nil, err
	}
	return append(converArray, s.Unit), nil
}

func (s *ScalerUnit) DataType() byte {
	return ScalerUnitIdent
}

func (s *ScalerUnit) Value() interface{} {
	return s
}

/*--------------------------------*/

type Selector interface {
	// Decode 解码
	Decode(buf *bytes.Reader) error
	// Encode 编码
	Encode() ([]byte, error)
	// GetDataIdent 获取编码
	GetDataIdent() byte
}

var _ Selector = (*Selector0)(nil)
var _ Selector = (*Selector1)(nil)
var _ Selector = (*Selector2)(nil)
var _ Selector = (*Selector3)(nil)
var _ Selector = (*Selector4)(nil)
var _ Selector = (*Selector5)(nil)
var _ Selector = (*Selector6)(nil)
var _ Selector = (*Selector7)(nil)
var _ Selector = (*Selector8)(nil)
var _ Selector = (*Selector9)(nil)
var _ Selector = (*Selector10)(nil)

type RSD struct {
	Selector Selector `json:"selector"`
}

func (R *RSD) decoder(buf *bytes.Reader) error {
	rsdType, err := buf.ReadByte()
	if err != nil {
		return err
	}
	switch rsdType {
	case 0:
		R.Selector = &Selector0{}
	case 1:
		R.Selector = &Selector1{}
	case 2:
		R.Selector = &Selector2{}
	case 3:
		R.Selector = &Selector3{}
	case 4:
		R.Selector = &Selector4{}
	case 5:
		R.Selector = &Selector5{}
	case 6:
		R.Selector = &Selector6{}
	case 7:
		R.Selector = &Selector7{}
	case 8:
		R.Selector = &Selector8{}
	case 9:
		R.Selector = &Selector9{}
	case 10:
		R.Selector = &Selector10{}
	}
	if R.Selector == nil {
		return errors.New("RSD Selector must be not null！")
	}
	err = R.Selector.Decode(buf)
	return err
}

func (R *RSD) encoder() ([]byte, error) {
	rsdArray, err := R.Selector.Encode()
	if err != nil {
		return nil, err
	}
	return append([]byte{R.Selector.GetDataIdent()}, rsdArray...), nil
}

func (R *RSD) DataType() byte {
	return RSDIdent
}

func (R *RSD) Value() interface{} {
	return R
}

type Selector0 struct{}

func (s *Selector0) Decode(buf *bytes.Reader) error {
	return nil
}

func (s *Selector0) Encode() ([]byte, error) {
	return []byte{0x00}, nil
}

func (s *Selector0) GetDataIdent() byte {
	return 0
}

// Selector1 选择对象的指定值
type Selector1 struct {
	Oad  []byte    `json:"oad"`  //对象属性描述符
	Data DataInter `json:"data"` //数值
}

func (s *Selector1) Decode(buf *bytes.Reader) error {
	s.Oad = make([]byte, 4)
	if _, err := buf.Read(s.Oad); err != nil {
		return errors.New("decode Selector1 err:" + err.Error())
	}
	//解析data
	//获取data的标识
	dataFlg, err := buf.ReadByte()
	if err != nil {
		return errors.New("decode Selector1 err:" + err.Error())
	}
	s.Data = dataTranslate(dataFlg)
	err = s.Data.decoder(buf)
	if err != nil {
		return errors.New("decode Selector1 err:" + err.Error())
	}
	return nil
}

func (s *Selector1) Encode() ([]byte, error) {
	encodeArray := append(s.Oad, s.Data.DataType())
	dataArray, err := s.Data.encoder()
	if err != nil {
		return nil, errors.New("encode Selector1 err:" + err.Error())
	}
	return append(encodeArray, dataArray...), nil
}

func (s *Selector1) GetDataIdent() byte {
	return 1
}

type Selector2 struct {
	Oad       []byte    `json:"oad"`        //对象属性描述符
	StartData DataInter `json:"start_data"` //起始值
	EndData   DataInter `json:"end_data"`   //结束值
	Interval  DataInter `json:"interval"`   //数据间隔
}

func (s *Selector2) Decode(buf *bytes.Reader) error {
	s.Oad = make([]byte, 4)
	if _, err := buf.Read(s.Oad); err != nil {
		return errors.New("decode Selector2 err:" + err.Error())
	}
	startDataIdent, err := buf.ReadByte()
	if err != nil {
		return errors.New("decode Selector2 err:" + err.Error())
	}
	s.StartData = dataTranslate(startDataIdent)
	if s.StartData == nil {
		return errors.New("selector2 start err！")
	}
	if err = s.StartData.decoder(buf); err != nil {
		return errors.New("decode Selector2 err:" + err.Error())
	}
	endDataIdent, err := buf.ReadByte()
	if err != nil {
		return err
	}
	s.EndData = dataTranslate(endDataIdent)
	if s.EndData == nil {
		return errors.New("selector2 end type err！")
	}
	if err = s.EndData.decoder(buf); err != nil {
		return err
	}
	intervalIdent, err := buf.ReadByte()
	if err != nil {
		return err
	}
	s.Interval = dataTranslate(intervalIdent)
	if s.Interval == nil {
		return errors.New("selector2 interval decode err！")
	}
	if err = s.Interval.decoder(buf); err != nil {
		return errors.New("decode Selector2 err:" + err.Error())
	}
	return nil
}

func (s *Selector2) Encode() ([]byte, error) {
	encodeArray := append(s.Oad, s.StartData.DataType())
	startDataArray, err := s.StartData.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, startDataArray...)
	encodeArray = append(encodeArray, s.EndData.DataType())
	endDataArray, err := s.EndData.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, endDataArray...)
	encodeArray = append(encodeArray, s.Interval.DataType())
	intervalArray, err := s.Interval.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, intervalArray...)
	return encodeArray, nil
}

func (s *Selector2) GetDataIdent() byte {
	return 2
}

// Selector3 多个选择对象区间内连续间隔值的并集
type Selector3 struct {
	Selectors []*Selector2 `json:"selectors"`
}

func (s *Selector3) Decode(buf *bytes.Reader) error {
	selectorLen, err := buf.ReadByte()
	if err != nil {
		return err
	}
	for i := 0; i < int(selectorLen); i++ {
		selector := &Selector2{}
		if err := selector.Decode(buf); err != nil {
			return err
		}
		s.Selectors = append(s.Selectors, selector)
	}
	return nil
}

func (s *Selector3) Encode() ([]byte, error) {
	encodeArray := []byte{byte(len(s.Selectors))}
	for _, selector := range s.Selectors {
		selectorArray, err := selector.Encode()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, selectorArray...)
	}
	return encodeArray, nil
}

func (s *Selector3) GetDataIdent() byte {
	return 3
}

// Selector4 指定表计集合、指定采集启动时间
type Selector4 struct {
	StartTime *DateTimes `json:"start_time"` //采集启动时间
	MS        *MS        `json:"ms"`         //表计集合
}

func (s *Selector4) Decode(buf *bytes.Reader) error {
	s.StartTime = &DateTimes{}
	if err := s.StartTime.decoder(buf); err != nil {
		return err
	}
	s.MS = &MS{}
	if err := s.MS.decoder(buf); err != nil {
		return err
	}
	return nil
}

func (s *Selector4) Encode() ([]byte, error) {
	startTimeArray, err := s.StartTime.encoder()
	if err != nil {
		return nil, err
	}
	msArray, err := s.MS.encoder()
	if err != nil {
		return nil, err
	}
	return append(startTimeArray, msArray...), nil
}

func (s *Selector4) GetDataIdent() byte {
	return 4
}

type Selector5 struct {
	Selector4
}

func (s *Selector5) GetDataIdent() byte {
	return 5
}

type Selector6 struct {
	StartTime *DateTimes `json:"start_time"`
	EndTime   *DateTimes `json:"end_time"`
	Ti        *TI        `json:"ti"`
	Ms        *MS        `json:"ms"`
}

func (s *Selector6) Decode(buf *bytes.Reader) error {
	s.StartTime = &DateTimes{}
	if err := s.StartTime.decoder(buf); err != nil {
		return err
	}
	s.EndTime = &DateTimes{}
	if err := s.EndTime.decoder(buf); err != nil {
		return err
	}
	s.Ti = &TI{}
	if err := s.Ti.decoder(buf); err != nil {
		return err
	}
	s.Ms = &MS{}
	if err := s.Ms.decoder(buf); err != nil {
		return err
	}
	return nil
}

func (s *Selector6) Encode() ([]byte, error) {
	var err error
	startTimeArray, err := s.StartTime.encoder()
	if err != nil {
		return nil, err
	}
	endTimeArray, err := s.EndTime.encoder()
	if err != nil {
		return nil, err
	}
	startTimeArray = append(startTimeArray, endTimeArray...)
	tsArray, err := s.Ti.encoder()
	if err != nil {
		return nil, err
	}
	startTimeArray = append(startTimeArray, tsArray...)
	msArray, err := s.Ms.encoder()
	if err != nil {
		return nil, err
	}
	startTimeArray = append(startTimeArray, msArray...)
	return startTimeArray, nil
}

func (s *Selector6) GetDataIdent() byte {
	return 6
}

type Selector7 struct {
	Selector6
}

func (s *Selector7) GetDataIdent() byte {
	return 7
}

type Selector8 struct {
	Selector6
}

func (s *Selector8) GetDataIdent() byte {
	return 8
}

type Selector9 struct {
	Last uint8 `json:"last_n"`
}

func (s *Selector9) Decode(buf *bytes.Reader) error {
	var err error
	s.Last, err = buf.ReadByte()
	return err
}

func (s *Selector9) Encode() ([]byte, error) {
	return []byte{byte(s.Last)}, nil
}

func (s *Selector9) GetDataIdent() byte {
	return 9
}

type Selector10 struct {
	Last uint8 `json:"last_n"`
	Ms   *MS   `json:"ms"`
}

func (s *Selector10) Decode(buf *bytes.Reader) error {
	var err error
	s.Last, err = buf.ReadByte()
	if err != nil {
		return err
	}
	err = s.Ms.decoder(buf)
	return err
}

func (s *Selector10) Encode() ([]byte, error) {
	msArray, err := s.Ms.encoder()
	if err != nil {
		return nil, err
	}
	return append([]byte{s.Last}, msArray...), nil
}

func (s *Selector10) GetDataIdent() byte {
	return 10
}

/*---------------------*/

type CSD struct {
	CsdType byte   `json:"csd_type"`
	Oad     []byte `json:"oad"`
}

func (C *CSD) decoder(buf *bytes.Reader) error {
	err := binary.Read(buf, binary.LittleEndian, &C.CsdType)
	if err != nil {
		return errors.New("decode CSD err:" + err.Error())
	}
	C.Oad = make([]byte, 4)
	err = binary.Read(buf, binary.BigEndian, &C.Oad)
	return err
}

func (C *CSD) encoder() ([]byte, error) {
	return append([]byte{C.CsdType}, C.Oad...), nil
}

func (C *CSD) DataType() byte {
	return CSDIdent
}

func (C *CSD) Value() interface{} {
	return C
}

/*--------------------------*/

type MSer interface {
	// Decode 解码
	Decode(buf *bytes.Reader) error
	// Encode 编码
	Encode() ([]byte, error)
	// GetDataIdent 获取编码
	GetDataIdent() byte
}

var _ MSer = (*MS0)(nil)
var _ MSer = (*MS1)(nil)
var _ MSer = (*MS2)(nil)
var _ MSer = (*MS3)(nil)
var _ MSer = (*MS4)(nil)
var _ MSer = (*MS5)(nil)
var _ MSer = (*MS6)(nil)
var _ MSer = (*MS7)(nil)

type MS struct {
	MSType byte `json:"ms_type"`
	Ms     MSer `json:"ms"`
}

func (M *MS) decoder(buf *bytes.Reader) error {
	err := binary.Read(buf, binary.BigEndian, &M.MSType)
	if err != nil {
		return errors.New("decode MS err : " + err.Error())
	}
	switch M.MSType {
	case 0:
		M.Ms = &MS0{}
	case 1:
		M.Ms = &MS1{}
	case 2:
		M.Ms = &MS2{}
	case 3:
		M.Ms = &MS3{}
	case 4:
		M.Ms = &MS4{}
	case 5:
		M.Ms = &MS5{}
	case 6:
		M.Ms = &MS6{}
	case 7:
		M.Ms = &MS7{}
	}
	if M.Ms == nil {
		return errors.New("MS type err！")
	}
	err = M.Ms.Decode(buf)
	if err != nil {
		return errors.New("decode MS err : " + err.Error())
	}
	return nil
}

func (M *MS) encoder() ([]byte, error) {
	msArray, err := M.Ms.Encode()
	if err != nil {
		return nil, errors.New("encode MS err : " + err.Error())
	}
	return append([]byte{M.Ms.GetDataIdent()}, msArray...), nil
}

func (M *MS) DataType() byte {
	return MSIdent
}

func (M *MS) Value() interface{} {
	return M
}

type MS0 struct {
}

func (M MS0) Decode(buf *bytes.Reader) error {
	return nil
}

func (M MS0) Encode() ([]byte, error) {
	return []byte{0x00}, nil
}

func (M MS0) GetDataIdent() byte {
	return 0
}

type MS1 struct {
}

func (M *MS1) Decode(buf *bytes.Reader) error {
	return nil
}

func (M *MS1) Encode() ([]byte, error) {
	return []byte{0x01}, nil
}

func (M *MS1) GetDataIdent() byte {
	return 1
}

type MS2 struct {
	Meters []uint8 `json:"meters"`
}

func (M *MS2) Decode(buf *bytes.Reader) error {
	metersLength, err := buf.ReadByte()
	if err != nil {
		return err
	}
	M.Meters = make([]uint8, metersLength)
	err = binary.Read(buf, binary.BigEndian, &M.Meters)
	return err
}

func (M *MS2) Encode() ([]byte, error) {
	return append([]byte{0x02}, M.Meters...), nil
}

func (M *MS2) GetDataIdent() byte {
	return 2
}

type MS3 struct {
	TSAs []*TSA `json:"tsas"`
}

func (M *MS3) Decode(buf *bytes.Reader) error {
	tasLen, err := buf.ReadByte()
	if err != nil {
		return err
	}
	M.TSAs = make([]*TSA, tasLen)
	for i := 0; i < int(tasLen); i++ {
		tsa := &TSA{}
		err = tsa.decoder(buf)
		if err != nil {
			return err
		}
		M.TSAs[i] = tsa
	}
	return nil
}

func (M *MS3) Encode() ([]byte, error) {
	if len(M.TSAs) == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := []byte{byte(len(M.TSAs))}
	for _, tsa := range M.TSAs {
		tsaArray, err := tsa.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, tsaArray...)
	}
	return encodeArray, nil
}

func (M *MS3) GetDataIdent() byte {
	return 3
}

type MS4 struct {
	Meters []uint16 `json:"meters"`
}

func (M *MS4) Decode(buf *bytes.Reader) error {
	metersLength, err := buf.ReadByte()
	if err != nil {
		return err
	}
	M.Meters = make([]uint16, metersLength)
	err = binary.Read(buf, binary.BigEndian, &M.Meters)
	return err
}

func (M *MS4) Encode() ([]byte, error) {
	if len(M.Meters) == 0 {
		return []byte{0x00}, nil
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, byte(len(M.Meters)))
	for _, meter := range M.Meters {
		err = binary.Write(buf, binary.BigEndian, meter)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (M *MS4) GetDataIdent() byte {
	return 4
}

type MS5 struct {
	Regions []*Region `json:"regions"`
}

func (M *MS5) Decode(buf *bytes.Reader) error {
	regionLength, err := buf.ReadByte()
	if err != nil {
		return err
	}
	M.Regions = make([]*Region, regionLength)
	for i := 0; i < int(regionLength); i++ {
		M.Regions[i] = &Region{}
		err = M.Regions[i].decoder(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (M *MS5) Encode() ([]byte, error) {
	if len(M.Regions) == 0 {
		return []byte{0x00}, nil
	}
	encodeArray := []byte{byte(len(M.Regions))}
	for _, region := range M.Regions {
		regionArray, err := region.encoder()
		if err != nil {
			return nil, err
		}
		encodeArray = append(encodeArray, regionArray...)
	}
	return encodeArray, nil
}

func (M *MS5) GetDataIdent() byte {
	return 5
}

type MS6 struct {
	MS5
}

func (M *MS6) GetDataIdent() byte {
	return 6
}

type MS7 struct {
	MS5
}

func (M *MS7) GetDataIdent() byte {
	return 7
}

/*---------------*/

type SID struct {
	Flag       uint32       `json:"flag"`
	Additional *OctetString `json:"additional"`
}

func (S *SID) decoder(buf *bytes.Reader) error {
	err := binary.Read(buf, binary.BigEndian, &S.Flag)
	if err != nil {
		return errors.New("decode SID err:" + err.Error())
	}
	S.Additional = &OctetString{}
	err = S.Additional.decoder(buf)
	return err
}

func (S *SID) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, S.Flag)
	if err != nil {
		return nil, errors.New("encode SID err: " + err.Error())
	}
	addArray, err := S.Additional.encoder()
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, addArray)
	if err != nil {
		return nil, errors.New("encode SID err: " + err.Error())
	}
	return buf.Bytes(), nil
}

func (S *SID) DataType() byte {
	return RSDIdent
}

func (S *SID) Value() interface{} {
	return S
}

/*----------------------------------*/

type SIDMAC struct {
	Sid *SID `json:"sid"`
	Mac *MAC `json:"mac"`
}

func (S *SIDMAC) decoder(buf *bytes.Reader) error {
	S.Sid = &SID{}
	err := S.Sid.decoder(buf)
	if err != nil {
		return err
	}
	S.Mac = &MAC{}
	err = S.Mac.decoder(buf)
	return err
}

func (S *SIDMAC) encoder() ([]byte, error) {
	sidArray, err := S.Sid.encoder()
	if err != nil {
		return nil, err
	}
	macArray, err := S.Mac.encoder()
	if err != nil {
		return nil, err
	}
	return append(sidArray, macArray...), nil
}

func (S *SIDMAC) DataType() byte {
	return SIDMACIdent
}

func (S *SIDMAC) Value() interface{} {
	return S
}

/*------------------------------*/

type COMDCB struct {
	Baud        byte `json:"baud"`
	Parity      byte `json:"parity"`
	DataBits    byte `json:"data_bits"`
	StopBits    byte `json:"stop_bits"`
	FlowControl byte `json:"flow_control"`
}

func (C *COMDCB) decoder(buf *bytes.Reader) error {
	return binary.Read(buf, binary.BigEndian, &C)
}

func (C *COMDCB) encoder() ([]byte, error) {
	return []byte{C.Baud, C.Parity, C.DataBits, C.StopBits, C.FlowControl}, nil
}

func (C *COMDCB) DataType() byte {
	return COMDCBIdent
}

func (C *COMDCB) Value() interface{} {
	return C
}

/*--------------------------------------*/

type RCSD struct {
	CSDs []*CSD `xml:"csd"`
}

func (R *RCSD) decoder(buf *bytes.Reader) error {
	var csdLen byte
	err := binary.Read(buf, binary.BigEndian, &csdLen)
	if err != nil {
		return err
	}
	R.CSDs = make([]*CSD, csdLen)
	for i := 0; i < int(csdLen); i++ {
		csd := &CSD{}
		if err := csd.decoder(buf); err != nil {
			return err
		}
		R.CSDs[i] = csd
	}
	return nil
}

func (R *RCSD) encoder() ([]byte, error) {
	if len(R.CSDs) == 0 {
		return []byte{0x00}, nil
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, byte(len(R.CSDs)))
	if err != nil {
		return nil, err
	}
	for _, csd := range R.CSDs {
		csdArray, err := csd.encoder()
		if err != nil {
			return nil, err
		}
		err = binary.Write(buf, binary.LittleEndian, csdArray)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (R *RCSD) DataType() byte {
	return RCSDIdent
}

func (R *RCSD) Value() interface{} {
	return R
}
