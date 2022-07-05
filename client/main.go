package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
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

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		var buf []byte
		fmt.Println("连续发整数1~10")
		for i := 1; i <= 10; i++ {
			buf = append(buf, IntToBytes(i)...)
		}
		for i := 0; i < len(buf); i++ {
			conn.Write(buf[i : i+1])
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		buf1 := make([]byte, 4)
		for {
			l1, err := io.ReadFull(conn, buf1)
			if err != nil {
				fmt.Println("出错啦！", err)
				break
			}
			fmt.Println("收到了server回复的：", l1, "个byte, 翻译过来就是：", BytesToInt(buf1))
		}
		wg.Done()
	}()
	wg.Wait()
}
