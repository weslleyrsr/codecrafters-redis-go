package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer closeListener(l)

	fmt.Println("Server listening to port 6379")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleMessage(conn)
	}
}

type Type byte

const (
	SimpleString Type = '+'
	Error        Type = '-'
	Integer      Type = ':'
	BulkString   Type = '$'
	Array        Type = '*'
)

func parseResp(respString string) (string, string) {
	var command string = ""
	var value string = ""

	switch respString[0] {
	default:
		command = "Unknown"
	case '+':
		re := regexp.MustCompile(`[^+].*[^(\\r\\n)]`)
		match := re.FindString(respString)
		command = match
	}

	return command, value
}

func handleMessage(conn net.Conn) {
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			closeConnection(conn)
			break
		}

		command, value := parseResp(strings.TrimSpace(netData))

		temp := strings.ToUpper(command)
		fmt.Println(temp)

		if temp == "STOP" {
			closeConnection(conn)
			break
		} else if temp == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		} else if temp == "ECHO" {
			conn.Write([]byte("-" + value + "\r\n"))
		} else {
			conn.Write([]byte("-INVALID COMMAND\r\n"))
		}
	}
}

func closeListener(listener net.Listener) {
	print("closing listener\n")
	listener.Close()
}

func closeConnection(connection net.Conn) {
	print("closing connection\n")
	connection.Close()
}
