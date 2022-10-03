package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

//错误处理函数,注意：函数名大写开头才能被导出使用
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToHexInt(num int64) []byte {
	//创建一段缓存
	buff := new(bytes.Buffer)
	//采用大端存储将num写入buff，返回报错信息
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
