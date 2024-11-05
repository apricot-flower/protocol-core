package dlt698

// CreateLogin 登录
// address 逻辑地址
// heartbeat 心跳周期
// ca 客户机地址
func CreateLogin(address string, heartbeat uint16, pid byte, ca byte) ([]byte, error) {
	return linkRequest(address, 0, heartbeat, pid, ca)
}

// CreateHeartBeat 心跳
// address 逻辑地址
// heartbeat 心跳周期
// ca 客户机地址
func CreateHeartBeat(address string, heartbeat uint16, pid byte, ca byte) ([]byte, error) {
	return linkRequest(address, 1, heartbeat, pid, ca)
}

// CreateLogOut 登出
// address 逻辑地址
// heartbeat 心跳周期
// ca 客户机地址
func CreateLogOut(address string, heartbeat uint16, pid byte, ca byte) ([]byte, error) {
	return linkRequest(address, 2, heartbeat, pid, ca)
}

func linkRequest(address string, linkRequestType byte, heartbeatCycle uint16, pid byte, ca byte) ([]byte, error) {
	time := &DateTime{}
	protocolDlt645Model := ProtocolDlt645Model{
		Control: &ControlRegion{"1", "0", "0", "0", "001"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: pid, Data: &LinkRequest{LinkRequestType: linkRequestType, HeartbeatCycle: heartbeatCycle, RequestTime: time.Build()}},
	}
	return protocolDlt645Model.Encoder()
}

// CreateLinkResponse 创建预连接响应的报文
// address 逻辑地址
// CA 客户机地址
// result 结果 0-成功 1-地址重复 2-非法设备 3-容量不足
// timing 时钟可信标志 0-不可信 1-可信
// requestTime LinkRequest中的请求时间
// pid
func CreateLinkResponse(address string, ca byte, result byte, timing byte, requestTime *DateTime, pid byte) ([]byte, error) {
	time := DateTime{}
	protocolDlt645Model := ProtocolDlt645Model{
		Control: &ControlRegion{"0", "0", "0", "0", "001"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: pid, Data: &LinkResponse{Timing: timing, Result: result, RequestTime: requestTime, ReceiveTime: time.Build(), ResponseTime: time.Build()}},
	}
	return protocolDlt645Model.Encoder()
}

// CreateConnectRequest 生成自定义的请求建立应用连接的报文
// address 逻辑地址
// CA 客户机地址
// pid
// connectRequest 请求建立应用连接对象
// timeTag 时间标签
func CreateConnectRequest(address string, ca byte, pid byte, connectRequest *ConnectRequest, timeTag *TimeTag) ([]byte, error) {
	protocolDlt645Model := ProtocolDlt645Model{
		Control: &ControlRegion{"1", "0", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: pid, Data: connectRequest, TimeTag: timeTag},
	}
	return protocolDlt645Model.Encoder()
}

// CreateConnectResponse 创建建立应用连接响应的报文
// address 逻辑地址
// ca 客户机地址
// pid
// connectResponse 应用连接响应对象
// timeTag 时间标签
func CreateConnectResponse(address string, ca byte, pid byte, connectResponse *ConnectResponse, timeTag *TimeTag) ([]byte, error) {
	protocolDlt645Model := ProtocolDlt645Model{
		Control: &ControlRegion{Dir: "0", Prm: "0", Framing: "0", Sc: "0", Func: "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: pid, Data: connectResponse, FollowReport: nil, TimeTag: timeTag},
	}
	return protocolDlt645Model.Encoder()
}

// CreateReleaseRequest 请求断开应用连接
// address 逻辑地址
// ca 客户机地址
// pid
// timeUnit 时间单位
// interval 间隔时间
func CreateReleaseRequest(address string, ca byte, pid byte, timeUnit byte, interval uint16) ([]byte, error) {
	dateTimes := &DateTimes{}
	releaseRequest := &ReleaseRequest{}
	protocolDlt645Model := ProtocolDlt645Model{
		Control: &ControlRegion{"1", "0", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: pid, Data: releaseRequest, TimeTag: &TimeTag{SendTime: dateTimes.Build(), Ti: &TI{TimeUnit: timeUnit, Interval: interval}}},
	}
	return protocolDlt645Model.Encoder()
}

// CreateReleaseResponse 断开应用连接响应
// address 逻辑地址
// ca 客户机地址
// pid
// timeUnit 时间单位
// interval 间隔时间
func CreateReleaseResponse(address string, ca byte, pid byte, timeUnit byte, interval uint16) ([]byte, error) {
	dateTimes := &DateTimes{}
	releaseResponse := &ReleaseResponse{
		Result: 0,
	}
	protocolDlt645Model := ProtocolDlt645Model{
		Control: &ControlRegion{"0", "0", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: pid, Data: releaseResponse, TimeTag: &TimeTag{SendTime: dateTimes.Build(), Ti: &TI{TimeUnit: timeUnit, Interval: interval}}},
	}
	return protocolDlt645Model.Encoder()
}

// CreateReleaseNotification  断开应用连接通知
// address 逻辑地址
// ca 客户机地址
// pid
// 时间标签
func CreateReleaseNotification(address string, ca byte, pid byte, timeTag *TimeTag) ([]byte, error) {
	dateTimes := &DateTimes{}
	releaseNotification := &ReleaseNotification{
		LinkedTime: dateTimes.Build(),
		ServerTime: dateTimes.Build(),
	}
	protocolDlt645Model := ProtocolDlt645Model{
		Control: &ControlRegion{"1", "0", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: pid, Data: releaseNotification, FollowReport: nil, TimeTag: timeTag},
	}
	return protocolDlt645Model.Encoder()
}

// CreateGetRequestNormal 请求读取一个对象属性
// address 逻辑地址
// ca 客户机地址
// pid
// oad 请求标志
// timeTag 时间标签域
func CreateGetRequestNormal(address string, ca byte, pid byte, oad []byte, timeTag *TimeTag) ([]byte, error) {
	getRequestNormal := &GetRequestNormal{
		OAD: oad,
	}
	protocolDlt645Model := ProtocolDlt645Model{
		Control: &ControlRegion{"1", "0", "0", "0", "011"},
		Address: &AddressRegion{AddressType: 0, Address: address, CA: ca},
		Data:    &APDU{Pid: pid, Data: getRequestNormal, TimeTag: timeTag},
	}
	return protocolDlt645Model.Encoder()
}
