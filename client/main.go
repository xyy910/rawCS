package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
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
		for i := 0; i < 4; i++ {
			buf := make([]byte, 0)
			ar1 := make([]int, 0)
			suijishu := rand.Intn(10)
			fmt.Println("这次要随机发：", suijishu, "个数字，让服务器给我算和")
			buf = append(buf, IntToBytes(suijishu)...)
			for k := 0; k < suijishu; k++ {
				a1 := rand.Intn(10)
				ar1 = append(ar1, a1)
				buf = append(buf, IntToBytes(a1)...)
			}
			buflen := len(buf)
			for j := 0; j < buflen; j++ {
				conn.Write(buf[j : j+1])
			}
			fmt.Println(suijishu, "个数字发完了，分别是：", ar1)
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
