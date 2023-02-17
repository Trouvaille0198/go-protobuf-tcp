package server

import (
	"fmt"
	"go-protobuf-tcp/protos"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

type Server struct {
	// conns []net.Conn
	ip   net.IP
	port int
}

func NewServer(ip net.IP, port int) *Server {
	return &Server{
		ip:   ip,
		port: port,
	}
}

func (s *Server) Listen() {
	// 阻塞监听
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
	if err != nil {
		log.Fatal("Error listening: ", err.Error())

	}
	fmt.Println("TCP server started")

	// 接受客户端连接
	for {
		conn, err := listener.Accept()
		// s.conns = append(s.conns, conn)
		if err != nil {
			log.Fatal("Error accepting: ", err.Error())
			continue
		}
		// 处理客户端请求
		go s.handleRequest(conn)
	}
}

func (s *Server) handleRequest(conn net.Conn) {
	// defer conn.Close() // response only once
	for {
		// 阻塞读取数据
		rcvBuf := make([]byte, 1024)
		n, err := conn.Read(rcvBuf)
		if err == io.EOF {
			fmt.Printf("%v: Client (%s) closed connection \n", time.Now().Format("2006-01-02 15:04:05"), conn.RemoteAddr().(*net.TCPAddr))
			conn.Close()
			break
		}
		if err != nil {
			log.Fatal("Error receiving: ", err.Error())
		}
		// 反序列化
		msg := &protos.Msg{}
		err = proto.Unmarshal(rcvBuf[:n], msg)
		if err != nil {
			fmt.Println("Error unmarshalling request: ", err.Error())
			return
		}
		// print msg
		fmt.Printf("%v (%s): %s \n", time.Unix(msg.Time, 0).Format("2006-01-02 15:04:05"), msg.SrcAddr, msg.Data)
		// 倒序
		dataCopy := make([]byte, len(msg.Data))
		copy(dataCopy, msg.Data)
		sendData := []rune(string(dataCopy))
		for i, j := 0, len(sendData)-1; i < j; i, j = i+1, j-1 {
			sendData[i], sendData[j] = sendData[j], sendData[i]
		}
		// 处理请求
		resMsg := &protos.Msg{
			Data:    []byte(string(sendData)),
			Time:    time.Now().Unix(),
			SrcAddr: fmt.Sprintf("%s:%d", s.ip, s.port),
			DstAddr: msg.SrcAddr,
			Type:    protos.TYPE_FROM_SERVER,
		}
		// 序列化响应数据
		resBuf, err := proto.Marshal(resMsg)
		if err != nil {
			log.Fatal("Error marshalling response:", err.Error())
			return
		}
		// 发送响应数据
		_, err = conn.Write(resBuf)
		if err != nil {
			log.Fatal("Error sending:", err.Error())
			return
		}
	}
}
