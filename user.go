package main

import (
	"net"
	"strings"
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
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		_, ok := u.server.UserMap[newName]
		if ok {
			u.server.SendMsg(u, "此用户名已被占用\n")
		} else {
			u.server.MapLock.Lock()
			delete(u.server.UserMap, u.Name)
			u.server.UserMap[newName] = u
			u.server.MapLock.Unlock()

			u.Name = newName
			u.server.SendMsg(u, "你的用户名已被改为"+u.Name+"\n")
		}
	} else if len(msg) > 3 && msg[:3] == "to|" {
		toUserName := strings.Split(msg, "|")[1]
		if toUserName == "" || len(strings.Split(msg, "|")) <= 2 {
			u.server.SendMsg(u, "消息格式错误\n")
			return
		}
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.server.SendMsg(u, "消息格式错误\n")
			return
		}
		toUser, ok := u.server.UserMap[toUserName]
		if !ok {
			u.server.SendMsg(u, "不存在此用户\n")
			return
		}
		u.server.SendMsg(toUser, "from "+u.Name+": "+content+"\n")

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
