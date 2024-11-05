package dlt698

import "encoding/hex"

type plugIn func() APDURegion

var apduMap = make(map[string]plugIn)

func translate(typeFlag ...byte) APDURegion {
	flag := hex.EncodeToString(typeFlag)
	if pi, ok := apduMap[flag]; ok {
		return pi()
	}
	return nil
}

/*-----------------------------------以下数据类型相关--------------------------------------*/
type dataPlugIn func() DataInter

var dataMap = make(map[byte]dataPlugIn)

func dataTranslate(typeFlag byte) DataInter {
	if pi, ok := dataMap[typeFlag]; ok {
		return pi()
	}
	return nil
}
