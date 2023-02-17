package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"go-protobuf-tcp/protos"

	"google.golang.org/protobuf/proto"
)

type Client struct {
	conn    net.Conn
	SrcPort int
	DstIp   net.IP // target ip
	DstPort int    // target port
}

func (c *Client) getLocalAddr() string {
	localAddr := c.conn.LocalAddr().(*net.TCPAddr)
	return localAddr.String()
}

func NewClient(DstIp net.IP, DstPort int, SrcPort int) *Client {
	return &Client{
		DstIp:   DstIp,
		DstPort: DstPort,
		SrcPort: SrcPort,
	}
}

func (c *Client) Connect() error {
	var conn net.Conn
	var err error
	if c.SrcPort == 0 {
		// 使用随机端口连接
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.DstIp, c.DstPort))
	} else {
		// 使用指定端口连接
		localAddr := &net.TCPAddr{Port: c.SrcPort}
		fmt.Println("Using port", localAddr.Port)

		d := net.Dialer{LocalAddr: localAddr}
		conn, err = d.Dial("tcp", fmt.Sprintf("%s:%d", c.DstIp, c.DstPort))
	}

	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// TCP 客户端
func (c *Client) Start() {
	defer c.conn.Close()
	inputReader := bufio.NewReader(os.Stdin)
	for {
		// 读取用户输入
		input, err := inputReader.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading input: ", err)
		}
		input = strings.TrimSpace(input)
		// 序列化
		sendMsg := &protos.Msg{
			Data:    []byte(input),
			Time:    time.Now().Unix(),
			SrcAddr: c.getLocalAddr(),
			DstAddr: fmt.Sprintf("%s:%d", c.DstIp, c.DstPort),
			Type:    protos.TYPE_FROM_CLIENT,
		}
		buf, err := proto.Marshal(sendMsg)
		if err != nil {
			log.Fatal("Error marshalling: ", err)
		}
		// send msg
		_, err = c.conn.Write(buf)
		if err != nil {
			return
		}

		if strings.ToLower(input) == "q" {
			log.Println("Quit.")
			return
		}
		// receive msg
		rcvBuf := make([]byte, 1024)
		n, err := c.conn.Read(rcvBuf[:])
		if err != nil {
			log.Fatal("Error receiving: ", err.Error())
		}
		// 反序列化
		rcvMsg := &protos.Msg{}
		err = proto.Unmarshal(rcvBuf[:n], rcvMsg)
		if err != nil {
			log.Fatal("Error unmarshalling request: ", err.Error())
			return
		}
		// print msg
		fmt.Printf("%v (%s): %s\n", time.Unix(rcvMsg.Time, 0).Format("2006-01-02 15:04:05"), rcvMsg.SrcAddr, rcvMsg.Data)
	}
}
