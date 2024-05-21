package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	UserMap map[string]*User
	MapLock sync.RWMutex
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,

		UserMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}

func (s *Server) Handler(conn net.Conn) {
	user := NewUser(conn, s)

	user.Online()

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("err:", err)
				return
			}

			msg := string(buf[:n-1])
			user.SendMsg(msg)
		}
	}()

	select {}
}

func (s *Server) BroadCast(user *User, msg string) {
	msg = user.Name + " " + msg
	s.Message <- msg
}

func (s *Server) SendMsg(user *User, msg string) {
	user.Conn.Write([]byte(msg))
}

func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.MapLock.Lock()
		for _, user := range s.UserMap {
			user.C <- msg
		}
		s.MapLock.Unlock()
	}
}

func (s *Server) Start() {
	// listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	//close
	defer listener.Close()

	go s.ListenMessage()

	//accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("err:", err)
			continue
		}
		//do handler
		go s.Handler(conn)
	}

}
