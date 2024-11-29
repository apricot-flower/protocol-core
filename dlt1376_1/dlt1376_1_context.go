package dlt1376_1

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
	"strconv"
)

type plugIn func() LinkDataInter

var afnMap = make(map[byte]plugIn)

func init() {
	//确认∕否认
	afnMap[ConfirmDenyIdent] = func() LinkDataInter {
		return &ConfirmDeny{}
	}
	//复位
	afnMap[ResetIdent] = func() LinkDataInter {
		return &Reset{}
	}
	//链路接口检测
	afnMap[LinkInterfaceDetectionIdent] = func() LinkDataInter {
		return &LinkInterfaceDetection{}
	}
	//中继站命令
	afnMap[RelayStationCommandIdent] = func() LinkDataInter {
		return &RelayStationCommand{}
	}
	//设置参数
	afnMap[0x04] = func() LinkDataInter {
		return nil
	}
	//控制命令
	afnMap[0x05] = func() LinkDataInter {
		return nil
	}
	//身份认证及密钥协商
	afnMap[0x06] = func() LinkDataInter {
		return nil
	}
	//备用
	afnMap[0x07] = func() LinkDataInter {
		return nil
	}
	//请求被级联终端主动上报
	afnMap[0x08] = func() LinkDataInter {
		return nil
	}
	//请求终端配置
	afnMap[0x09] = func() LinkDataInter {
		return nil
	}
	//查询参数
	afnMap[0x0A] = func() LinkDataInter {
		return nil
	}
	//请求任务数据
	afnMap[0x0B] = func() LinkDataInter {
		return nil
	}
	//请求1类数据（实时数据）
	afnMap[0x0C] = func() LinkDataInter {
		return nil
	}
	//请求2类数据（历史数据）
	afnMap[0x0D] = func() LinkDataInter {
		return nil
	}
	//请求3类数据（事件数据）
	afnMap[0x0E] = func() LinkDataInter {
		return nil
	}
	//文件传输
	afnMap[0x0F] = func() LinkDataInter {
		return nil
	}
	//数据转发
	afnMap[0x10] = func() LinkDataInter {
		return nil
	}
}

func translate(afn byte) LinkDataInter {
	if inter, ok := afnMap[afn]; ok {
		return inter()
	}
	return nil
}

func analyzeUnit(buf *bytes.Reader) ([]uint64, []uint64, error) {
	var daDt [4]byte
	err := binary.Read(buf, binary.LittleEndian, &daDt)
	if err != nil {
		return nil, nil, errors.New("analyze data_unit err: " + err.Error())
	}
	var pn []uint64
	da1 := daDt[0]
	da2 := daDt[1]
	//解析pn
	if da1 == 0 && da2 == 0 {
		pn = append(pn, 0)
	} else if da1 == 0xff && da2 == 0 {
		for index := 1; index < 2040; index++ {
			pn = append(pn, uint64(index))
		}
	} else {
		da1Str := fmt.Sprintf("%08b", da1)
		for da1Index, da1Value := range da1Str {
			if da1Value == '1' {
				p := uint64(da2+1)*8 - (7 - (7 - uint64(da1Index)))
				pn = append(pn, p)
			}
		}
	}
	var fn []uint64
	dt1 := daDt[2]
	dt2 := daDt[3]
	dt1Str := fmt.Sprintf("%08b", dt1)
	for dt1Index, dt1Value := range dt1Str {
		if dt1Value == '1' {
			f := uint64(dt2+1)*8 - (7 - (7 - uint64(dt1Index)))
			if f <= 248 {
				fn = append(fn, f)
			}
		}
	}
	sort.Slice(pn, func(i, j int) bool {
		return pn[i] < pn[j]
	})
	sort.Slice(fn, func(i, j int) bool {
		return fn[i] < fn[j]
	})
	return pn, fn, nil
}

func encodeUnit(pn uint64, fn uint64) ([]byte, error) {
	//1. 先解析pn
	pns, err := encodePnFn(pn)
	if err != nil {
		return nil, err
	}
	//在解析fn
	if fn > 248 {
		return nil, errors.New("fn must > 248")
	}
	fns, err := encodePnFn(fn)
	if err != nil {
		return nil, err
	}
	return append(pns, fns...), nil
}

func encodePnFn(pn uint64) ([]byte, error) {
	if pn == 0 {
		return []byte{0x00, 0x00}, nil
	}
	//正经解析
	if pn > 2040 {
		return nil, errors.New("pn must <= 2040")
	}
	da2 := byte(pn / 8)
	//解析da1
	da1Index := int(pn % 8)
	var bitIndex int
	if da1Index == 0 {
		bitIndex = 7
	} else if da1Index == 1 {
		bitIndex = 0
	} else {
		bitIndex = da1Index - 1
	}
	da1, err := setBit(bitIndex)
	if err != nil {
		return nil, err
	}
	return []byte{da1, da2}, nil
}

func setBit(n int) (byte, error) {
	s := "00000000"
	runes := []rune(s)
	// 设置第 n 位为 '1'
	runes[len(s)-1-n] = '1'
	// 将 rune 切片转换回字符串
	s = string(runes)
	value, err := strconv.ParseInt(s, 2, 8)
	if err != nil {
		return 0, err
	}
	return byte(value), nil
}
