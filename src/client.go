package main

import (
	"fmt"
	"net"
	"time"
	"bytes"
	"encoding/binary"
	"math/rand"
)

const (
	addr = "127.0.0.1:8080"
)

func main() {
	for i := 0; i < 10000; i++ {
		go f()
	}
	select {

	}
}

func f() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("连接服务端失败:", err.Error())
		return
	}
	fmt.Println("已连接服务器")
	defer conn.Close()
	send(conn)
}

func send(conn net.Conn) {
	for {
		rd :=rand.Intn(16)

		if rd > 0 {
			data := make([]byte, rd)
			var n int = len(data)
			bb := make([]byte, 4)
			binary.BigEndian.PutUint32(bb, uint32(n))
			b := bytes.NewBuffer([]byte{})
			b.Write(bb)
			b.Write(data)
			conn.Write(b.Bytes())

			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}
}
