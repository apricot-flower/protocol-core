package dlt645

import "errors"

var _ Analyser = (*WriteAnalyzer)(nil)

type WriteAnalyzer struct {
	Pa           byte
	Password     []byte
	OperatorCode []byte
	Data         []byte
}

func (w *WriteAnalyzer) Decode(buf []byte) error {
	if len(buf) < 8 {
		return errors.New("0x14 write data func: pa + password + operator length must > 8")
	}
	w.Pa = buf[0]
	w.Password = buf[1:4]
	w.OperatorCode = buf[4:8]
	w.overturn(w.OperatorCode)
	w.Data = buf[8:]
	return nil
}

func (w *WriteAnalyzer) Encode() ([]byte, error) {
	if w.OperatorCode == nil || len(w.OperatorCode) != 4 {
		return nil, errors.New("0x14 write data func: operator code is nil or operator length != 4")
	}
	if w.Password == nil || len(w.Password) != 3 {
		return nil, errors.New("0x14 write data func: password is nil or password length != 3")
	}
	encodeArray := []byte{w.Pa}
	encodeArray = append(encodeArray, w.Password...)
	encodeArray = append(encodeArray, w.OperatorCode...)
	return append(encodeArray, w.Data...), nil
}

func (w *WriteAnalyzer) GetValue() interface{} {
	return w
}

// JoinByteParams 添加byte或者[]byte  这个是追加，不是覆盖
// overturn 是否翻转data
// data 设置的数据，这个加载后并不会负载原来的数据，而是追加到原来数据的后面
func (w *WriteAnalyzer) JoinByteParams(overturn bool, data ...byte) {
	if overturn {
		w.overturn(data)
	}
	w.Data = append(w.Data, data...)
}

// JoinFloatParams 添加浮点数  这个是追加，不是覆盖
// length 浮点数占用的字节数量
// rate 倍率
// data 设置的数据(原始值)，这个加载后并不会负载原来的数据，而是追加到原来数据的后面
func (w *WriteAnalyzer) JoinFloatParams(length byte, rate float64, data ...float64) error {
	floating := &FloatAnalyser{Length: length, Rate: rate, Value: data}
	arr, err := floating.Encode()
	if err != nil {
		return err
	}
	w.Data = append(w.Data, arr...)
	return nil
}

// JoinTwoBitArrayParams 添加二维数组  这个是追加，不是覆盖
// overturn 是否翻转每一个数组
// data 设置的数据，这个加载后并不会负载原来的数据，而是追加到原来数据的后面
func (w *WriteAnalyzer) JoinTwoBitArrayParams(overturn bool, data ...[]byte) {
	for _, d := range data {
		w.JoinByteParams(overturn, d...)
	}
}

func (w *WriteAnalyzer) overturn(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func (w *WriteAnalyzer) Clean() {
	w.Pa = 0
	w.Password = nil
	w.OperatorCode = nil
	w.Data = nil
}
