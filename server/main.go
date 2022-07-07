package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

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

type XieyiBody struct {
	Method string
	Ars    []int
}

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

func printAddMany(conn *net.TCPConn) {
	zongchangdu := make([]byte, 4)
	var err error
	for {
		_, err := io.ReadFull(conn, zongchangdu)
		if err != nil {
			break
		}
		l2 := BytesToInt(zongchangdu)
		fmt.Println("客户端说，他要给我发", l2, "个数字，让我算和")
		allbytes := make([]byte, l2*4)
		_, err = io.ReadFull(conn, allbytes)
		if err != nil {
			break
		}
		sum := 0
		var ar1 []int
		for i := 0; i < l2; i++ {
			a1 := BytesToInt(allbytes[i*4 : (i+1)*4])
			ar1 = append(ar1, a1)
			sum += a1
		}
		fmt.Println("算出来了，和是：", sum, "数组是：", ar1)
		conn.Write(IntToBytes(sum))
	}
	if err == io.EOF {
		fmt.Println("完了，完了，芭比Q了", err)
	}
}

func printFuzaStruct(conn *net.TCPConn) {
	zongchangdu := make([]byte, 4)
	var err error
	for {
		_, err := io.ReadFull(conn, zongchangdu)
		if err != nil {
			break
		}
		l2 := BytesToInt(zongchangdu)
		allbytes := make([]byte, l2)
		_, err = io.ReadFull(conn, allbytes)
		if err != nil {
			break
		}
		var fzs XieyiBody
		err = json.Unmarshal(allbytes, &fzs)
		if err != nil {
			break
		}
		fmt.Println("客户端说，他要给我发", l2, "个byte，", len(fzs.Ars), "让我算和")
		gansha := fzs.Method
		res := 0
		switch gansha {
		case "add":
			res = arSums(fzs.Ars)
			break
		case "multi":
			res = arMulti(fzs.Ars)
			break
		default:
			res = -1
		}
		fmt.Println("算出来了，客户端让", gansha, len(fzs.Ars), fzs.Ars, "结果是：", res)
		conn.Write(IntToBytes(res))
	}
	if err == io.EOF {
		fmt.Println("完了，完了，芭比Q了", err)
	}
}

func arSums(a []int) int {
	sum := 0
	for i := 0; i < len(a); i++ {
		sum += a[i]
	}
	return sum
}

func arMulti(a []int) int {
	sum := 1
	for i := 0; i < len(a); i++ {
		sum *= a[i]
	}
	return sum
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
		go printFuzaStruct(conn)
	}
}
