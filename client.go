package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	return client
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务端ip")
	flag.IntVar(&serverPort, "port", 8080, "设置服务端port")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))

	if err != nil {
		fmt.Println("连接失败\n")
		return
	}

	client.Conn = conn

	fmt.Println("连接成功\n")

	select {}
}
