package dlt645

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
)

var _ Analyser = (*FloatAnalyser)(nil)

type FloatAnalyser struct {
	Length byte
	Rate   float64
	Value  []float64
	Mark   []byte
}

func (f *FloatAnalyser) Clean() {
	f.Value = nil
	f.Mark = nil
}

func (f *FloatAnalyser) Decode(buf []byte) error {
	dataArray := f.splitBytes(buf)
	f.Value = make([]float64, len(dataArray))
	for index, data := range dataArray {
		dataStr := hex.EncodeToString(data)
		floating, err := strconv.ParseFloat(dataStr, 64)
		if err != nil {
			return err
		}
		f.Value[index] = floating * f.Rate
	}
	return nil
}

func (f *FloatAnalyser) Encode() ([]byte, error) {
	var err error
	buf := new(bytes.Buffer)
	for _, v := range f.Value {
		floatValue := v / f.Rate
		var dataStr string
		switch f.Length {
		case 0:
			return nil, nil
		case 1:
			dataStr = f.formatIntWithLeadingZeros(int8(floatValue))
		case 2:
			dataStr = f.formatIntWithLeadingZeros(int16(floatValue))
		case 4:
			dataStr = f.formatIntWithLeadingZeros(int32(floatValue))
		case 8:
			dataStr = f.formatIntWithLeadingZeros(int64(floatValue))
		}
		if len(dataStr)%2 != 0 {
			dataStr = "0" + dataStr
		}
		dataArray, err := hex.DecodeString(dataStr)
		if err != nil {
			return nil, err
		}
		for ii, j := 0, len(dataArray)-1; ii < j; ii, j = ii+1, j-1 {
			dataArray[ii], dataArray[j] = dataArray[j], dataArray[ii]
		}
		err = binary.Write(buf, binary.LittleEndian, dataArray)
		if err != nil {
			return nil, err
		}
	}
	if f.Mark != nil {
		err = binary.Write(buf, binary.LittleEndian, f.Mark)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), err
}

func (f *FloatAnalyser) GetValue() interface{} {
	return f.Value
}

func (f *FloatAnalyser) formatIntWithLeadingZeros(value interface{}) string {
	var width int
	// 根据值的类型确定宽度
	switch value.(type) {
	case int8:
		width = 2
	case int16:
		width = 4
	case int32:
		width = 8
	case int64:
		width = 16
	default:
		return "Unsupported type"
	}
	// 格式化整数并添加前导零
	return fmt.Sprintf("%0*d", width, value)
}

func (f *FloatAnalyser) splitBytes(data []byte) [][]byte {
	var result [][]byte
	for i := 0; i < len(data); i += int(f.Length) {
		end := i + int(f.Length)
		if end > len(data) {
			end = len(data)
		}
		child := data[i:end]
		for ii, j := 0, len(child)-1; ii < j; ii, j = ii+1, j-1 {
			child[ii], child[j] = child[j], child[ii]
		}
		result = append(result, child)
	}
	return result
}
