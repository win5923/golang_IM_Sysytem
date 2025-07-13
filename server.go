package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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

	user := NewUser(conn, this)

	user.Online(this)

	// Listen if the user is active
	isLive := make(chan bool)

	// Recieve messages from user
	go this.ProcessMessage(user, conn, isLive)

	// BLock handler
	for {
		select {
		case <-isLive:
			// User is active, reset the timer
			// Do nothing, just reset the timer

		case <-time.After(time.Second * 300):
			// Overtime

			// Kick the user out

			user.SendMsg("You are kicked out due to inactivity\n")

			// Delete the user from online map
			close(user.C)
			// Close the connection
			conn.Close()

			return
		}
	}
}

// Precess the message from user
func (this *Server) ProcessMessage(user *User, conn net.Conn, isLive chan bool) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			user.Offline(this)
			return
		}

		if err != nil && err != io.EOF {
			fmt.Println("Conn.Read err:", err)
			return
		}

		// Process the message (remove '\n')
		msg := string(buf[:n-1])

		// User's message handling logic
		user.DoMessage(msg)

		// If the user receives a message, it means the user is still active
		isLive <- true
	}
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
