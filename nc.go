package main

import (
	"net"
	"log"
	"bufio"
	"os"
	"strings"
	"flag"
	"strconv"
)

var connections [] net.Conn

var localConnection net.Conn

var port *int

var connectionStr string

func handelConn(c net.Conn) {
	defer c.Close()
	for {
		var buf = make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			log.Println("conn read error:", err)
			return
		}
		for _, c := range connections {
			c.Write([]byte(string(buf[:n])))
		}

	}
}

func getInput() {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, err := inputReader.ReadString('\n')
		if err != nil {
			log.Println("read stdin error!")
		}
		if strings.Index(input, "echo") == -1 {
			localConnection.Write([]byte(input))
		}

	}

}

func getLocalOutput() {
	for {
		var buf = make([]byte, 1024)
		n, err := localConnection.Read(buf)
		if err != nil {
			log.Println("local read error:", err)
			return
		}
		print(string(buf[:n]))

	}
}

func main() {

	port = flag.Int("port", 8888, "port")
	connectionStr = ":"
	connectionStr += strconv.Itoa(*port)

	go startServer()
	go localConn()
	getInput()

}

func startServer() {
	l, err := net.Listen("tcp", connectionStr)
	if err != nil {
		log.Println("Error listen:", err)
		return
	}
	defer l.Close()
	log.Println("Server start")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Accept error:", err)
		}
		conn.Write([]byte("Hello\n"))
		connections = append(connections, conn)
		go handelConn(conn)
	}
}

func localConn() {
	localConnection, _ = net.Dial("tcp", connectionStr)
	log.Println("connect to server!")
	go getLocalOutput()
}
