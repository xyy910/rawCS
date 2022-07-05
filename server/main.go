package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func printShow(conn *net.TCPConn) {
	buff := make([]byte, 1024)
	var res []byte
	var err error
	for {
		changdu, err := conn.Read(buff)
		res = append(res, buff[0:changdu]...)
		fmt.Println("收到了", changdu, "err是：", err, "现在res的长度是: ", len(res))
		xianyoude := len(res)
		if xianyoude/4 > 0 {
			var res1 []int
			i := 0
			for ; i < xianyoude/4; i++ {
				res1 = append(res1, BytesToInt(res[i*4:(i+1)*4]))
			}
			res = res[i*4:]
			fmt.Println("结果是：", res1)
		} else {
			break
		}
	}
	if err == io.EOF {
		fmt.Println("碰到EOF了 到底了", err)
	}

	fmt.Println("长度不大于0了, err 是： ", err)
}

func printAddTwo(conn *net.TCPConn) {
	buff := make([]byte, 8)
	var err error
	for {
		l1, err := io.ReadFull(conn, buff)
		if err != nil {
			break
		}
		fmt.Println("接收到了：", l1, "个byte")
		v1 := BytesToInt(buff[0:4])
		v2 := BytesToInt(buff[4:])
		sum := v1 + v2
		fmt.Println(v1, v2, "和是：", sum)
		conn.Write(IntToBytes(sum))
	}
	if err == io.EOF {
		fmt.Println("完了，完了，芭比Q了", err)
	}
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func main() {
	address := net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8090,
	}
	listener, err := net.ListenTCP("tcp4", &address)
	if err != nil {
		log.Fatalln("监听出错啦：", err)
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatalln("accept出错啦：", err)
		}
		fmt.Println("咦！有个臭活现上钩了！", conn.RemoteAddr())
		go printAddTwo(conn)
	}
}
