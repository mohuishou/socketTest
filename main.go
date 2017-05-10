package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type clients struct {
	Lon  string `json:"lon"`
	Lat  string `json:"lat"`
	Fall int    `json:"fall"`
}

type user struct {
	ID    int    `json:"id"`
	Types string `json:"type"`
}

type returnData struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

var clientMap = make(map[int]*list.List)

func main() {
	//建立socket，监听端口
	netListen, err := net.Listen("tcp", "localhost:1024")
	checkError(err)
	defer netListen.Close()

	logs("Waiting for clients")

	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		logs(conn.RemoteAddr().String(), " tcp connect success")

		go handleConnection(conn)
	}
}

//处理连接，获取链接的类型，并且分发
func handleConnection(conn net.Conn) {

	buffer := make([]byte, 2048)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			logs(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		u := new(user)
		json.Unmarshal(buffer[:n], u)
		switch u.Types {
		case "app":
			if clientMap[u.ID] == nil {
				clientMap[u.ID] = list.New()
			}
			handleReturn(conn, 1, "app连接成功！")
			go appHandle(conn, clientMap[u.ID])
			return
		case "client":
			if clientMap[u.ID] == nil {
				clientMap[u.ID] = list.New()
			}
			handleReturn(conn, 1, "client连接成功！")
			go clientHandle(conn, clientMap[u.ID])
			return
		default:
			logs("客户端类型错误：", string(buffer[:n]))
			handleReturn(conn, 0, "客户端类型错误！")
			return
		}
	}
}

//对连接返回信息
func handleReturn(conn net.Conn, status int, msg string) {
	words, err := json.Marshal(&returnData{Status: status, Msg: msg})
	checkError(err)
	conn.Write(words)
}

//appHandle 用于处理app的连接
func appHandle(conn net.Conn, cq *list.List) {
	defer conn.Close()
	logs("appHandle收到请求：")
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
func clientHandle(conn net.Conn, cq *list.List) {
	logs("clientHandle收到请求：")
	defer conn.Close()
	buffer := make([]byte, 2048)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			handleReturn(conn, 0, "数据错误！"+err.Error())
			logs(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		c := new(clients)
		json.Unmarshal(buffer[:n], c)

		if c.Lat == "" || c.Lon == "" {
			logs("数据格式错误：", c)
			continue
		}

		cq.PushBack(c)

		//最多缓存十条数据
		if cq.Len() > 10 {
			cq.Remove(cq.Front())
		}
	}
}

func logs(v ...interface{}) {
	log.Println(v...)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
