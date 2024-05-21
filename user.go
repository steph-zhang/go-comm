package main

import (
	"net"
)

type User struct {
	Name   string
	Addr   string
	Conn   net.Conn
	C      chan string
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		Conn:   conn,
		C:      make(chan string),
		server: server,
	}

	go user.ListenMessage()
	return user
}

func (u *User) Online() {
	u.server.MapLock.Lock()
	u.server.UserMap[u.Name] = u
	u.server.MapLock.Unlock()
	u.server.BroadCast(u, "已上线")
}

func (u *User) Offline() {
	delete(u.server.UserMap, u.Name)
	u.server.BroadCast(u, u.Name+"下线")
}

func (u *User) SendMsg(msg string) {
	if msg == "who" {
		u.server.MapLock.Lock()
		for _, user := range u.server.UserMap {
			u.server.SendMsg(u, user.Name+"在线\n")
		}
		u.server.MapLock.Unlock()
	} else {
		u.server.BroadCast(u, msg)
	}
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.Conn.Write([]byte(msg + "\n"))
	}
}
