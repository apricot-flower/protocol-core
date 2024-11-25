package dlt1376_1

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
)

const (
	COMFIRM_DENY             byte = 0x00
	RESET_COMMAND            byte = 0x01
	LINK_INTERFACE_DETECTION byte = 0x02 //链路接口检测
	RELAY_STATION_COMMAND    byte = 0x03
)

var Afns map[byte]PostAFN

func init() {
	Afns = make(map[byte]PostAFN)
}

func Translate(flag byte) Afn {
	if afnFunc, ok := Afns[flag]; ok {
		return afnFunc()
	}
	return nil
}

type PostAFN func() Afn

type Dlt13761Data struct {
	P    uint64
	F    uint64
	Data Dlt13761DataInter
}

type Afn interface {
	Decode(buf *bytes.Reader) error   //解码
	Encode() ([]byte, error)          //编码
	Idents() ([]*Dlt13761Data, error) //获取数据
	Flag() (byte, string)             //获取AFN和他的说明
	HasAux() bool                     //是否存在附加信息域
	Append(*Dlt13761Data) error       //添加一个数据项
}

func analyzeUnit(daDt []byte) ([]uint64, []uint64) {
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
	return pn, fn
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

func checkData(data ...byte) bool {
	for _, d := range data {
		if d == 0xEE {
			return false
		}
	}
	return true
}

func sortByF(dataSlice []*Dlt13761Data) {
	sort.Slice(dataSlice, func(i, j int) bool {
		return dataSlice[i].F < dataSlice[j].F
	})
}
