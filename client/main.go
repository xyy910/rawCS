package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8090")
	if err != nil {
		log.Fatalln("第一步就出错啦！", err)
	}
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		log.Fatalln("建立连接出错啦：", err)
	}

	//fmt.Println("发送字符串：时代少年团-渐暖")
	//n, err := conn.Write([]byte("时代少年团-渐暖"))
	//if err != nil {
	//	log.Fatalln("发送出错啦!", err)
	//}
	//fmt.Println("客户端发送了：", n)

	var buf []byte
	fmt.Println("连续发整数1~10")
	for i := 1; i <= 10; i++ {
		buf = append(buf, IntToBytes(i)...)
	}
	for i := 0; i < len(buf); i++ {
		conn.Write(buf[i : i+1])
	}

}
