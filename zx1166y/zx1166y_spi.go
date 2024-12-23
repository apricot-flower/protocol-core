package zx1166y

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/ecc1/spi"
	"strings"
)

type SPICodec struct {
	device *spi.Device
	Dev    string
	Mode   uint8
	Speed  int
}

func (s *SPICodec) Open() error {
	device, err := spi.Open(s.Dev, s.Speed, 0)
	if err != nil {
		return err
	}
	s.device = device
	err = s.device.SetMaxSpeed(s.Speed)
	if err != nil {
		return err
	}
	err = s.device.SetMode(s.Mode)
	return err
}

func (s *SPICodec) Close() error {
	if s.device == nil {
		return nil
	}
	return s.device.Close()
}

// TerminalActiveReport2 终端主动上报 第2个步骤
func (s *SPICodec) TerminalActiveReport2(secureFlag string, attachData []byte, data []byte, mac []byte) ([]byte, error) {
	tx, err := s.encode(secureFlag, attachData, data, mac)
	if err != nil {
		return nil, err
	}
	return s.TransferBytes(tx, 2048)
}

// TerminalActiveReport1 终端主动上报 第1个步骤
func (s *SPICodec) TerminalActiveReport1(data []byte) (resultData []byte, mac []byte, err error) {
	tx := "80140103"
	txArr, err := s.encode(tx, data)
	if err != nil {
		return nil, nil, err
	}
	result, err := s.TransferBytes(txArr, 4096)
	if err != nil {
		return nil, nil, err
	}
	length := len(result)
	if length < 4 {
		return nil, nil, errors.New("返回结果容量错误！")
	}
	return data[:length-4], data[length-4:], nil
}

// UpdateSessionTimeLimit 更新会话时效门限
func (s *SPICodec) UpdateSessionTimeLimit(data []byte) error {
	tx := "81340105"
	txArr, err := s.encode(tx, data)
	if err != nil {
		return err
	}
	_, err = s.TransferBytes(txArr, 1024)
	return err
}

// VerifyBySid 5.3.4 若数据验证信息为SID
func (s *SPICodec) VerifyBySid(secureFlag string, attachData, data []byte) ([]byte, error) {
	return s.VerifyBySidMac(secureFlag, attachData, data, nil)
}

// VerifyBySidMac 5.3.4 若数据验证信息为SID_MAC
// secureFlag 安全标识
// attachData 附加数据
// data
// mac
func (s *SPICodec) VerifyBySidMac(secureFlag string, attachData, data, mac []byte) ([]byte, error) {
	tx, err := s.encode(secureFlag, attachData, data, mac)
	if err != nil {
		return nil, err
	}
	rx, err := s.TransferBytes(tx, 4096)
	if err != nil {
		return nil, err
	}
	return rx, nil
}

// VerifySelectSecureFlag 5.3.3 获取安全标识
// lastValidBit 组地址或广播地址最后一位有效位
func (s *SPICodec) VerifySelectSecureFlag(lastValidBit byte) string {
	secureFlag := `8016480`
	switch lastValidBit {
	case 1:
		secureFlag = secureFlag + "01"
	case 2:
		secureFlag = secureFlag + "02"
	case 3:
		secureFlag = secureFlag + "03"
	case 4:
		secureFlag = secureFlag + "04"
	case 5:
		secureFlag = secureFlag + "05"
	case 6:
		secureFlag = secureFlag + "06"
	case 7:
		secureFlag = secureFlag + "07"
	case 8:
		secureFlag = secureFlag + "08"
	case 9:
		secureFlag = secureFlag + "09"
	default:
		secureFlag = secureFlag + "0A"
	}
	return secureFlag
}

// CertificateUpdate 证书更新
// enData 证书内容
// sid 安全标识
// attachData 附加数据
func (s *SPICodec) CertificateUpdate(enData, sid, attachData []byte) error {
	if strings.ToUpper(hex.EncodeToString(sid)) != "81300203" {
		return errors.New("安全标识错误，应为：[81300203] ！")
	}
	tx, err := s.encode("81300203", attachData, enData)
	if err != nil {
		return err
	}
	_, err = s.TransferBytes(tx, 1024)
	return err
}

// TerminalSymmetricKeyUpdate 终端对称密钥更新
// secureFlag 安全标识
// attachData 附加数据
// mac
// enData 密文
func (s *SPICodec) TerminalSymmetricKeyUpdate(secureFlag, attachData, mac, enData []byte) ([]byte, error) {
	//验证安全标识
	if strings.ToUpper(hex.EncodeToString(secureFlag)) != "812E0000" {
		return nil, errors.New("安全标识错误，应为：[812E0000] ！")
	}
	tx, err := s.encode("812E0000", attachData, enData, mac)
	if err != nil {
		return nil, err
	}
	data, err := s.TransferBytes(tx, 4096)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Encrypt 加密  5.3.7
// encrypt 解密类型
// data
func (s *SPICodec) Encrypt(encrypt EncryptType, data []byte) ([]byte, error) {
	tx := "801C00" + string(encrypt)
	txArr, err := s.encode(tx, data)
	if err != nil {
		return nil, err
	}
	return s.TransferBytes(txArr, 4096)
}

// ReadTerminal 抄读终端
// rn 抄读随机数，由主站下发
// data要抄读的数据
func (s *SPICodec) ReadTerminal(rn, data []byte) ([]byte, error) {
	tx, err := s.encode("800E4002", rn, data)
	if err != nil {
		return nil, err
	}
	return s.TransferBytes(tx, 1024)
}

// SessionKeyConnect 建立应用连接（会话密钥协商）
// ucSessionData：服务器随机数，48 字节
// ucSign:服务器签名信息，Len-48 字节
func (s *SPICodec) SessionKeyConnect(ucOutSessionInit, ucOutSign []byte) (ucSessionData string, ucSign string, err error) {
	if len(ucOutSessionInit) == 32 {
		return "", "", errors.New("ucOutSessionInit长度应为32！")
	}
	tx, err := s.encode("80020000", ucOutSessionInit, ucOutSign)
	if err != nil {
		return "", "", err
	}
	data, err := s.TransferBytes(tx, 1024)
	if err != nil {
		return "", "", err
	}
	ucSessionData = hex.EncodeToString(data[:48])
	ucSign = hex.EncodeToString(data[48:])
	return ucSessionData, ucSign, nil
}

// SelectMasterStationCertificate 获取主站证书
func (s *SPICodec) SelectMasterStationCertificate() (string, error) {
	tx := `558036000C000045`
	data, err := s.TransferString(tx, 4096)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

// SelectTerminalCertificate 获取终端证书
func (s *SPICodec) SelectTerminalCertificate() (string, error) {
	tx := `558036000B000042`
	data, err := s.TransferString(tx, 4096)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

// SelectESAMInfos 主站下发获取ESAM信息命令
func (s *SPICodec) SelectESAMInfos() (*TESABInfo, error) {
	tx := `55803600FF0000B6`
	data, err := s.TransferString(tx, 1024)
	if err != nil {
		return nil, err
	}
	info := &TESABInfo{}
	err = info.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *SPICodec) TransferString(tx string, length int) ([]byte, error) {
	txArr, err := hex.DecodeString(tx)
	if err != nil {
		return nil, err
	}
	return s.TransferBytes(txArr, length)
}

func (s *SPICodec) TransferBytes(tx []byte, length int) ([]byte, error) {
	txArr := make([]byte, length)
	copy(txArr, tx)
	rxArr := make([]byte, length)
	err := s.device.Transfer(txArr, rxArr)
	if err != nil {
		return nil, err
	}
	rxArr, err = s.decode(rxArr)
	if err != nil {
		return nil, err
	}
	return rxArr, nil
}

// 拆包器
func (s *SPICodec) decode(array []byte) ([]byte, error) {
	buf := bytes.NewReader(array)
	var err error
	var start byte
	for start != 0x55 {
		err = binary.Read(buf, binary.BigEndian, &start)
		if err != nil {
			return nil, err
		}
	}
	//读取状态码
	status := make([]byte, 2)
	err = binary.Read(buf, binary.BigEndian, &status)
	if err != nil {
		return nil, err
	}
	//验证错误码
	if hex.EncodeToString(status) != "9000" {
		return nil, errors.New("不是正确报文，返回了错误码，错误码为：" + hex.EncodeToString(status))
	}
	//解析长度
	length := make([]byte, 2)
	err = binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}
	if uint16(length[0])<<8|uint16(length[1]) == 0 {
		//长度为0，没有数据，证书更新的时候会出现这种情况
		return nil, nil
	}
	data := make([]byte, uint16(length[0])<<8|uint16(length[1]))
	err = binary.Read(buf, binary.BigEndian, &data)
	if err != nil {
		return nil, err
	}
	var csValue byte
	err = binary.Read(buf, binary.BigEndian, &csValue)
	if err != nil {
		return nil, err
	}
	//todo 计算cs
	csData := append(status, length...)
	csData = append(csData, data...)
	if s.Cs(csData) != csValue {
		return nil, errors.New("cs错误！")
	}
	return data, nil
}

// 封装报文
func (s *SPICodec) encode(tx string, data ...[]byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	txArr, err := hex.DecodeString(tx)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, txArr)
	if err != nil {
		return nil, err
	}
	var length uint16
	var dataArr []byte
	for _, v := range data {
		length += uint16(len(v))
		dataArr = append(dataArr, v...)
	}
	err = binary.Write(buf, binary.BigEndian, length)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, dataArr)
	if err != nil {
		return nil, err
	}
	txArr = buf.Bytes()
	buf.Reset()
	cs := s.Cs(txArr)
	txArr = append([]byte{0x55}, txArr...)
	txArr = append(txArr, cs)
	return txArr, nil
}

// Cs 计算校验值
func (s *SPICodec) Cs(data []byte) byte {
	var xorSum byte = 0
	for _, b := range data {
		xorSum ^= b
	}

	// 取反最终的异或值
	invertedXorSum := ^xorSum
	return invertedXorSum
}

// 终端抄读电表

// ReadMeter1 终端抄读电表第1个步骤
func (s *SPICodec) ReadMeter1() ([]byte, error) {
	tx, err := hex.DecodeString("800400100000")
	if err != nil {
		return nil, err
	}
	csValue := s.Cs(tx)
	tx = append([]byte{0x55}, tx...)
	tx = append(tx, csValue)
	return s.TransferBytes(tx, 4096)
}

// ReadMeter8 终端抄读电表第8个步骤
func (s *SPICodec) ReadMeter8(meterId, rand, data, mac []byte) error {
	tx, err := s.encode("800E4887", meterId, rand, data, mac)
	if err != nil {
		return err
	}
	_, err = s.TransferBytes(tx, 4096)
	return err
}

// ReadMeter9 终端抄读电表第9个步骤
func (s *SPICodec) ReadMeter9(meterId, rn, data []byte) ([]byte, error) {
	tx, err := s.encode("800C4807", meterId, rn, data)
	if err != nil {
		return nil, err
	}
	return s.TransferBytes(tx, 4096)
}

// ReadMeter10 终端抄读电表第10个步骤
func (s *SPICodec) ReadMeter10(meterId, rn, data, mac []byte) ([]byte, error) {
	tx, err := s.encode("80124807", meterId, rn, data, mac)
	if err != nil {
		return nil, err
	}
	return s.TransferBytes(tx, 4096)
}

type TESABInfo struct {
	ESAMNumber                    string `json:"esam_number"`                      //ESAM序列号
	ESAMVersion                   string `json:"esam_version"`                     //ESAM版本号
	SymmetricKeyVersion           string `json:"symmetric_key_version"`            //对称密钥版本
	MainStationCertificateVersion byte   `json:"main_station_certificate_version"` //主站证书版本号
	TerminalCertificateVersion    byte   `json:"terminal_certificate_version"`     //终端证书版本号
	SessionTimeLimit              uint32 `json:"session_time_limit"`               //会话时效门限
	SessionTimeRemainingTime      uint32 `json:"session_time_remaining_time"`      //会话时效剩余时间
	ASCTR                         uint32 `json:"ASCTR"`                            //单地址应用协商计数器
	ARCTR                         uint32 `json:"ARCTR"`                            //主动上报计数器
	AGSEQ                         uint32 `json:"AGSEQ"`                            //应用广播通信序列号
	TerminalCertificateNumber     string `json:"terminal_certificate_number"`      //终端证书序列号
	MainStationCertificateNumber  string `json:"main_station_certificate_number"`  //主站证书序列号
}

func (t *TESABInfo) Decode(buf *bytes.Reader) error {
	//ESAM序列号
	ESAMNumber := make([]byte, 8)
	err := binary.Read(buf, binary.BigEndian, &ESAMNumber)
	if err != nil {
		return err
	}
	t.ESAMNumber = hex.EncodeToString(ESAMNumber)
	//ESAM版本号
	ESAMVersion := make([]byte, 4)
	err = binary.Read(buf, binary.BigEndian, &ESAMVersion)
	if err != nil {
		return err
	}
	t.ESAMVersion = hex.EncodeToString(ESAMVersion)
	//对称密钥版本
	SymmetricKeyVersion := make([]byte, 16)
	err = binary.Read(buf, binary.BigEndian, &SymmetricKeyVersion)
	if err != nil {
		return err
	}
	t.SymmetricKeyVersion = hex.EncodeToString(SymmetricKeyVersion)
	//主站证书版本号
	err = binary.Read(buf, binary.BigEndian, &t.MainStationCertificateVersion)
	if err != nil {
		return err
	}
	//终端证书版本号
	err = binary.Read(buf, binary.BigEndian, &t.TerminalCertificateVersion)
	if err != nil {
		return err
	}
	//会话时效门限
	err = binary.Read(buf, binary.BigEndian, &t.SessionTimeLimit)
	if err != nil {
		return err
	}
	//会话时效剩余时间
	err = binary.Read(buf, binary.BigEndian, &t.SessionTimeRemainingTime)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &t.ASCTR)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &t.ARCTR)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &t.AGSEQ)
	if err != nil {
		return err
	}
	TerminalCertificateNumber := make([]byte, 16)
	err = binary.Read(buf, binary.BigEndian, &TerminalCertificateNumber)
	if err != nil {
		return err
	}
	t.TerminalCertificateNumber = hex.EncodeToString(TerminalCertificateNumber)
	MainStationCertificateNumber := make([]byte, 16)
	err = binary.Read(buf, binary.BigEndian, &MainStationCertificateNumber)
	if err != nil {
		return err
	}
	t.MainStationCertificateNumber = hex.EncodeToString(MainStationCertificateNumber)
	return nil
}

// EncryptType 加密类型
type EncryptType string

const (
	Plaintext_MAC     EncryptType = "11" //明文+MAC
	CiphertextEncrypt EncryptType = "96" //密文
	Ciphertext        EncryptType = "97" //密文+MAC
)
