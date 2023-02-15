package main

import (
	"fmt"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

func main() {
	// 连接到服务器端
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Fail to connect: ", err)
		return
	}
	defer conn.Close()

	// 创建一个ProtoBuf消息
	req := &Msg{
		Data:   []byte("hello world"),
		Time:   time.Now().Unix(),
		FromIp: "1.1.1.1",
		ToIp:   "2.2.2.2",
	}

	data, err := proto.Marshal(req)
	if err != nil {
		fmt.Println("Fail to marshal: ", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Fail to send: ", err)
		return
	}

	// 读取响应数据
	respBuf := make([]byte, 1024)
	n, err := conn.Read(respBuf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	// 反序列化响应数据
	resp := &Msg{}
	err = proto.Unmarshal(respBuf[:n], resp)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err.Error())
		return
	}

	// 打印响应结果
	fmt.Println(resp.Data)
}
