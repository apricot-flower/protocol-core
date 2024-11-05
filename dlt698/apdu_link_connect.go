package dlt698

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var _ FrameRegion = (*FactoryVersion)(nil)
var _ FrameRegion = (*ConnectResponseInfo)(nil)
var _ FrameRegion = (*SecurityData)(nil)

var _ ConnectMechanismInfo = (*NullSecurity)(nil)
var _ ConnectMechanismInfo = (*PasswordSecurity)(nil)
var _ ConnectMechanismInfo = (*SymmetrySecurity)(nil)
var _ ConnectMechanismInfo = (*SignatureSecurity)(nil)

var _ APDURegion = (*LinkRequest)(nil)
var _ APDURegion = (*LinkResponse)(nil)
var _ APDURegion = (*ConnectRequest)(nil)
var _ APDURegion = (*ConnectResponse)(nil)
var _ APDURegion = (*ReleaseRequest)(nil)
var _ APDURegion = (*ReleaseResponse)(nil)
var _ APDURegion = (*ReleaseNotification)(nil)

const (
	LinkRequestIdent         string = "01"
	LinkResponseIdent        string = "81"
	ConnectRequestIdent      string = "02"
	ConnectResponseIdent     string = "82"
	ReleaseRequestIdent      string = "03"
	ReleaseResponseIdent     string = "83"
	ReleaseNotificationIdent string = "84"
)

func init() {
	apduMap[LinkRequestIdent] = func() APDURegion {
		return new(LinkRequest)
	}
	apduMap[LinkResponseIdent] = func() APDURegion {
		return new(LinkResponse)
	}
	apduMap[ConnectRequestIdent] = func() APDURegion {
		return new(ConnectRequest)
	}
	apduMap[ConnectResponseIdent] = func() APDURegion {
		return new(ConnectResponse)
	}
	apduMap[ReleaseRequestIdent] = func() APDURegion {
		return new(ReleaseRequest)
	}
	apduMap[ReleaseResponseIdent] = func() APDURegion {
		return new(ReleaseResponse)
	}
	apduMap[ReleaseNotificationIdent] = func() APDURegion {
		return new(ReleaseNotification)
	}
}

type LinkRequest struct {
	LinkRequestType byte      `json:"link_request_type"` //请求类型
	HeartbeatCycle  uint16    `json:"heartbeat_cycle"`   //心跳周期
	RequestTime     *DateTime `json:"request_time"`      //请求时间
}

func (l *LinkRequest) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, &l.LinkRequestType); err != nil {
		return err
	}
	if l.LinkRequestType != 0x00 && l.LinkRequestType != 0x01 && l.LinkRequestType != 0x02 {
		return errors.New("link_request_type must be 0x00(login) or 0x01(heartbeat) or 0x02(logout)")
	}
	if err := binary.Read(buf, binary.BigEndian, &l.HeartbeatCycle); err != nil {
		return err
	}
	l.RequestTime = &DateTime{}
	err := l.RequestTime.decoder(buf)
	return err
}

func (l *LinkRequest) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, l.LinkRequestType)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, l.HeartbeatCycle)
	if err != nil {
		return nil, err
	}
	rtArray, err := l.RequestTime.encoder()
	err = binary.Write(buf, binary.BigEndian, rtArray)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (l *LinkRequest) APDUType() string {
	return LinkRequestIdent
}

func (l *LinkRequest) APDUMark() string {
	return "link_request"
}

func (l *LinkRequest) hasFollowReport() bool {
	return false
}

func (l *LinkRequest) hasTimeTag() bool {
	return false
}

/*---------------------------------------------------linkResponse-----------------------------------------------------*/

type LinkResponse struct {
	Timing       byte      `json:"timing"`        //时钟是否可信
	Result       byte      `json:"result"`        //结果
	RequestTime  *DateTime `json:"request_time"`  //请求时间
	ReceiveTime  *DateTime `json:"receive_time"`  //收到时间
	ResponseTime *DateTime `json:"response_time"` //响应
}

func (l *LinkResponse) decoder(buf *bytes.Reader) error {
	var resultValue byte
	if err := binary.Read(buf, binary.BigEndian, &resultValue); err != nil {
		return errors.New("decode linkResponse's Timing err:" + err.Error())
	}
	//获取bit7
	l.Timing = resultValue & (1 << 7)
	//获取bit0~bit2
	l.Result = resultValue & 0x07
	l.RequestTime = &DateTime{}
	if err := l.RequestTime.decoder(buf); err != nil {
		return err
	}
	l.ReceiveTime = &DateTime{}
	if err := l.ReceiveTime.decoder(buf); err != nil {
		return err
	}
	l.ResponseTime = &DateTime{}
	if err := l.ResponseTime.decoder(buf); err != nil {
		return err
	}
	return nil
}

func (l *LinkResponse) encoder() ([]byte, error) {
	var result byte = 0
	result |= l.Timing << 7
	result |= l.Result & 0x07
	encodeArray := []byte{result}
	requestTimeArray, err := l.RequestTime.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, requestTimeArray...)
	receiveTime, err := l.ReceiveTime.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, receiveTime...)
	responseTimeArray, err := l.ResponseTime.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray = append(encodeArray, responseTimeArray...)
	return encodeArray, nil
}

func (l *LinkResponse) APDUType() string {
	return LinkResponseIdent
}

func (l *LinkResponse) APDUMark() string {
	return "link_response"
}

func (l *LinkResponse) hasFollowReport() bool {
	return false
}

func (l *LinkResponse) hasTimeTag() bool {
	return false
}

/*--------------------------------CONNECT-Request---------------------------------*/

// ConnectRequest 请求建立应用连接
type ConnectRequest struct {
	ExpectVersion           uint16               `json:"expect_version"`             //期望的协议版本号
	ProtocolBlock           []byte               `json:"protocol_block"`             //期望的协议一致性块
	FuncBlock               []byte               `json:"func_block"`                 //期望的功能一致性块
	ClientSendMaxSize       uint16               `json:"client_send_max_size"`       //客户机发送帧最大尺寸
	ClientReceiveMaxSize    uint16               `json:"client_receive_max_size"`    //客户机接收帧的最大尺寸
	ClientReceiveWindowSize uint8                `json:"client_receive_window_size"` //客户机接收帧最大窗口尺寸
	ClientHandleMaxSize     uint16               `json:"client_handle_max_size"`     //客户机最大可处理帧尺寸
	LinkTimeOut             uint32               `json:"link_time_out"`              //期望的应用连接超时时间
	ConnectMechanismInfo    ConnectMechanismInfo `json:"connect_mechanism_info"`     //认证请求对象
}

func (c *ConnectRequest) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, &c.ExpectVersion); err != nil {
		return errors.New("decode ConnectRequest's ExpectVersion err:" + err.Error())
	}
	c.ProtocolBlock = make([]byte, 8)
	if err := binary.Read(buf, binary.BigEndian, &c.ProtocolBlock); err != nil {
		return errors.New("decode ConnectRequest's ProtocolBlock err:" + err.Error())
	}
	c.FuncBlock = make([]byte, 16)
	if err := binary.Read(buf, binary.BigEndian, &c.FuncBlock); err != nil {
		return errors.New("decode ConnectRequest's FuncBlock err:" + err.Error())
	}
	if err := binary.Read(buf, binary.BigEndian, &c.ClientSendMaxSize); err != nil {
		return errors.New("decode ConnectRequest's ClientSendMaxSize err:" + err.Error())
	}
	if err := binary.Read(buf, binary.BigEndian, &c.ClientReceiveMaxSize); err != nil {
		return errors.New("decode ConnectRequest's ClientReceiveMaxSize err:" + err.Error())
	}
	if err := binary.Read(buf, binary.BigEndian, &c.ClientReceiveWindowSize); err != nil {
		return errors.New("decode ConnectRequest's ClientReceiveWindowSize err:" + err.Error())
	}
	if err := binary.Read(buf, binary.BigEndian, &c.ClientHandleMaxSize); err != nil {
		return errors.New("decode ConnectRequest's ClientHandleMaxSize err:" + err.Error())
	}
	if err := binary.Read(buf, binary.BigEndian, &c.LinkTimeOut); err != nil {
		return errors.New("decode ConnectRequest's LinkTimeOut err:" + err.Error())
	}
	var connectMechanismInfoType byte
	if err := binary.Read(buf, binary.BigEndian, &connectMechanismInfoType); err != nil {
		return errors.New("decode ConnectRequest's connectMechanismInfoType err:" + err.Error())
	}
	switch connectMechanismInfoType {
	case 0:
		c.ConnectMechanismInfo = &NullSecurity{}
	case 1:
		c.ConnectMechanismInfo = &PasswordSecurity{}
	case 2:
		c.ConnectMechanismInfo = &SignatureSecurity{}
	case 3:
		c.ConnectMechanismInfo = &SymmetrySecurity{}
	}
	if c.ConnectMechanismInfo == nil {
		return errors.New("ConnectMechanismInfo no such type")
	}
	return c.ConnectMechanismInfo.decoder(buf)
}

func (c *ConnectRequest) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &c.ExpectVersion); err != nil {
		return nil, errors.New("encode connectRequest's ExpectVersion err:" + err.Error())
	}
	if err := binary.Write(buf, binary.BigEndian, &c.ProtocolBlock); err != nil {
		return nil, errors.New("encode connectRequest's ProtocolBlock err:" + err.Error())
	}
	if err := binary.Write(buf, binary.BigEndian, &c.FuncBlock); err != nil {
		return nil, errors.New("encode connectRequest's FuncBlock err:" + err.Error())
	}
	if err := binary.Write(buf, binary.BigEndian, &c.ClientSendMaxSize); err != nil {
		return nil, errors.New("encode connectRequest's ClientSendMaxSize err:" + err.Error())
	}
	if err := binary.Write(buf, binary.BigEndian, &c.ClientReceiveMaxSize); err != nil {
		return nil, errors.New("encode connectRequest's ClientReceiveMaxSize err:" + err.Error())
	}
	if err := binary.Write(buf, binary.BigEndian, &c.ClientReceiveWindowSize); err != nil {
		return nil, errors.New("encode connectRequest's ClientReceiveWindowSize err:" + err.Error())
	}
	if err := binary.Write(buf, binary.BigEndian, &c.ClientHandleMaxSize); err != nil {
		return nil, errors.New("encode connectRequest's ClientHandleMaxSize err:" + err.Error())
	}
	if err := binary.Write(buf, binary.BigEndian, &c.LinkTimeOut); err != nil {
		return nil, errors.New("encode connectRequest's LinkTimeOut err:" + err.Error())
	}
	cmiArray, err := c.ConnectMechanismInfo.encoder()
	if err != nil {
		return nil, err
	}
	cmiArray = append([]byte{c.ConnectMechanismInfo.Sign()}, cmiArray...)
	return append(buf.Bytes(), cmiArray...), nil
}

func (c *ConnectRequest) APDUType() string {
	return ConnectRequestIdent
}

func (c *ConnectRequest) APDUMark() string {
	return "connect_request"
}

func (c *ConnectRequest) hasFollowReport() bool {
	return false
}

func (c *ConnectRequest) hasTimeTag() bool {
	return true
}

type ConnectMechanismInfo interface {
	FrameRegion
	Sign() byte
}

type NullSecurity struct {
}

func (n *NullSecurity) decoder(_ *bytes.Reader) error {
	return nil
}

func (n *NullSecurity) encoder() ([]byte, error) {
	return nil, nil
}

func (n *NullSecurity) Sign() byte {
	return 0
}

type PasswordSecurity struct {
	Password *VisibleString `json:"password"` //密码
}

func (p *PasswordSecurity) decoder(buf *bytes.Reader) error {
	p.Password = &VisibleString{}
	return p.Password.decoder(buf)
}

func (p *PasswordSecurity) encoder() ([]byte, error) {
	return p.Password.encoder()
}

func (p *PasswordSecurity) Sign() byte {
	return 1
}

type SymmetrySecurity struct {
	Secret    *OctetString `json:"secret1"`    //密文1
	Signature *OctetString `json:"signature1"` //客户机签名1
}

func (s SymmetrySecurity) decoder(buf *bytes.Reader) error {
	s.Secret = &OctetString{}
	if err := s.Secret.decoder(buf); err != nil {
		return err
	}
	s.Signature = &OctetString{}
	if err := s.Signature.decoder(buf); err != nil {
		return err
	}
	return nil
}

func (s SymmetrySecurity) encoder() ([]byte, error) {
	secret1Array, err := s.Secret.encoder()
	if err != nil {
		return nil, err
	}
	signature1Array, err := s.Signature.encoder()
	if err != nil {
		return nil, err
	}
	return append(secret1Array, signature1Array...), nil
}

func (s SymmetrySecurity) Sign() byte {
	return 2
}

type SignatureSecurity struct {
	Secret    *OctetString `json:"secret2"`    //密文2
	Signature *OctetString `json:"signature2"` //客户机签名2
}

func (s SignatureSecurity) decoder(buf *bytes.Reader) error {
	s.Secret = &OctetString{}
	if err := s.Secret.decoder(buf); err != nil {
		return err
	}
	s.Signature = &OctetString{}
	if err := s.Signature.decoder(buf); err != nil {
		return err
	}
	return nil
}

func (s SignatureSecurity) encoder() ([]byte, error) {
	secret1Array, err := s.Secret.encoder()
	if err != nil {
		return nil, err
	}
	signature1Array, err := s.Signature.encoder()
	if err != nil {
		return nil, err
	}
	return append(secret1Array, signature1Array...), nil
}

func (s SignatureSecurity) Sign() byte {
	return 3
}

/*----------------------CONNECT-Response--------------------*/

type ConnectResponse struct {
	FactoryVersion             *FactoryVersion      `json:"factory_version"`                //服务器厂商版本信息
	ExpectVersion              uint16               `json:"expect_version"`                 //商定的协议版本号
	ProtocolBlock              []byte               `json:"protocol_block"`                 //商定的协议一致性块
	FuncBlock                  []byte               `json:"func_block"`                     //商定的功能一致性块
	ServerSendMaxSize          uint16               `json:"server_send_max_size"`           //服务器发送帧最大尺寸
	ServerReceiveMaxSize       uint16               `json:"server_receive_max_size"`        //服务器接收帧最大尺寸
	ServerReceiveWindowMaxSize uint8                `json:"server_receive_window_max_size"` //服务器接收帧最大窗口尺寸
	ServerHandleMaxSize        uint16               `json:"server_handle_max_size"`         //服务器最大可处理帧尺寸
	LinkTimeOut                uint32               `json:"link_time_out"`                  //商定的应用连接超时时间
	ConnectResponseInfo        *ConnectResponseInfo `json:"connect_response_info"`          //连接响应对象

}

func (c *ConnectResponse) decoder(buf *bytes.Reader) error {
	c.FactoryVersion = &FactoryVersion{}
	if err := c.FactoryVersion.decoder(buf); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.ExpectVersion); err != nil {
		return err
	}
	c.ProtocolBlock = make([]byte, 8)
	if err := binary.Read(buf, binary.BigEndian, &c.ProtocolBlock); err != nil {
		return err
	}
	c.FuncBlock = make([]byte, 16)
	if err := binary.Read(buf, binary.BigEndian, &c.FuncBlock); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.ServerSendMaxSize); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.ServerReceiveMaxSize); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.ServerReceiveWindowMaxSize); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.ServerHandleMaxSize); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.LinkTimeOut); err != nil {
		return err
	}
	c.ConnectResponseInfo = &ConnectResponseInfo{}
	return c.ConnectResponseInfo.decoder(buf)
}

func (c *ConnectResponse) encoder() ([]byte, error) {
	buf := new(bytes.Buffer)
	factoryVersionArray, err := c.FactoryVersion.encoder()
	if err != nil {
		return nil, errors.New("encode connectResponse's factoryVersionArray err :" + err.Error())
	}
	err = binary.Write(buf, binary.BigEndian, factoryVersionArray)
	if err != nil {
		return nil, errors.New("encode connectResponse's factoryVersionArray err :" + err.Error())
	}
	err = binary.Write(buf, binary.BigEndian, c.ExpectVersion)
	if err != nil {
		return nil, errors.New("encode connectResponse's expectVersion err :" + err.Error())
	}
	err = binary.Write(buf, binary.BigEndian, c.ProtocolBlock)
	if err != nil {
		return nil, errors.New("encode connectResponse's protocolBlock err :" + err.Error())
	}
	err = binary.Write(buf, binary.BigEndian, c.FuncBlock)
	if err != nil {
		return nil, errors.New("encode connectResponse's funcBlock err :" + err.Error())
	}
	err = binary.Write(buf, binary.BigEndian, c.ServerSendMaxSize)
	if err != nil {
		return nil, errors.New("encode connectResponse's serverSendMaxSize err :" + err.Error())
	}
	err = binary.Write(buf, binary.BigEndian, c.ServerReceiveMaxSize)
	if err != nil {
		return nil, errors.New("encode connectResponse's serverReceiveMaxSize err :" + err.Error())
	}
	err = binary.Write(buf, binary.BigEndian, c.ServerReceiveWindowMaxSize)
	if err != nil {
		return nil, errors.New("encode connectResponse's serverReceiveWindowMaxSize err :" + err.Error())
	}
	err = binary.Write(buf, binary.BigEndian, c.ServerHandleMaxSize)
	if err != nil {
		return nil, errors.New("encode connectResponse's serverHandleMaxSize err :" + err.Error())
	}
	err = binary.Write(buf, binary.BigEndian, c.LinkTimeOut)
	if err != nil {
		return nil, errors.New("encode connectResponse's linkTimeOut err :" + err.Error())
	}
	connectResponseInfoArray, err := c.ConnectResponseInfo.encoder()
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, connectResponseInfoArray)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *ConnectResponse) APDUType() string {
	return ConnectResponseIdent
}

func (c *ConnectResponse) APDUMark() string {
	return "connect_response"
}

func (c *ConnectResponse) hasFollowReport() bool {
	return true
}

func (c *ConnectResponse) hasTimeTag() bool {
	return true
}

type FactoryVersion struct {
	ManuCode            *VisibleString `json:"manu_code"`             //厂商代码
	SoftwareVersion     *VisibleString `json:"software_version"`      //软件版本号
	SoftwareVersionDate *VisibleString `json:"software_version_date"` //厂家版本日期
	HardwareVersion     *VisibleString `json:"hardware_version"`      //硬件版本号
	HardwareVersionDate *VisibleString `json:"hardware_version_date"` //硬件版本日期
	ExtendedInfo        *VisibleString `json:"extended_info"`         //厂商扩展信息
}

func (f *FactoryVersion) decoder(buf *bytes.Reader) error {
	f.ManuCode = &VisibleString{}
	if err := f.ManuCode.DecodeByLen(buf, 4); err != nil {
		return err
	}
	f.SoftwareVersion = &VisibleString{}
	if err := f.SoftwareVersion.DecodeByLen(buf, 4); err != nil {
		return err
	}
	f.SoftwareVersionDate = &VisibleString{}
	if err := f.SoftwareVersionDate.DecodeByLen(buf, 6); err != nil {
		return err
	}
	f.HardwareVersion = &VisibleString{}
	if err := f.HardwareVersion.DecodeByLen(buf, 4); err != nil {
		return err
	}
	f.HardwareVersionDate = &VisibleString{}
	if err := f.HardwareVersionDate.DecodeByLen(buf, 6); err != nil {
		return err
	}
	f.ExtendedInfo = &VisibleString{}
	if err := f.ExtendedInfo.DecodeByLen(buf, 8); err != nil {
		return err
	}
	return nil
}

func (f *FactoryVersion) encoder() ([]byte, error) {
	manuCodeArray, _ := f.ManuCode.encoder()
	softwareVersionArray, _ := f.SoftwareVersion.encoder()
	softwareVersionDateArray, _ := f.SoftwareVersionDate.encoder()
	hardwareVersionArray, _ := f.HardwareVersion.encoder()
	hardwareVersionDateArray, _ := f.HardwareVersionDate.encoder()
	extendedInfoArray, _ := f.ExtendedInfo.encoder()
	encodeArray := append(manuCodeArray[1:], softwareVersionArray[1:]...)
	encodeArray = append(encodeArray, softwareVersionDateArray[1:]...)
	encodeArray = append(encodeArray, hardwareVersionArray[1:]...)
	encodeArray = append(encodeArray, hardwareVersionDateArray[1:]...)
	encodeArray = append(encodeArray, extendedInfoArray[1:]...)
	return encodeArray, nil
}

// ConnectResponseInfo 连接响应对象
type ConnectResponseInfo struct {
	ConnectResult byte          `json:"connect_result"` //应用连接请求认证的结果
	SecurityData  *SecurityData `json:"security_data"`  //认证附加信息
}

func (c *ConnectResponseInfo) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, &c.ConnectResult); err != nil {
		return errors.New("decode connectResponseInfo's ConnectResult err:" + err.Error())
	}
	var securityDataFlg byte
	if err := binary.Read(buf, binary.BigEndian, &securityDataFlg); err != nil {
		return errors.New("decode connectResponseInfo's SecurityDataFlg err:" + err.Error())
	}

	if securityDataFlg == 0x01 {
		c.SecurityData = &SecurityData{}
		err := c.SecurityData.decoder(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ConnectResponseInfo) encoder() ([]byte, error) {
	if &c.SecurityData == nil {
		return []byte{c.ConnectResult, 0x00}, nil
	}
	encodeSignArray := append([]byte{c.ConnectResult}, 0x01)
	serverSignArray, err := c.SecurityData.encoder()
	if err != nil {
		return nil, err
	}
	encodeSignArray = append(encodeSignArray, serverSignArray...)
	return encodeSignArray, nil
}

type SecurityData struct {
	RN         *OctetString `json:"rn"`          //服务器随机数
	ServerSign *OctetString `json:"server_sign"` //服务器签名
}

func (s *SecurityData) decoder(buf *bytes.Reader) error {
	s.RN = &OctetString{}
	err := s.RN.decoder(buf)
	if err != nil {
		return err
	}
	s.ServerSign = &OctetString{}
	err = s.ServerSign.decoder(buf)
	return err
}

func (s *SecurityData) encoder() ([]byte, error) {
	rnArray, err := s.RN.encoder()
	if err != nil {
		return nil, err
	}
	serverSignArray, err := s.ServerSign.encoder()
	if err != nil {
		return nil, err
	}
	return append(rnArray, serverSignArray...), nil
}

/*------------------RELEASE-Request---------------*/

// ReleaseRequest 请求断开应用连接
type ReleaseRequest struct {
}

func (r *ReleaseRequest) decoder(_ *bytes.Reader) error {
	return nil
}

func (r *ReleaseRequest) encoder() ([]byte, error) {
	return nil, nil
}

func (r *ReleaseRequest) APDUType() string {
	return ReleaseRequestIdent
}

func (r *ReleaseRequest) APDUMark() string {
	return "release_request"
}

func (r *ReleaseRequest) hasFollowReport() bool {
	return false
}

func (r *ReleaseRequest) hasTimeTag() bool {
	return true
}

/*------------------RELEASE-Response---------------*/

type ReleaseResponse struct {
	Result byte `json:"result"` //结果
}

func (r *ReleaseResponse) decoder(buf *bytes.Reader) error {
	if err := binary.Read(buf, binary.BigEndian, &r.Result); err != nil {
		return errors.New("decode ReleaseResponse err:" + err.Error())
	}
	return nil
}

func (r *ReleaseResponse) encoder() ([]byte, error) {
	return []byte{r.Result}, nil
}

func (r *ReleaseResponse) APDUType() string {
	return ReleaseResponseIdent
}

func (r *ReleaseResponse) APDUMark() string {
	return "release_response"
}

func (r *ReleaseResponse) hasFollowReport() bool {
	return true
}

func (r *ReleaseResponse) hasTimeTag() bool {
	return true
}

// ReleaseNotification 断开应用连接通知
type ReleaseNotification struct {
	LinkedTime *DateTimes `json:"linked_time"` //应用连接建立时间
	ServerTime *DateTimes `json:"server_time"` //服务器当前时间
}

func (r *ReleaseNotification) decoder(buf *bytes.Reader) error {
	r.LinkedTime = &DateTimes{}
	if err := r.LinkedTime.decoder(buf); err != nil {
		return err
	}
	r.ServerTime = &DateTimes{}
	if err := r.ServerTime.decoder(buf); err != nil {
		return err
	}
	return nil
}

func (r *ReleaseNotification) encoder() ([]byte, error) {
	linkedTimeArray, err := r.LinkedTime.encoder()
	if err != nil {
		return nil, err
	}
	serverTimeArray, err := r.ServerTime.encoder()
	if err != nil {
		return nil, err
	}
	encodeArray := append(linkedTimeArray, serverTimeArray...)
	return encodeArray, nil
}

func (r *ReleaseNotification) APDUType() string {
	return ReleaseNotificationIdent
}

func (r *ReleaseNotification) APDUMark() string {
	return "release_notification"
}

func (r *ReleaseNotification) hasFollowReport() bool {
	return true
}

func (r *ReleaseNotification) hasTimeTag() bool {
	return true
}
