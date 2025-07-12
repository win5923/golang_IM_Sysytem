package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// Create a new user
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	// Start listening for messages in a separate goroutine
	go user.ListenMessage()

	return user
}

// User's online logic
func (this *User) Online(server *Server) {
	// User is online, Add user to online map
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this

	this.server.mapLock.Unlock()

	// Broadcast welcome message
	this.server.Broadcast(this, "has online")
}

// User's offline logic
func (this *User) Offline(server *Server) {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)

	this.server.mapLock.Unlock()

	this.server.Broadcast(this, "has offline")
}

// User's message handling logic
func (this *User) DoMessage(msg string) {
	// Broadcast the message to all users
	this.server.Broadcast(this, msg)
}

// Listen for channel messages of the user
// Once a message is received, it sends it to the user's connection
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		if _, err := this.conn.Write([]byte(msg + "\n")); err != nil {
			return
		}
	}
}
