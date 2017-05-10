package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	logs("启动server")
	go main()
	time.Sleep(time.Second)

	logs("创建app客户端")
	appConn := newClient()
	defer appConn.Close()
	send(appConn, &user{ID: 1, Types: "app"})
	go appTestHandle(appConn)
	time.Sleep(time.Second)

	logs("创建client客户端")
	clientConn := newClient()
	send(clientConn, &user{ID: 1, Types: "client"})
	time.Sleep(time.Second)
	send(clientConn, &clients{Fall: 0, Lat: "sss", Lon: "111"})
	time.Sleep(time.Second)

	defer clientConn.Close()
}

func newClient() net.Conn {
	server := "127.0.0.1:1024"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	return conn
}

func appTestHandle(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			logs(conn.RemoteAddr().String(), "appTestHandle connection error: ", err)
			return
		}
		logs(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))
	}
}

func send(conn net.Conn, c interface{}) {
	words, err := json.Marshal(c)
	checkError(err)
	conn.Write(words)
}
