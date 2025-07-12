package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// Users map to store connected users
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// Message channel to broadcast messages
	Message chan string
}

// create a server interface
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// Listen for messages and broadcast them to all users
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		// Broadcast message to all online users
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	// handle the connection
	// fmt.Println("Client connected from:", conn.RemoteAddr().String())

	user := NewUser(conn)
	// User is online, Add user to online map
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user

	this.mapLock.Unlock()

	// Broadcast welcome message
	this.Broadcast(user, "has online")

	// BLock handler
	select {}
}

// start the server interface
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

	// Start listening for messages in a separate goroutine
	go this.ListenMessager()

	for {
		// accept
		// This will block so we need to handle it in a goroutine
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// do handler
		go this.Handler(conn)
	}
}
