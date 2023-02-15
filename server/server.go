package server

import (
	"fmt"
	"go-protobuf-tcp/protos"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

type Server struct {
	conn net.Conn
	ip   net.IP
	port int
}

func NewServer(ip net.IP, port int) *Server {
	return &Server{
		ip:   ip,
		port: port,
	}
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", "127.0.0.1:9527")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return nil
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
	// defer conn.Close()
	for {

		rcvBuf := make([]byte, 1024)
		n, err := conn.Read(rcvBuf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}

		msg := &protos.Msg{}
		err = proto.Unmarshal(rcvBuf[:n], msg)
		if err != nil {
			fmt.Println("Error unmarshalling request:", err.Error())
			return
		}

		// 处理请求
		resMsg := &protos.Msg{
			Data:  []byte("Goood"),
			Time:  time.Now().Unix(),
			SrcIp: "?",
			DstIp: "?",
		}

		// 序列化响应数据
		resBuf, err := proto.Marshal(resMsg)
		if err != nil {
			fmt.Println("Error marshalling response:", err.Error())
			return
		}

		// 发送响应数据
		_, err = conn.Write(resBuf)
		if err != nil {
			fmt.Println("Error writing:", err.Error())
			return
		}
	}
}
