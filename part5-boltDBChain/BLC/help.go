package BLC

import (
	"bytes"
	"encoding/binary"
)

func IntToByte(num int64) []byte {
	buff := new(bytes.Buffer)
	//大写端输入
	binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes()
}
