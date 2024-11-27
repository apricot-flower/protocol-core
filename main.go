package main

import (
	"fmt"
	"protocol-core/dlt1376_1"
)

func main() {
	frame := `68 b6 00 b6 00 68 4a 31 07 01 00 02 04 ec 00 00 04 01 02 01 04 00 64 00 02 05 01 c8 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 0c 12 04 09 21 00 01 16`
	dlt13761Statute, err := dlt1376_1.Decode(frame)
	if err != nil {
		fmt.Println("err: " + err.Error())
		return
	}
	fmt.Printf("%#v\n", dlt13761Statute)
}
