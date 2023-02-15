package main

import (
	"fmt"
	"net"

	"google.golang.org/protobuf/proto"
)

func main() {
	// 监听TCP端口
	listener, err := net.Listen("tcp", "127.0.0.1:9527")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}

	fmt.Println("TCP server started")

	// 接受客户端连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		// 处理客户端请求
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// 读取客户端请求
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	// 反序列化请求数据
	req := &Msg{}
	err = proto.Unmarshal(buf[:n], req)
	if err != nil {
		fmt.Println("Error unmarshalling request:", err.Error())
		return
	}

	// 处理请求
	resp := &Msg{Data: []byte("Gooood")}

	// 序列化响应数据
	respBytes, err := proto.Marshal(resp)
	if err != nil {
		fmt.Println("Error marshalling response:", err.Error())
		return
	}

	// 发送响应数据
	_, err = conn.Write(respBytes)
	if err != nil {
		fmt.Println("Error writing:", err.Error())
		return
	}
}
