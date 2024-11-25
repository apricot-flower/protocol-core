package dlt698

import (
	"errors"
	"fmt"
	"strconv"
)

// CreateLogin 登录
// address 逻辑地址
// heartbeat 心跳周期
// ca 客户机地址
func CreateLogin(address string, heartbeat uint16, Pid byte, ca byte) ([]byte, error) {
	return linkRequest(address, 0, heartbeat, Pid, ca)
}

// CreateHeartBeat 心跳
// address 逻辑地址
// heartbeat 心跳周期
// ca 客户机地址
func CreateHeartBeat(address string, heartbeat uint16, Pid byte, ca byte) ([]byte, error) {
	return linkRequest(address, 1, heartbeat, Pid, ca)
}

// CreateLogOut 登出
// address 逻辑地址
// heartbeat 心跳周期
// ca 客户机地址
func CreateLogOut(address string, heartbeat uint16, Pid byte, ca byte) ([]byte, error) {
	return linkRequest(address, 2, heartbeat, Pid, ca)
}

func linkRequest(address string, linkRequestType byte, heartbeatCycle uint16, Pid byte, ca byte) ([]byte, error) {
	time := &DateTime{}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "0", "0", "0", "001"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: &LinkRequest{LinkRequestType: linkRequestType, HeartbeatCycle: heartbeatCycle, RequestTime: time.Build()}},
	}
	return protocolDlt698Model.Encoder()
}

// CreateLinkResponse 创建预连接响应的报文
// address 逻辑地址
// CA 客户机地址
// result 结果 0-成功 1-地址重复 2-非法设备 3-容量不足
// timing 时钟可信标志 0-不可信 1-可信
// requestTime LinkRequest中的请求时间
// Pid
func CreateLinkResponse(address string, ca byte, result byte, timing byte, requestTime *DateTime, Pid byte) ([]byte, error) {
	time := DateTime{}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"0", "0", "0", "0", "001"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: &LinkResponse{Timing: timing, Result: result, RequestTime: requestTime, ReceiveTime: time.Build(), ResponseTime: time.Build()}},
	}
	return protocolDlt698Model.Encoder()
}

// CreateConnectRequest 生成自定义的请求建立应用连接的报文
// address 逻辑地址
// CA 客户机地址
// Pid
// connectRequest 请求建立应用连接对象
// timeTag 时间标签
func CreateConnectRequest(address string, ca byte, Pid byte, connectRequest *ConnectRequest, timeTag *TimeTag) ([]byte, error) {
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "0", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: connectRequest, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateConnectResponse 创建建立应用连接响应的报文
// address 逻辑地址
// ca 客户机地址
// Pid
// connectResponse 应用连接响应对象
// timeTag 时间标签
func CreateConnectResponse(address string, ca byte, Pid byte, connectResponse *ConnectResponse, timeTag *TimeTag) ([]byte, error) {
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{Dir: "0", Prm: "0", Framing: "0", Sc: "0", Func: "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: connectResponse, FollowReport: nil, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateReleaseRequest 请求断开应用连接
// address 逻辑地址
// ca 客户机地址
// Pid
// timeUnit 时间单位
// interval 间隔时间
func CreateReleaseRequest(address string, ca byte, Pid byte, timeUnit byte, interval uint16) ([]byte, error) {
	dateTimes := &DateTimes{}
	releaseRequest := &ReleaseRequest{}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "0", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: releaseRequest, TimeTag: &TimeTag{SendTime: dateTimes.Build(), Ti: &TI{TimeUnit: timeUnit, Interval: interval}}},
	}
	return protocolDlt698Model.Encoder()
}

// CreateReleaseResponse 断开应用连接响应
// address 逻辑地址
// ca 客户机地址
// Pid
// timeUnit 时间单位
// interval 间隔时间
func CreateReleaseResponse(address string, ca byte, Pid byte, timeUnit byte, interval uint16) ([]byte, error) {
	dateTimes := &DateTimes{}
	releaseResponse := &ReleaseResponse{
		Result: 0,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"0", "0", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: releaseResponse, TimeTag: &TimeTag{SendTime: dateTimes.Build(), Ti: &TI{TimeUnit: timeUnit, Interval: interval}}},
	}
	return protocolDlt698Model.Encoder()
}

// CreateReleaseNotification  断开应用连接通知
// address 逻辑地址
// ca 客户机地址
// Pid
// 时间标签
func CreateReleaseNotification(address string, ca byte, Pid byte, timeTag *TimeTag) ([]byte, error) {
	dateTimes := &DateTimes{}
	releaseNotification := &ReleaseNotification{
		LinkedTime: dateTimes.Build(),
		ServerTime: dateTimes.Build(),
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "0", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: releaseNotification, FollowReport: nil, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateGetRequestNormal 请求读取一个对象属性
// address 逻辑地址
// ca 客户机地址
// Pid
// oad 请求标志
// timeTag 时间标签域
func CreateGetRequestNormal(address string, ca byte, Pid byte, oad []byte, timeTag *TimeTag) ([]byte, error) {
	getRequestNormal := &GetRequestNormal{
		OAD: oad,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "0", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: getRequestNormal, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateGetRequestNormalList 请求读取若干个对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// oad 请求标志
func CreateGetRequestNormalList(address string, ca byte, Pid byte, timeTag *TimeTag, oad ...[]byte) ([]byte, error) {
	for i := 0; i < len(oad); i++ {
		if len(oad[i]) != 4 {
			return nil, errors.New("oad length is not 4")
		}
	}
	getRequestNormalList := &GetRequestNormalList{
		OADs: oad,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "0", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: getRequestNormalList, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateGetRequestRecord 请求读取一个记录型对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// oad
// selector RSD类型
// rcsd
// timeTag 时间标签域
func CreateGetRequestRecord(address string, ca byte, Pid byte, oad []byte, selector Selector, rcsd *RCSD, timeTag *TimeTag) ([]byte, error) {
	if len(oad) != 4 {
		return nil, errors.New("oad length is not 4")
	}
	getRecord := &GetRecord{
		OAD:  oad,
		Rsd:  &RSD{Selector: selector},
		Rcsd: rcsd,
	}
	getRequestRecord := &GetRequestRecord{GetRecord: getRecord}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"0", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: getRequestRecord, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateGetRequestRecordList 请求读取若干个记录型对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
func CreateGetRequestRecordList(address string, ca byte, Pid byte, timeTag *TimeTag, getRecord ...*GetRecord) ([]byte, error) {
	getRequestRecordList := &GetRequestRecordList{GetRecords: getRecord}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"0", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: getRequestRecordList, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateGetRequestNext 读取下一帧数据
// address 服务器地址SA
// ca 客户机地址
// Pid
// lastId 已接收的最后分帧序号
func CreateGetRequestNext(address string, ca byte, Pid byte, lastId uint16) ([]byte, error) {
	getRequestNext := &GetRequestNext{LastId: lastId}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"0", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: getRequestNext},
	}
	return protocolDlt698Model.Encoder()
}

// CreateGetRequestMD5 请求读取一个对象属性MD5值
// address 服务器地址SA
// ca 客户机地址
// Pid
// oad
func CreateGetRequestMD5(address string, ca byte, Pid byte, oad []byte) ([]byte, error) {
	if len(oad) != 4 {
		return nil, errors.New("oad length is not 4")
	}
	getRequestMD5 := &GetRequestMD5{OAD: oad}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"0", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: getRequestMD5},
	}
	return protocolDlt698Model.Encoder()
}

// CreateGetResponseNormal 响应读取一个对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// oad
// 数据
func CreateGetResponseNormal(address string, ca byte, Pid byte, oad []byte, data DataInter, timeTag *TimeTag, followReport *FollowReport) ([]byte, error) {
	if len(oad) != 4 {
		return nil, errors.New("oad length is not 4")
	}
	getResponseNormal := &GetResponseNormal{
		ResultNormal: &ResultNormal{
			OAD: oad,
			GetResult: &GetResult{
				Data: data,
			},
		},
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: getResponseNormal, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// CreateGetResponseNormalList 响应读取若干个对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// resultNormal 结果集
func CreateGetResponseNormalList(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, resultNormal ...*ResultNormal) ([]byte, error) {
	getResponseNormalList := &GetResponseNormalList{
		ResultNormals: resultNormal,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: getResponseNormalList, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// CreateGetResponseRecord 响应读取一个记录型对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// resultRecord 结果集
func CreateGetResponseRecord(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, resultRecord *ResultRecord) ([]byte, error) {
	getResponseRecord := &GetResponseRecord{ResultRecord: resultRecord}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: getResponseRecord, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// CreateGetResponseRecordList 响应读取若干个记录型对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// resultRecord 结果集
func CreateGetResponseRecordList(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, resultRecord ...*ResultRecord) ([]byte, error) {
	getResponseRecordList := &GetResponseRecordList{ResultRecords: resultRecord}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: getResponseRecordList, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// GetResponseNextError 响应读取分帧传输的下一帧 错误信息用
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// endFlag 末帧标志
// dar 错误标志
func GetResponseNextError(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, endFlag byte, frameNumber uint16, dar *DAR) ([]byte, error) {
	gr := &GetResponseNext{
		EndFlag:     endFlag,
		FrameNumber: frameNumber,
		Dar:         dar,
	}
	return getResponseNext(address, ca, Pid, timeTag, followReport, gr)
}

// GetResponseNextResultNormal 响应读取分帧传输的下一帧 ResultNormal
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// endFlag 末帧标志
// rn 对象属性
func GetResponseNextResultNormal(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, endFlag byte, frameNumber uint16, rn ...*ResultNormal) ([]byte, error) {
	dataSlice := make([]DataInter, len(rn))
	for i, item := range rn {
		dataSlice[i] = item
	}
	gr := &GetResponseNext{EndFlag: endFlag, FrameNumber: frameNumber, Data: dataSlice}
	return getResponseNext(address, ca, Pid, timeTag, followReport, gr)
}

// GetResponseNextResultRecord 响应读取分帧传输的下一帧 ResultRecord
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// endFlag 末帧标志
// rn 对象属性
func GetResponseNextResultRecord(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, endFlag byte, frameNumber uint16, rn ...*ResultRecord) ([]byte, error) {
	dataSlice := make([]DataInter, len(rn))
	for i, item := range rn {
		dataSlice[i] = item
	}
	gr := &GetResponseNext{EndFlag: endFlag, FrameNumber: frameNumber, Data: dataSlice}
	return getResponseNext(address, ca, Pid, timeTag, followReport, gr)
}

func getResponseNext(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, gr *GetResponseNext) ([]byte, error) {
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: gr, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// GetResponseMD5Error 响应读取对象属性MD5值 错误响应
// address 服务器地址SA
// ca 客户机地址
// Pid
// errCode 错误标志
func GetResponseMD5Error(address string, ca byte, Pid byte, oad []byte, errCode byte, timeTag *TimeTag, followReport *FollowReport) ([]byte, error) {
	if len(oad) != 4 {
		return nil, errors.New("oad length is not 4")
	}
	getResponseMD5 := &GetResponseMD5{
		Oad:  oad,
		Data: &DAR{Data: errCode},
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: getResponseMD5, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// CreateGetResponseMD5 响应读取对象属性MD5值 正常响应
// address 服务器地址SA
// ca 客户机地址
// Pid
// errCode 错误标志
func CreateGetResponseMD5(address string, ca byte, Pid byte, oad []byte, str string, timeTag *TimeTag, followReport *FollowReport) ([]byte, error) {
	if len(oad) != 4 {
		return nil, errors.New("oad length is not 4")
	}
	getResponseMD5 := &GetResponseMD5{
		Oad:  oad,
		Data: &OctetString{Data: str},
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: getResponseMD5, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// CreateSetRequestNormal 请求设置一个对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// oad
// data 设置的数据
// timeTag 时间标签域
func CreateSetRequestNormal(address string, ca byte, Pid byte, oad []byte, data DataInter, timeTag *TimeTag) ([]byte, error) {
	if len(oad) != 4 {
		return nil, errors.New("oad length is not 4")
	}
	setRequestNormal := &SetRequestNormal{
		Oad:  oad,
		Data: data,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: setRequestNormal, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateSetRequestNormalList 请求设置若干个对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// sn 设置数据集
func CreateSetRequestNormalList(address string, ca byte, Pid byte, timeTag *TimeTag, sn ...*SetRequestNormal) ([]byte, error) {
	setRequestNormalList := &SetRequestNormalList{
		Data: sn,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: setRequestNormalList, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateSetThenGetRequestNormalList 请求设置后读取若干个对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// setItem 设置数据集
func CreateSetThenGetRequestNormalList(address string, ca byte, Pid byte, timeTag *TimeTag, setItem ...*SetThenGetRequestItem) ([]byte, error) {
	setThenGetRequestNormalList := &SetThenGetRequestNormalList{
		Data: setItem,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: setThenGetRequestNormalList, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateSetResponseNormal 响应设置一个对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// oad
// result 设置结果
func CreateSetResponseNormal(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, oad []byte, result byte) ([]byte, error) {
	if len(oad) != 4 {
		return nil, errors.New("oad length is not 4")
	}
	setResponseNormal := &SetResponseNormal{
		Oad: oad,
		Dar: &DAR{Data: result},
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: setResponseNormal, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// CreateSetResponseNormalList 响应设置若干个对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// sr 设置结果集
func CreateSetResponseNormalList(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, sr ...*SetResponseNormal) ([]byte, error) {
	setResponseNormalList := &SetResponseNormalList{
		Data: sr,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: setResponseNormalList, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// CreateSetThenGetResponseNormalListItem 创建SetThenGetResponseNormalList的子对象
// OAD 一个设置的对象属性描述符
// result 设置执行结果
// road ResultNormal中的oad
// data ResultNormal中的数据
func CreateSetThenGetResponseNormalListItem(oad []byte, result byte, road []byte, data DataInter) (*SetThenGetResponseNormalListItem, error) {
	if len(oad) != 4 || len(road) != 4 {
		return nil, errors.New("oad length is not 4")
	}
	setThenGetResponseNormalListItem := &SetThenGetResponseNormalListItem{
		Oad: oad,
		Dar: &DAR{Data: result},
		ResultNormal: &ResultNormal{
			OAD:       road,
			GetResult: &GetResult{Data: data},
		},
	}
	return setThenGetResponseNormalListItem, nil
}

// CreateSetThenGetResponseNormalList 响应设置若干个对象属性以及读取若干个对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// data 结果集
func CreateSetThenGetResponseNormalList(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, data ...*SetThenGetResponseNormalListItem) ([]byte, error) {
	setThenGetResponseNormalList := &SetThenGetResponseNormalList{
		Data: data,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: setThenGetResponseNormalList, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// CreateActionRequestNormal 服务器地址SA 请求操作一个对象方法
// address 服务器地址SA
// ca 客户机地址
// Pid
// oi 对象标识
// funcMark  方法标识
// mode 操作模式
// data 操作值
// timeTag 时间标签域
func CreateActionRequestNormal(address string, ca byte, Pid byte, oi []byte, funcMark byte, mode byte, data DataInter, timeTag *TimeTag) ([]byte, error) {
	if len(oi) != 4 {
		return nil, errors.New("oi size != 4")
	}
	oia1 := fmt.Sprintf("%04b", oi[0]&0x0F) + fmt.Sprintf("%04b", oi[1]&0x0F)
	oia2 := fmt.Sprintf("%04b", oi[2]&0x0F) + fmt.Sprintf("%04b", oi[3]&0x0F)
	result, err := strconv.ParseUint(oia1+oia2, 2, 16)
	if err != nil {
		return nil, err
	}
	actionRequestNormal := &ActionRequestNormal{
		OMD: &OMD{
			Oi:       &OI{Data: uint16(result)},
			FuncMark: funcMark,
			Mode:     mode,
		},
		Data: data,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"0", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: actionRequestNormal, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

func ActionRequestNormalItem(oi []byte, funcMark byte, mode byte, data DataInter) (*ActionRequestNormal, error) {
	if len(oi) != 4 {
		return nil, errors.New("oi size != 4")
	}
	oia1 := fmt.Sprintf("%04b", oi[0]&0x0F) + fmt.Sprintf("%04b", oi[1]&0x0F)
	oia2 := fmt.Sprintf("%04b", oi[2]&0x0F) + fmt.Sprintf("%04b", oi[3]&0x0F)
	result, err := strconv.ParseUint(oia1+oia2, 2, 16)
	if err != nil {
		return nil, err
	}
	actionRequestNormal := &ActionRequestNormal{
		OMD: &OMD{
			Oi:       &OI{Data: uint16(result)},
			FuncMark: funcMark,
			Mode:     mode,
		},
		Data: data,
	}
	return actionRequestNormal, nil
}

// CreateActionRequestNormalList 请求操作若干个对象方法
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签
// data 结果集
func CreateActionRequestNormalList(address string, ca byte, Pid byte, timeTag *TimeTag, data ...*ActionRequestNormal) ([]byte, error) {
	actionRequestNormalLit := &ActionRequestNormalList{
		Data: data,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"0", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: actionRequestNormalLit, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateActionThenGetRequestNormalListItem 创建一个ActionThenGetRequestNormal
// oi 对象标识
// funcMark  方法标识
// mode 操作模式
// data 操作值
// oad 一个读取的对象属性描述符
// daley 读取延时
func CreateActionThenGetRequestNormalListItem(oi []byte, funcMark byte, mode byte, data DataInter, oad []byte, daley uint8) (*ActionThenGetRequestNormalListItem, error) {
	if len(oi) != 4 {
		return nil, errors.New("oi size != 4")
	}
	if len(oad) != 4 {
		return nil, errors.New("oad size != 4")
	}
	oia1 := fmt.Sprintf("%04b", oi[0]&0x0F) + fmt.Sprintf("%04b", oi[1]&0x0F)
	oia2 := fmt.Sprintf("%04b", oi[2]&0x0F) + fmt.Sprintf("%04b", oi[3]&0x0F)
	result, err := strconv.ParseUint(oia1+oia2, 2, 16)
	if err != nil {
		return nil, err
	}
	omd := &OMD{
		Oi:       &OI{Data: uint16(result)},
		FuncMark: funcMark,
		Mode:     mode,
	}
	actionThenGetRequestNormalListItem := &ActionThenGetRequestNormalListItem{
		Omd:   omd,
		Data:  data,
		Oad:   oad,
		Daley: daley,
	}
	return actionThenGetRequestNormalListItem, nil
}

// CreateActionThenGetRequestNormalList 请求操作若干个对象方法后读取若干个对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签
// data 结果集
func CreateActionThenGetRequestNormalList(address string, ca byte, Pid byte, timeTag *TimeTag, data ...*ActionThenGetRequestNormalListItem) ([]byte, error) {
	actionThenGetRequestNormalList := &ActionThenGetRequestNormalList{
		Data: data,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"0", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: actionThenGetRequestNormalList, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateActionResponseNormal 响应操作一个对象方法
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// funcMark  方法标识
// mode 操作模式
// result 操作执行结果
// data 操作值
func CreateActionResponseNormal(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, oi []byte, funcMark byte, mode byte, result byte, data DataInter) ([]byte, error) {
	if len(oi) != 4 {
		return nil, errors.New("oi size != 4")
	}
	oia1 := fmt.Sprintf("%04b", oi[0]&0x0F) + fmt.Sprintf("%04b", oi[1]&0x0F)
	oia2 := fmt.Sprintf("%04b", oi[2]&0x0F) + fmt.Sprintf("%04b", oi[3]&0x0F)
	value, err := strconv.ParseUint(oia1+oia2, 2, 16)
	if err != nil {
		return nil, err
	}
	actionResponseNormal := &ActionResponseNormal{
		Omd: &OMD{
			Oi:       &OI{Data: uint16(value)},
			FuncMark: funcMark,
			Mode:     mode,
		},
		DAR:  &DAR{Data: result},
		Data: data,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: actionResponseNormal, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// CreateActionResponseNormalItem 创建一个ActionResponseNormal
// oi 对象标识
// funcMark  方法标识
// mode 操作模式
// result 操作结果
func CreateActionResponseNormalItem(oi []byte, funcMark byte, mode byte, result byte, data DataInter) (*ActionResponseNormal, error) {
	if len(oi) != 4 {
		return nil, errors.New("oi size != 4")
	}
	oia1 := fmt.Sprintf("%04b", oi[0]&0x0F) + fmt.Sprintf("%04b", oi[1]&0x0F)
	oia2 := fmt.Sprintf("%04b", oi[2]&0x0F) + fmt.Sprintf("%04b", oi[3]&0x0F)
	value, err := strconv.ParseUint(oia1+oia2, 2, 16)
	if err != nil {
		return nil, err
	}
	actionResponseNormal := &ActionResponseNormal{
		Omd: &OMD{
			Oi:       &OI{Data: uint16(value)},
			FuncMark: funcMark,
			Mode:     mode,
		},
		DAR:  &DAR{Data: result},
		Data: data,
	}
	return actionResponseNormal, err
}

// CreateActionResponseNormalList 响应操作若干个对象方法
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// data 结果集
func CreateActionResponseNormalList(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, data ...*ActionResponseNormal) ([]byte, error) {
	actionResponseNormalList := &ActionResponseNormalList{
		Data: data,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: actionResponseNormalList, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// CreateActionThenGetResponseNormalListItem ActionThenGetResponseNormalList中的属性
// oi 对象标识
// funcMark  方法标识
// mode 操作模式
// result 操作执行结果
// data1 操作返回数据
// oad
// data2 resultNormal中的数据
func CreateActionThenGetResponseNormalListItem(oi []byte, funcMark byte, mode byte, result byte, data1 DataInter, oad []byte, data2 DataInter) (*ActionThenGetResponseNormalListItem, error) {
	if len(oi) != 4 {
		return nil, errors.New("oi size != 4")
	}
	if len(oad) != 4 {
		return nil, errors.New("oad size != 4")
	}
	oia1 := fmt.Sprintf("%04b", oi[0]&0x0F) + fmt.Sprintf("%04b", oi[1]&0x0F)
	oia2 := fmt.Sprintf("%04b", oi[2]&0x0F) + fmt.Sprintf("%04b", oi[3]&0x0F)
	value, err := strconv.ParseUint(oia1+oia2, 2, 16)
	if err != nil {
		return nil, err
	}
	omd := &OMD{
		Oi:       &OI{Data: uint16(value)},
		FuncMark: funcMark,
		Mode:     mode,
	}
	dar := &DAR{Data: result}
	resultNormal := &ResultNormal{
		OAD: oad,
		GetResult: &GetResult{
			Data: data2,
		},
	}
	actionThenGetResponseNormalListItem := &ActionThenGetResponseNormalListItem{
		Omd:          omd,
		Dar:          dar,
		Data:         data1,
		ResultNormal: resultNormal,
	}
	return actionThenGetResponseNormalListItem, nil
}

// CreateActionThenGetResponseNormalList 响应操作若干个对象方法后读取若干个属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// data 结果集
func CreateActionThenGetResponseNormalList(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, data ...*ActionThenGetResponseNormalListItem) ([]byte, error) {
	actionThenGetResponseNormalList := &ActionThenGetResponseNormalList{
		Data: data,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: actionThenGetResponseNormalList, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// CreateReportNotificationList 通知上报若干个对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// data 结果集
func CreateReportNotificationList(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, data ...*ResultNormal) ([]byte, error) {
	reportNotificationList := &ReportNotificationList{
		Data: data,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "0", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: reportNotificationList, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// CreateReportNotificationRecordList 通知上报若干个记录型对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// data 结果集
func CreateReportNotificationRecordList(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, data ...*ResultRecord) ([]byte, error) {
	reportNotificationRecordList := &ReportNotificationRecordList{
		Data: data,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "0", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: reportNotificationRecordList, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// CreateReportNotificationTransData 通知上报透明数据
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// oad
// 数据
func CreateReportNotificationTransData(address string, ca byte, Pid byte, timeTag *TimeTag, followReport *FollowReport, oad []byte, value string) ([]byte, error) {
	if len(oad) != 4 {
		return nil, errors.New("oad size != 4")
	}
	reportNotificationTransData := &ReportNotificationTransData{
		Oad:  oad,
		Data: &OctetString{Data: value},
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "0", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: reportNotificationTransData, TimeTag: timeTag, FollowReport: followReport},
	}
	return protocolDlt698Model.Encoder()
}

// CreateReportResponseList 响应上报若干个对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// oad
func CreateReportResponseList(address string, ca byte, Pid byte, timeTag *TimeTag, oad ...[]byte) ([]byte, error) {
	for _, o := range oad {
		if len(o) != 4 {
			return nil, errors.New("oad child size != 4")
		}
	}
	reportResponseList := &ReportResponseList{
		Data: oad,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: reportResponseList, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateReportResponseRecordList 响应上报若干个记录型对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// followReport 跟随上报信息域
// oad
func CreateReportResponseRecordList(address string, ca byte, Pid byte, timeTag *TimeTag, oad ...[]byte) ([]byte, error) {
	for _, o := range oad {
		if len(o) != 4 {
			return nil, errors.New("oad child size != 4")
		}
	}
	reportResponseRecordList := &ReportResponseRecordList{
		ReportResponseList: ReportResponseList{Data: oad},
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: reportResponseRecordList, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateReportResponseTransData 响应上报透明数据
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
func CreateReportResponseTransData(address string, ca byte, Pid byte, timeTag *TimeTag) ([]byte, error) {
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{"1", "1", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: &ReportResponseTransData{}, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateProxyGetRequestList 请求代理读取若干个服务器的若干个对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// timeOut 代理整个请求的超时时间
// data 结果集
func CreateProxyGetRequestList(address string, ca byte, Pid byte, timeTag *TimeTag, timeOut uint16, data ...*ProxyGetRequestListItem) ([]byte, error) {
	proxyGetRequestList := &ProxyGetRequestList{
		TimeOut: &LongUnsigned{Data: timeOut},
		Data:    data,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{Dir: "0", Prm: "1", Framing: "0", Sc: "0", Func: "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: proxyGetRequestList, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}

// CreateProxyGetRequestRecord 请求代理读取一个服务器的一个记录型对象属性
// address 服务器地址SA
// ca 客户机地址
// Pid
// timeTag 时间标签域
// timeOut 代理整个请求的超时时间
// TSA 目标服务器地址
// OAD 对象属性描述符
// RSD 记录行选择描述符
// RCSD 记录列选择描述符
func CreateProxyGetRequestRecord(address string, ca byte, Pid byte, timeTag *TimeTag, timeOut uint16, tsa string, oad []byte, selector Selector, rcsd *RCSD) ([]byte, error) {
	proxyGetRequestRecord := &ProxyGetRequestRecord{
		TimeOut: timeOut,
		Tsa:     tsa,
		Oad:     oad,
		Rsd:     &RSD{Selector: selector},
		Rcsd:    rcsd,
	}
	protocolDlt698Model := ProtocolDlt698Model{
		Control: &ControlRegion{Dir: "0", Prm: "1", Framing: "0", Sc: "0", Func: "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: Pid, Data: proxyGetRequestRecord, TimeTag: timeTag},
	}
	return protocolDlt698Model.Encoder()
}
