package dlt1376_1

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func Decode(frameStr string) (*Dlt13761statute, error) {
	frameStr = strings.ReplaceAll(frameStr, " ", "")
	frameBytes, err := hex.DecodeString(frameStr)
	if err != nil {
		return nil, err
	}
	return DecodeBytes(frameBytes)
}

func DecodeBytes(frame []byte) (*Dlt13761statute, error) {
	statute := &Dlt13761statute{}
	err := statute.Decode(frame)
	if err != nil {
		return nil, err
	}
	return statute, nil
}

// BuildTp 生成一个时间标签
// pfc 启动帧帧序号计数器PFC
// delayed 允许发送传输延时时间
func BuildTp(pfc byte, delayed byte) (*Tp, error) {
	tp := &Tp{PFC: pfc, Delayed: delayed}
	err := tp.BuildByNow()
	return tp, err
}

// Reset 生成复位命令
// address 地址
// PSEQorRSEQ 帧序列号
// hasTp 是否包含时间标签
// fn
// pw 消息认证码字段
func Reset(address string, PSEQorRSEQ byte, tp *Tp, fn uint64, pw []byte) ([]byte, error) {
	control := &ControlFiled{DIR: "0", PRM: "1", FCBorACD: "0", FCV: "1", FuncCode: "0001"}
	addressFiled := &AddressFiled{AddressType: "0", MSA: "0001010"}
	err := addressFiled.Build(address)
	if err != nil {
		return nil, err
	}
	resetCommand := &ResetCommand{}
	_ = resetCommand.Append(&Dlt13761Data{P: 0, F: fn, Data: nil})
	tpFlag := "0"
	if tp != nil {
		tpFlag = "1"
	}
	data := &LinkUserData{Seq: &SEQ{TpV: tpFlag, FIR: "1", FIN: "1", CON: "1", PSEQorRSEQ: fmt.Sprintf("%04b", PSEQorRSEQ&0xF)}, DataUnit: resetCommand, Aux: &AUX{PW: pw, TP: tp}}
	statute := &Dlt13761statute{Control: control, Address: addressFiled, Data: data}
	return statute.Encode()
}

// CreateLinkInterfaceDetection 链路接口检测
// address 地址
// PSEQorRSEQ 帧序列号
func CreateLinkInterfaceDetection(address string, PSEQorRSEQ byte, fn uint64) ([]byte, error) {
	control := &ControlFiled{DIR: "1", PRM: "1", FCBorACD: "0", FCV: "0", FuncCode: "1001"}
	addressFiled := &AddressFiled{AddressType: "0", MSA: "0"}
	err := addressFiled.Build(address)
	if err != nil {
		return nil, err
	}
	linkInterfaceDetection := &LinkInterfaceDetection{}
	err = linkInterfaceDetection.Append(&Dlt13761Data{P: 0, F: fn, Data: nil})
	if err != nil {
		return nil, err
	}
	data := &LinkUserData{Seq: &SEQ{TpV: "0", FIR: "1", FIN: "1", CON: "1", PSEQorRSEQ: fmt.Sprintf("%04b", PSEQorRSEQ&0xF)}, DataUnit: linkInterfaceDetection}
	statute := &Dlt13761statute{Control: control, Address: addressFiled, Data: data}
	return statute.Encode()
}

// CreateRelayStationCommandDownF1 中继命令下行,中继站工作状态控制
// address 地址
// PSEQorRSEQ 帧序列号
// fn
// switchControl 值班机/备份机切换控制
// allowFlag 值班机中继转发允许标志
func CreateRelayStationCommandDownF1(address string, PSEQorRSEQ byte, fn uint64, switchControl byte, allowFlag byte) ([]byte, error) {
	control := &ControlFiled{DIR: "0", PRM: "1", FCBorACD: "0", FCV: "1", FuncCode: "0001"}
	addressFiled := &AddressFiled{AddressType: "0", MSA: "0001010"}
	err := addressFiled.Build(address)
	if err != nil {
		return nil, err
	}
	relayStationCommand := &RelayStationCommand{}
	relayStationCommand.Direction("0")
	err = relayStationCommand.Append(&Dlt13761Data{P: 0, F: fn, Data: &RelayStationCommandDownF1{SwitchControl: switchControl, AllowFlag: allowFlag}})
	if err != nil {
		return nil, err
	}
	data := &LinkUserData{Seq: &SEQ{TpV: "0", FIR: "1", FIN: "1", CON: "1", PSEQorRSEQ: fmt.Sprintf("%04b", PSEQorRSEQ&0xF)}, DataUnit: relayStationCommand}
	statute := &Dlt13761statute{Control: control, Address: addressFiled, Data: data}
	return statute.Encode()
}

//中继命令上行
