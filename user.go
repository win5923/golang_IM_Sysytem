package main

import (
	"net"
	"strings"
)

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

// Send a message to the user
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// User's message handling logic
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// List all online users
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + " is online...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// Rename user
		newName := strings.Split(msg, "|")[1]

		// Check if the new name already exists
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("This name already exists, please choose another one.\n")
		} else {
			// Rename the user
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			// Notify the user of the successful rename
			this.SendMsg("You have successfully renamed to " + this.Name + "\n")
		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		// Private message to another user
		// Format: to|username|message

		// 1. get the recipient's name
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("The format is incorrect, please use: to|username|message\n")
			return
		}

		// 2. Get the recipient user depending on the name
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("The user " + remoteName + " cannot be found.\n")
			return
		}

		// 3. Send the message to the recipient
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("The message cannot be empty.\n")
			return
		}
		remoteUser.SendMsg(this.Name + " says: " + content + "\n")

	} else {
		this.server.Broadcast(this, msg)
	}
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
