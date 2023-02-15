package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"go-protobuf-tcp/protos"

	"google.golang.org/protobuf/proto"
)

type Client struct {
	conn net.Conn
	ip   net.IP
	port int
}

func NewClient(ip net.IP, port int) *Client {
	// if ip == nil {
	// 	ip = net.ParseIP("127.0.0.1")
	// }
	return &Client{
		ip:   ip,
		port: port,
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.ip, c.port))
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// TCP 客户端
func (c *Client) Start() {
	defer c.conn.Close() // 关闭TCP连接
	inputReader := bufio.NewReader(os.Stdin)
	for {
		// 读取用户输入
		input, err := inputReader.ReadString('\n')
		input = strings.TrimSpace(input)
		if err != nil {
			return
		}
		if strings.ToLower(input) == "q" {
			return
		}

		msg := &protos.Msg{
			Data:  []byte(input),
			Time:  time.Now().Unix(),
			SrcIp: "?",
			DstIp: "?",
		}
		buf, err := proto.Marshal(msg)
		if err != nil {
			fmt.Println("Fail to marshal: ", err)
			return
		}

		_, err = c.conn.Write(buf) // send msg
		if err != nil {
			return
		}

		rcvBuf := make([]byte, 1024)
		n, err := c.conn.Read(rcvBuf[:])
		if err != nil {
			fmt.Println("recv failed, err:", err)
			return
		}
		fmt.Println(string(rcvBuf[:n]))
	}
}
