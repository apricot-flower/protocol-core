package dlt645

import (
	"encoding/hex"
	"errors"
	"strings"
)

// DecodeBytes 解码字节数组
func DecodeBytes(frame []byte) {
	d := &Dlt645statute{}
	err := d.Decode(frame)
	if err != nil {
		if FrameDecodeErrorHandle != nil {
			FrameDecodeErrorHandle(hex.EncodeToString(frame), "", 0, err)
		}
		return
	}
	ControlOverTurn(hex.EncodeToString(frame), d.Address, d.Control, d.Ident, d.Data)
}

// DecodeString 解码字符串
func DecodeString(frame string) {
	frameString := strings.ReplaceAll(frame, " ", "")
	frameBytes, err := hex.DecodeString(frameString)
	if err != nil {
		if FrameDecodeErrorHandle != nil {
			FrameDecodeErrorHandle(frame, "", 0, err)
		}
		return
	}
	DecodeBytes(frameBytes)
}

// MasterReadMeterData1 请求读电能表数据  帧格式 1
// address 电能表地址
// ident 要去取的数据标识
func MasterReadMeterData1(address string, ident string) ([]byte, error) {
	dlt645statute := Dlt645statute{Address: address, Control: MasterReadMeterData, Ident: ident}
	return dlt645statute.Encode()

}

// MasterReadMeterData2 请求读电能表数据  帧格式 2
// address 电能表地址
// ident 要去取的数据标识
// loadNumber 负荷记录块数
func MasterReadMeterData2(address string, ident string, loadNumber byte) ([]byte, error) {
	dlt645statute := Dlt645statute{Address: address, Control: MasterReadMeterData, Ident: ident, Data: &ByteArrayAnalyzer{Value: []byte{loadNumber}}}
	return dlt645statute.Encode()

}

// MasterReadMeterData3 请求读电能表数据  帧格式 3
// address 电能表地址
// ident 要去取的数据标识
// loadNumber 负荷记录块数
// minute 分钟
// hour 时
// day 日
// month 月
// year 年
func MasterReadMeterData3(address string, ident string, loadNumber byte, minute byte, hour byte, day byte, month byte, year byte) ([]byte, error) {
	dlt645statute := Dlt645statute{Address: address, Control: MasterReadMeterData, Ident: ident, Data: &ByteArrayAnalyzer{Value: []byte{loadNumber, minute, hour, day, month, year}}}
	return dlt645statute.Encode()
}

// MasterReadMeterDataNormalResponseFloat 请求读电能表数据的 从站正常应答
// address 电能表地址
// framing true有后续数据， false无后续数据
// ident 数据标识
// value 数据值 原始浮点数
func MasterReadMeterDataNormalResponseFloat(address string, framing bool, ident string, value ...float64) ([]byte, error) {
	floatAnalyzer := Translate(ident)
	if floatAnalyzer == nil {
		return nil, errors.New(ident + "never enroll")
	}
	if analyzer, ok := floatAnalyzer.(*FloatAnalyser); ok {
		analyzer.Value = value
	} else {
		return nil, errors.New(ident + "value type error")
	}
	control := MasterReadMeterNormalResponseNoNext
	if framing {
		control = MasterReadMeterNormalResponseNext
	}
	defer Put(ident, floatAnalyzer)
	dlt645statute := Dlt645statute{Address: address, Control: control, Ident: ident, Data: floatAnalyzer}
	return dlt645statute.Encode()
}

// CreateMasterReadMeterDataErrorResponse 请求读电能表数据的 从站异常应答
// address 电能表数据
// errValue 错误编码
func CreateMasterReadMeterDataErrorResponse(address string, errValue byte) ([]byte, error) {
	dlt645statute := Dlt645statute{Address: address, Control: MasterReadMeterDataErrorResponse, Data: &ByteArrayAnalyzer{Value: []byte{errValue}}}
	return dlt645statute.Encode()
}

// CreateMasterReadNext 主站请求帧 请求读后续数据
func CreateMasterReadNext(address string, ident string, seq byte) ([]byte, error) {
	dlt645statute := Dlt645statute{Address: address, Control: MasterReadNext, Ident: ident, Data: &ByteArrayAnalyzer{Value: []byte{seq}}}
	return dlt645statute.Encode()
}

// MasterReadNextNormalResponseFloat 主站请求读后续数据的从站正常应答帧 浮点数
func MasterReadNextNormalResponseFloat(address string, framing bool, ident string, seq byte, value ...float64) ([]byte, error) {
	floatAnalyzer := Translate(ident)
	if floatAnalyzer == nil {
		return nil, errors.New(ident + "never enroll")
	}
	if analyzer, ok := floatAnalyzer.(*FloatAnalyser); ok {
		analyzer.Value = value
		analyzer.Mark = []byte{seq}
	} else {
		return nil, errors.New(ident + "value type error")
	}
	defer Put(ident, floatAnalyzer)
	control := MasterReadNextNormalResponseFloatNoNext
	if framing {
		control = MasterReadNextNormalResponseFloatNext
	}
	dlt645statute := Dlt645statute{Address: address, Control: control, Ident: ident, Data: floatAnalyzer}
	return dlt645statute.Encode()
}

// CreateMasterReadNextError 主站请求读后续数据的从站异常应答帧
// address 电能表地址
// errValue 异常编码
func CreateMasterReadNextError(address string, errValue byte) ([]byte, error) {
	dlt645statute := Dlt645statute{Address: address, Control: MasterReadNextError, Data: &ByteArrayAnalyzer{Value: []byte{errValue}}}
	return dlt645statute.Encode()
}

// MasterRequestMeterSetFloatData 主站向从站请求设置数据 设置浮点数
// address 电能表地址
// ident 数据标识
// auth 写块
func MasterRequestMeterSetFloatData(address string, ident string, auth *WriteAnalyzer) ([]byte, error) {
	dlt645statute := Dlt645statute{Address: address, Control: MasterRequestMeterSetData, Ident: ident, Data: auth}
	return dlt645statute.Encode()
}

// MasterRequestMeterSetFloatDataNormalResponse 主站向从站请求设置数据,从站正常返回
// address 电能表地址
func MasterRequestMeterSetFloatDataNormalResponse(address string) ([]byte, error) {
	dlt645statute := Dlt645statute{Address: address, Control: MasterRequestMeterSetDataNormalResponse}
	return dlt645statute.Encode()
}

// CreateMasterRequestMeterSetFloatDataErrorResponse 主站向从站请求设置数据,从站异常返回
// address 电能表地址
func CreateMasterRequestMeterSetFloatDataErrorResponse(address string, errValue byte) ([]byte, error) {
	dlt645statute := Dlt645statute{Address: address, Control: MasterRequestMeterSetFloatDataErrorResponse, Data: &ByteArrayAnalyzer{Value: []byte{errValue}}}
	return dlt645statute.Encode()
}

// CreateMasterReadMeterAddress 请求读电能表通信地址，仅支持点对点通信
func CreateMasterReadMeterAddress() ([]byte, error) {
	dlt645statute := Dlt645statute{Address: "AAAAAAAAAAAA", Control: MasterReadMeterAddress}
	return dlt645statute.Encode()
}
