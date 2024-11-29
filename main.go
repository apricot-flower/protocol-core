package main

import (
	"encoding/hex"
	"fmt"
	"protocol-core/dlt1376_1"
)

func main() {
	frame := `6832003200680402010403140371000002009816`
	dlt13761Statute := &dlt1376_1.StatuteDlt13761{}
	err := dlt13761Statute.DecodeByStr(frame)
	if err != nil {
		fmt.Println("err: " + err.Error())
		return
	}
	fmt.Printf("%#v\n", dlt13761Statute)

	//cd()
	//reset()
	//link()
	cmd()
}

func cmd() {
	frameArr, err := dlt1376_1.CreateRelayStationCommandDownF234("01020304", 1, 2)
	if err != nil {
		fmt.Println("err: " + err.Error())
		return
	}
	fmt.Printf(hex.EncodeToString(frameArr))
}

func link() {
	frameArr, err := dlt1376_1.CreateLink("01020304", 3, 1)
	if err != nil {
		fmt.Println("err: " + err.Error())
		return
	}
	fmt.Printf(hex.EncodeToString(frameArr))
}

func reset() {
	frameArr, err := dlt1376_1.CreateResetDown("01020304", 0, 11, 1, 0, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, true)
	if err != nil {
		fmt.Println("err: " + err.Error())
		return
	}
	fmt.Printf(hex.EncodeToString(frameArr))
}

func cd() {
	frameArr, err := dlt1376_1.CreateConfirmDenyF1OrF2("01020304", dlt1376_1.UP, 0, 11, 1, 0, true)
	if err != nil {
		fmt.Println("err: " + err.Error())
		return
	}
	fmt.Printf(hex.EncodeToString(frameArr))
}
