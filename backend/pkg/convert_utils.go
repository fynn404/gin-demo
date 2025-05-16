package pkg

import "strconv"

func ConvertStr2Int(str string) int64 {
	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic("convert str2int error")
	}
	return res
}
