package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// Create a new user
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// Start listening for messages in a separate goroutine
	go user.ListenMessage()

	return user
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
