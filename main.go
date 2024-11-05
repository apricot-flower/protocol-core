package main

import (
	"encoding/hex"
	"fmt"
	"protocol-core/dlt698"
)

func main() {
	frame := `684200830506050403020100b50402100010ffffffffffffffffffffffffffffffffffffffffffffffff080007ec011f4000001c20000107e80b050f310405001e0f1116`
	dlt698Mode := dlt698.ProtocolDlt645Model{}
	err := dlt698Mode.DecodeByStr(frame)
	if err != nil {
		fmt.Println("dlt698 frame decode error:" + err.Error())
	}
	encode()
}

func encode() {
	fmt.Println("===========ENCODE===========")
	login, err := dlt698.CreateLogOut("010203040506", 300, 0, 0)
	if err != nil {
		fmt.Println("login error:" + err.Error())
	}
	fmt.Println("login:" + hex.EncodeToString(login))
	time := dlt698.DateTime{}
	linkResponse, err := dlt698.CreateLinkResponse("010203040506", 0, 0, 0, time.Build(), 2)
	fmt.Println("linkResponse:" + hex.EncodeToString(linkResponse))
	//应用连接
	connectRequest := &dlt698.ConnectRequest{
		ExpectVersion:           16,
		ProtocolBlock:           []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		FuncBlock:               []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		ClientSendMaxSize:       2048,
		ClientReceiveMaxSize:    2028,
		ClientReceiveWindowSize: 1,
		ClientHandleMaxSize:     8000,
		LinkTimeOut:             7200,
		ConnectMechanismInfo:    &dlt698.NullSecurity{},
	}
	dateTimes := dlt698.DateTimes{}
	tt := &dlt698.TimeTag{SendTime: dateTimes.Build(), Ti: &dlt698.TI{TimeUnit: 5, Interval: 30}}
	connectRequestFrame, err := dlt698.CreateConnectRequest("010203040506", 0, 16, connectRequest, tt)
	if err != nil {
		fmt.Println("connectRequestFrame error:" + err.Error())
	}
	fmt.Println("connectRequest:" + hex.EncodeToString(connectRequestFrame))

	fmt.Println("============================")
}
