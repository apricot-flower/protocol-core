package dlt645

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
	"sync"
)

type Analyser interface {
	Decode(buf []byte) error
	Encode() ([]byte, error)
	GetValue() interface{}
	Clean()
}

var Dlt698Codec = make(map[string]*sync.Pool)

func Translate(ident string) Analyser {
	if pool, ok := Dlt698Codec[strings.ToUpper(ident)]; ok {
		if analyzer, ok := pool.Get().(Analyser); ok {
			return analyzer
		}
	}
	return nil
}

func Put(ident string, analyser Analyser) {
	if pool, ok := Dlt698Codec[strings.ToUpper(ident)]; ok {
		analyser.Clean()
		pool.Put(analyser)
	}
}

const (
	MasterReadMeterData                         byte = 0x11
	MasterReadMeterNormalResponseNoNext         byte = 0x91
	MasterReadMeterNormalResponseNext           byte = 0xB1
	MasterReadMeterDataErrorResponse            byte = 0xD1
	MasterReadNext                              byte = 0x12
	MasterReadNextNormalResponseFloatNext       byte = 0xB2
	MasterReadNextNormalResponseFloatNoNext     byte = 0x92
	MasterReadNextError                         byte = 0xD2
	MasterRequestMeterSetData                   byte = 0x14
	MasterRequestMeterSetDataNormalResponse     byte = 0x94
	MasterRequestMeterSetFloatDataErrorResponse byte = 0xD4
	MasterReadMeterAddress                      byte = 0x13
	MasterReadMeterAddressNormalResponse        byte = 0x93
	MasterWriteMeterAddress                     byte = 0x15
	MasterWriteMeterAddressNormalResponse       byte = 0x95
	BroadcastingProofreadTime                   byte = 0x08
)

// IdentHandle 解析标识
func IdentHandle(buf *bytes.Reader) (string, error) {
	identArray := make([]byte, 4)
	err := binary.Read(buf, binary.LittleEndian, &identArray)
	if err != nil {
		return "", err
	}
	for i, j := 0, len(identArray)-1; i < j; i, j = i+1, j-1 {
		identArray[i], identArray[j] = identArray[j], identArray[i]
	}
	return strings.ToUpper(hex.EncodeToString(identArray)), nil
}

var FrameDecodeErrorHandle FrameDecodeError

type CallBackHandle interface{}

var controlHandles = make(map[byte]CallBackHandle)

// MasterReadMeterDataHandle 主站请求帧 请求读电能表数据 控制码：C=11H
type MasterReadMeterDataHandle func(address string, control byte, ident string, seq byte, param ...byte)
type EnrollNormalResponseFloatHandle func(address string, control byte, ident string, seq byte, value ...float64)

func EnrollControlHandle(controlCode byte, handle CallBackHandle) error {
	switch controlCode {
	case MasterReadMeterData, MasterReadNext, MasterRequestMeterSetData:
		if _, ok := handle.(MasterReadMeterDataHandle); !ok {
			return errors.New("handle is not a MasterReadMeterDataHandle")
		}
		controlHandles[controlCode] = handle
	case MasterReadMeterNormalResponseNoNext, MasterReadMeterNormalResponseNext, MasterReadNextNormalResponseFloatNext, MasterReadNextNormalResponseFloatNoNext:
		if _, ok := handle.(EnrollNormalResponseFloatHandle); !ok {
			return errors.New("control handle is not a MasterReadMeterDataHandle")
		}
		controlHandles[controlCode] = handle

	default:
		return errors.New("this control code is not supported")
	}
	return nil
}

func EnrollNormalResponseFloat(length byte, rate float64, handle CallBackHandle, ident ...string) error {
	floatPool := &sync.Pool{
		New: func() interface{} {
			return &FloatAnalyser{Length: length, Rate: rate}
		}}
	for _, id := range ident {
		Dlt698Codec[strings.ToUpper(id)] = floatPool
	}
	err := EnrollControlHandle(MasterReadMeterNormalResponseNoNext, handle)
	if err != nil {
		return err
	}
	err = EnrollControlHandle(MasterReadNextNormalResponseFloatNext, handle)
	if err != nil {
		return err
	}
	err = EnrollControlHandle(MasterReadNextNormalResponseFloatNoNext, handle)
	if err != nil {
		return err
	}
	err = EnrollControlHandle(MasterReadMeterNormalResponseNext, handle)
	return err
}

func ControlOverTurn(frame string, address string, control byte, ident string, data Analyser) {
	switch control {
	case MasterReadMeterData:
		if handle, ok := controlHandles[control].(MasterReadMeterDataHandle); ok {
			if data == nil {
				handle(address, control, ident, 0)
			}
			if analyzer, ok := data.(*ByteArrayAnalyzer); ok {
				handle(address, control, ident, 0, analyzer.Value...)
			}
		}
	case MasterReadMeterNormalResponseNoNext, MasterReadMeterNormalResponseNext:
		if handle, ok := controlHandles[control].(EnrollNormalResponseFloatHandle); ok {
			if analyzer, ok := data.(*FloatAnalyser); ok {
				handle(address, control, ident, 0, analyzer.Value...)
			}
		}
	case MasterReadMeterDataErrorResponse, MasterReadNextError:
		if FrameDecodeErrorHandle != nil {
			FrameDecodeErrorHandle(frame, address, control, nil)
		}
	case MasterReadNext:
		if handle, ok := controlHandles[control].(MasterReadMeterDataHandle); ok {
			if data == nil {
				handle(address, control, ident, 0)
			}
			if analyzer, ok := data.(*ByteArrayAnalyzer); ok {
				handle(address, control, ident, analyzer.Value[0])
			}
		}
	case MasterReadNextNormalResponseFloatNext, MasterReadNextNormalResponseFloatNoNext:
		if handle, ok := controlHandles[control].(EnrollNormalResponseFloatHandle); ok {
			if analyzer, ok := data.(*FloatAnalyser); ok {
				handle(address, control, ident, analyzer.Mark[0], analyzer.Value...)
			}
		}
	case MasterRequestMeterSetData:
		if handle, ok := controlHandles[control].(MasterReadMeterDataHandle); ok {
			if analyzer, ok := data.(*WriteAnalyzer); ok {
				handle(address, control, ident, 0, analyzer.Data...)
			}
		}
	default:
		if FrameDecodeErrorHandle != nil {
			FrameDecodeErrorHandle(frame, address, control, errors.New("控制域未注册或错误！"))
		}
	}
}

// FrameDecodeError 统一的异常解析处理方法
type FrameDecodeError func(frame string, address string, errCode byte, err error)

func EnrollFrameDecodeErrorHandle(handle FrameDecodeError) {
	FrameDecodeErrorHandle = handle
}
