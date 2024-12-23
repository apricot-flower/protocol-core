package broker

import (
	"bufio"
	"encoding/binary"
	"errors"
)

type Decoder struct {
}

func (d *Decoder) Decode(reader *bufio.Reader) ([]byte, error) {
	var err error
	var frameStartChar byte
	var frameEndChar byte
	err = binary.Read(reader, binary.BigEndian, &frameStartChar)
	if err != nil {
		return nil, err
	}
	if frameStartChar != startChar {
		return nil, errors.New("invalid start char")
	}
	var frameArray []byte
	//拆eventId
	eventIdArray := make([]byte, 8)
	err = binary.Read(reader, binary.BigEndian, &eventIdArray)
	if err != nil {
		return nil, err
	}
	frameArray = append(frameArray, eventIdArray...)
	//拆idLength
	idLengthArray := make([]byte, 2)
	err = binary.Read(reader, binary.BigEndian, &idLengthArray)
	if err != nil {
		return nil, err
	}
	frameArray = append(frameArray, idLengthArray...)
	//拆id
	idArray := make([]byte, binary.BigEndian.Uint16(idLengthArray)+2)
	err = binary.Read(reader, binary.BigEndian, &idArray)
	if err != nil {
		return nil, err
	}
	frameArray = append(frameArray, idArray...)
	//拆length
	length := make([]byte, 4)
	err = binary.Read(reader, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}
	frameArray = append(frameArray, length...)
	//拆数据
	dataArray := make([]byte, binary.BigEndian.Uint32(length)+2)
	err = binary.Read(reader, binary.BigEndian, &dataArray)
	if err != nil {
		return nil, err
	}
	frameArray = append(frameArray, dataArray...)
	//拆分needResponse和cs
	endArray := make([]byte, 2)
	err = binary.Read(reader, binary.BigEndian, &endArray)
	if err != nil {
		return nil, err
	}
	frameArray = append(frameArray, endArray...)
	err = binary.Read(reader, binary.BigEndian, &frameEndChar)
	if err != nil {
		return nil, err
	}
	if frameEndChar != endChar {
		return nil, errors.New("invalid end char")
	}
	return frameArray, nil
}
