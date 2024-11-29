package dlt1376_1

import "fmt"

//确认否认报文示例：68 32 00 32 00 68 0B 03 44 04 00 00 00 61 00 00 01 00 B8 16

// FrameDirection 报文方向
type FrameDirection byte

const (
	DOWN FrameDirection = 0 //下行
	UP   FrameDirection = 1 //上行
)

// CreateConfirmDenyF1OrF2 创建一个确认否认报文, F1：全部确认, F2：全部否认
// address 地址
// direction 报文方向
// FCBorACD
// funcCode 控制域中的功能码
// f F1， F2
// frameNumber 序号
// tp 是否包含时间标签
func CreateConfirmDenyF1OrF2(address string, direction FrameDirection, fCBorACD byte, funcCode byte, f byte, frameNumber byte, tp bool) ([]byte, error) {
	dlt13761Statute, err := build(address, direction, fCBorACD, funcCode)
	if err != nil {
		return nil, err
	}
	if tp {
		dlt13761Statute.Data.Seq = &SEQ{TpV: "1", FIR: "1", FIN: "1", CON: "0", PSEQorRSEQ: fmt.Sprintf("%04b", frameNumber&0x0F)}
		tpv := &Tp{}
		err = tpv.BuildByNow()
		if err != nil {
			return nil, err
		}
		dlt13761Statute.Data.Aux = &AUX{TP: tpv}
	} else {
		dlt13761Statute.Data.Seq = &SEQ{TpV: "0", FIR: "1", FIN: "1", CON: "0", PSEQorRSEQ: fmt.Sprintf("%04b", frameNumber&0x0F)}
	}
	dlt13761Statute.Data.Data = &ConfirmDeny{F: uint64(f)}
	return dlt13761Statute.Encode()
}

// CreateResetDown 创建下行复位命令
// address 地址
// FCBorACD
// funcCode 控制域中的功能码
// f F1， F2
// frameNumber 序号
// pw
// tp 是否包含时间标签
func CreateResetDown(address string, fCBorACD byte, funcCode byte, f byte, frameNumber byte, pw []byte, tp bool) ([]byte, error) {
	dlt13761Statute, err := build(address, DOWN, fCBorACD, funcCode)
	if err != nil {
		return nil, err
	}
	if tp {
		dlt13761Statute.Data.Seq = &SEQ{TpV: "1", FIR: "1", FIN: "1", CON: "0", PSEQorRSEQ: fmt.Sprintf("%04b", frameNumber&0x0F)}
		tpv := &Tp{}
		err = tpv.BuildByNow()
		if err != nil {
			return nil, err
		}
		dlt13761Statute.Data.Aux = &AUX{TP: tpv, PW: pw}
	} else {
		dlt13761Statute.Data.Seq = &SEQ{TpV: "0", FIR: "1", FIN: "1", CON: "0", PSEQorRSEQ: fmt.Sprintf("%04b", frameNumber&0x0F)}
	}
	dlt13761Statute.Data.Data = &Reset{F: uint64(f)}
	return dlt13761Statute.Encode()
}

// CreateLink 创建链路接口检测命令
// address 地址
// f F1， F2, F3
// frameNumber 序号
func CreateLink(address string, f byte, frameNumber byte) ([]byte, error) {
	dlt13761Statute, err := build(address, UP, 0, 9)
	if err != nil {
		return nil, err
	}
	dlt13761Statute.Data.Seq = &SEQ{TpV: "0", FIR: "1", FIN: "1", CON: "1", PSEQorRSEQ: fmt.Sprintf("%04b", frameNumber&0x0F)}
	dlt13761Statute.Data.Data = &LinkInterfaceDetection{F: uint64(f)}
	return dlt13761Statute.Encode()
}

// CreateRelayStationCommandDownF1 创建下行中继站命令-中继站工作状态控制
// address 地址
// frameNumber 序号
// MSSwitchControl D0～D1 值班机/备份机切换控制：D0=0、D1=0：表示不切换；D0=1、D1=1：表示切换；D0、D1为其他是无效
// DutyForwardFlag D2～D3 值班机中继转发允许标志：D2=0、D3=0：表示不允许；D2=1、D3=1：表示允许；D2、D3为其他是无效
func CreateRelayStationCommandDownF1(address string, frameNumber byte, mSSwitchControl string, dutyForwardFlag string) ([]byte, error) {
	dlt13761Statute, err := build(address, DOWN, 0, 4)
	if err != nil {
		return nil, err
	}
	dlt13761Statute.Data.Seq = &SEQ{TpV: "0", FIR: "1", FIN: "1", CON: "1", PSEQorRSEQ: fmt.Sprintf("%04b", frameNumber&0x0F)}
	dlt13761Statute.Data.Data = &RelayStationCommand{F: 1, Cmd: &ControlOfWorkingStatusOfRelayStation{MSSwitchControl: mSSwitchControl, DutyForwardFlag: dutyForwardFlag}}
	return dlt13761Statute.Encode()
}

// CreateRelayStationCommandDownF234 创建下行中继站命令-f2,f3,f4
// address 地址
// frameNumber 序号
// f
func CreateRelayStationCommandDownF234(address string, frameNumber byte, f uint64) ([]byte, error) {
	dlt13761Statute, err := build(address, DOWN, 0, 4)
	if err != nil {
		return nil, err
	}
	dlt13761Statute.Data.Seq = &SEQ{TpV: "0", FIR: "1", FIN: "1", CON: "1", PSEQorRSEQ: fmt.Sprintf("%04b", frameNumber&0x0F)}
	dlt13761Statute.Data.Data = &RelayStationCommand{F: f}
	return dlt13761Statute.Encode()
}

// 构建并不完整的地址
func build(address string, direction FrameDirection, fCBorACD byte, funcCode byte) (*StatuteDlt13761, error) {
	dlt13761Statute := &StatuteDlt13761{}
	dlt13761Statute.Address = &AddressFiled{AddressType: "0", MSA: "0001010"}
	err := dlt13761Statute.Address.Build(address)
	if err != nil {
		return nil, err
	}
	FCBorACD := "1"
	if fCBorACD == 0 {
		FCBorACD = "0"
	}
	if direction == DOWN {
		dlt13761Statute.Control = &ControlFiled{DIR: "0", PRM: "0", FCBorACD: FCBorACD, FCV: "0", FuncCode: fmt.Sprintf("%04b", funcCode&0x0F)}
	} else {
		dlt13761Statute.Control = &ControlFiled{DIR: "1", PRM: "1", FCBorACD: FCBorACD, FCV: "0", FuncCode: fmt.Sprintf("%04b", funcCode&0x0F)}
	}
	dlt13761Statute.Data = &ApplicationLayer{}
	return dlt13761Statute, nil
}
