package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int // operation flag
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn

	return client
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1. Public chat")
	fmt.Println("2. Private chat")
	fmt.Println("3. Update the user name")
	fmt.Println("0. Exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("Invalid input, please try again.")
		return false
	}
}

// Search the online users
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	if _, err := client.conn.Write([]byte(sendMsg)); err != nil {
		fmt.Println("conn.Write:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()

	fmt.Println(">>>>>Enter the user name you want to chat with [UserName], exit to quit:")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>Enter your message (type 'exit' to quit):")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			// Send the private message if not empty
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				if _, err := client.conn.Write([]byte(sendMsg)); err != nil {
					fmt.Println("conn.Write:", err)
					return
				}
			}
			chatMsg = ""
			fmt.Println(">>>>>Enter your message (type 'exit' to quit):")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUsers()
		fmt.Println(">>>>>Enter the user name you want to chat with [UserName], exit to quit:")
		fmt.Scanln(&remoteName)
	}

}

func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>>Enter your message (type 'exit' to quit):")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {

		// Send the message if not empty
		if len(chatMsg) > 0 {
			sendMsg := chatMsg + "\n"
			if _, err := client.conn.Write([]byte(sendMsg)); err != nil {
				fmt.Println("conn.Write:", err)
				return
			}
		}

		chatMsg = ""
		fmt.Println(">>>>>Enter your message (type 'exit' to quit):")
		fmt.Scanln(&chatMsg)
	}

}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>Enter your new user name:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	if _, err := client.conn.Write([]byte(sendMsg)); err != nil {
		fmt.Println("conn.Write:", err)
		return false
	}

	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for !client.menu() {
			// Loop until a valid input is received
		}

		// Handle the user's choice
		switch client.flag {
		case 1:
			// Public chat
			client.PublicChat()
		case 2:
			// Private chat
			client.PrivateChat()
		case 3:
			// Update the user name
			client.UpdateName()
		}

	}
}

// Handle the server's response
func (client *Client) DealResponse() {
	// Copy server response to standard output
	// Block forever
	io.Copy(os.Stdout, client.conn)

	// Same as
	// for
	// {
	// 		buf := make([]byte, 4096)
	// 		n, err := client.conn.Read(buf)
	// 		fmt.Print(buf)
	// }
}

var serverIP string
var serverPort int

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "Server IP address")
	flag.IntVar(&serverPort, "port", 8888, "Server port")
}

func main() {
	flag.Parse()
	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println("Failed to connect to server...")
		return
	}

	// Create another goroutine to handle server responses
	go client.DealResponse()

	fmt.Println("Successfully connected to server...")

	// clinet logic here
	client.Run()
}
