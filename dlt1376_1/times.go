package dlt1376_1

import (
	"fmt"
	"strconv"
	"time"
)

type Dlt13761DataInter interface {
	Decode(buf ...byte) error
	Encode() ([]byte, error)
}

var _ Dlt13761DataInter = (*TimeA1)(nil)
var _ Dlt13761DataInter = (*TimeA15)(nil)

type TimeA1 struct {
	Second    byte
	Minute    byte
	Hour      byte
	Day       byte
	Week      byte
	Month     byte
	weekMonth byte
	Year      byte
}

func (t *TimeA1) Decode(buf ...byte) error {
	second := buf[0]
	minute := buf[1]
	hour := buf[2]
	day := buf[3]
	weekMonth := buf[4]
	year := buf[5]
	//解析秒
	t.Second = extractDigits(second)
	t.Minute = extractDigits(minute)
	t.Hour = extractDigits(hour)
	t.Day = extractDigits(day)
	//解析星期-月
	bits5To7, bit4, bits0To3 := t.extractBits(weekMonth)
	t.Week = bits5To7
	t.Month = bit4*10 + bits0To3
	//解析年
	t.Year = extractDigits(year)
	return nil
}

func (t *TimeA1) extractBits(b byte) (bits5To7 byte, bit4 byte, bits0To3 byte) {
	// 提取 bit5 到 bit7
	bits5To7 = (b & 0xE0) >> 5
	// 提取 bit4
	bit4 = (b & 0x10) >> 4
	// 提取 bit0 到 bit3
	bits0To3 = b & 0x0F
	return
}

func extractDigits(b byte) byte {
	tens := (b >> 4) & 0xF
	units := b & 0xF
	result := tens*10 + units
	return result
}

func (t *TimeA1) Encode() ([]byte, error) {
	return []byte{t.Second, t.Minute, t.Hour, t.Day, t.weekMonth, t.Year}, nil
}

func (t *TimeA1) Build() error {
	var err error
	// 获取当前时间
	currentTime := time.Now()
	// 获取年份
	year := currentTime.Year() % 100
	t.Year, err = montage(year)
	if err != nil {
		return err
	}
	month := currentTime.Month()
	weekDay := currentTime.Weekday()
	t.weekMonth = t.combineBits(byte(weekDay), byte(month/10), byte(month%10))
	// 获取日
	day := currentTime.Day()
	t.Day, err = montage(day)
	if err != nil {
		return err
	}
	// 获取小时
	hour := currentTime.Hour()
	t.Hour, err = montage(hour)
	if err != nil {
		return err
	}
	// 获取分钟
	minute := currentTime.Minute()
	t.Minute, err = montage(minute)
	if err != nil {
		return err
	}
	// 获取秒
	second := currentTime.Second()
	t.Second, err = montage(second)
	return err
}

// 新的 combineBits 函数，用于将三个部分重新组合成一个字节
func (t *TimeA1) combineBits(bits5To7 byte, bit4 byte, bits0To3 byte) byte {
	// 将 bits5To7 左移 5 位
	part1 := bits5To7 << 5
	// 将 bit4 左移 4 位
	part2 := bit4 << 4
	// bits0To3 保持不变
	part3 := bits0To3

	// 合并所有部分
	result := part1 | part2 | part3
	return result
}

func montage(data int) (byte, error) {
	value := byte(data)
	tens := value / 10  // 十位
	units := value % 10 // 个位
	binStr := fmt.Sprintf("%04b", tens) + fmt.Sprintf("%04b", units)
	result, err := strconv.ParseUint(binStr, 2, 8)
	if err != nil {
		return 0, err
	}
	return byte(result), nil
}

type TimeA15 struct {
	Minute byte
	Hour   byte
	Day    byte
	Month  byte
	Year   byte
}

func (t *TimeA15) Decode(buf ...byte) error {
	minute := buf[0]
	hour := buf[1]
	day := buf[2]
	month := buf[3]
	year := buf[4]
	t.Minute = extractDigits(minute)
	t.Hour = extractDigits(hour)
	t.Day = extractDigits(day)
	t.Month = extractDigits(month)
	t.Year = extractDigits(year)
	return nil

}

func (t *TimeA15) Encode() ([]byte, error) {
	minute, err := montage(int(t.Minute))
	if err != nil {
		return nil, err
	}
	hour, err := montage(int(t.Hour))
	if err != nil {
		return nil, err
	}
	day, err := montage(int(t.Day))
	if err != nil {
		return nil, err
	}
	month, err := montage(int(t.Month))
	if err != nil {
		return nil, err
	}
	year, err := montage(int(t.Year))
	if err != nil {
		return nil, err
	}
	return []byte{minute, hour, day, month, year}, nil
}
