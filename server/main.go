package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type clients struct {
	Lon  string
	Lat  string
	Fall int
}

type user struct {
	ID    int
	Types string
}

var cq = list.New()

func main() {
	//建立socket，监听端口
	netListen, err := net.Listen("tcp", "localhost:1024")
	checkError(err)
	defer netListen.Close()

	log("Waiting for clients")

	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		log(conn.RemoteAddr().String(), " tcp connect success")

		go handleConnection(conn)
	}
}

//处理连接，获取链接的类型，并且分发
func handleConnection(conn net.Conn) {

	buffer := make([]byte, 2048)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		u := new(user)
		json.Unmarshal(buffer[:n], u)
		switch u.Types {
		case "app":
			go appHandle(conn, u.ID)
			return
		case "client":
			go clientHandle(conn, u.ID)
			return
		default:
			log("客户端类型错误：", string(buffer[:n]))
			return
		}
	}
}

//appHandle 用于处理app的连接
func appHandle(conn net.Conn, id int) {
	defer conn.Close()
	log("appHandle收到请求：")
	//不停的从list读取数据
	for {
		if cq.Len() > 0 {
			if c, ok := cq.Front().Value.(*clients); ok {
				cq.Remove(cq.Front())
				words, err := json.Marshal(c)
				checkError(err)
				conn.Write(words)
			}

		}
	}
}

//clientHandle 用于处理硬件的连接,不断从硬件读取数据
func clientHandle(conn net.Conn, id int) {
	log("clientHandle收到请求：")
	defer conn.Close()
	buffer := make([]byte, 2048)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		c := new(clients)
		json.Unmarshal(buffer[:n], c)

		if c.Lat == "" || c.Lon == "" {
			log("数据格式错误：", c)
			continue
		}

		cq.PushBack(c)
	}
}

func log(v ...interface{}) {
	log.Println(v...)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
