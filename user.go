package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	Conn net.Conn
	C    chan string
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		Conn: conn,
		C:    make(chan string),
	}

	go user.ListenMessage()
	return user
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.Conn.Write([]byte(msg + "\n"))
	}
}
