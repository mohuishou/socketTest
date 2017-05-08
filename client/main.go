package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type Clients struct {
	Lon  string
	Lat  string
	Fall int
}

type user struct {
	ID    int
	Types string
}

func sender(conn net.Conn) {
	c := &user{ID: 1, Types: "client"}
	words, err := json.Marshal(c)
	CheckError(err)
	Log(words)
	conn.Write(words)
	words, err = json.Marshal(&Clients{Fall: 0, Lat: "sss", Lon: "111"})
	conn.Write(words)
	fmt.Println("send over")

}

func main() {
	server := "127.0.0.1:1024"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	// for {
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	handleConnection(conn)
	// }
	fmt.Println("connect success")
	sender(conn)

}

func handleConnection(conn net.Conn) {

	buffer := make([]byte, 2048)
	sender(conn)

	// for {

	n, err := conn.Read(buffer)

	if err != nil {
		Log(conn.RemoteAddr().String(), " connection error: ", err)
		return
	}
	Log(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))

	// }

}

func Log(v ...interface{}) {
	log.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
