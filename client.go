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
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("connect server error:", err)
		return nil
	}
	client.conn = conn
	return client

}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) updateName() bool {

	fmt.Print("请输入新的用户名：")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("更新用户名失败：", err)
		return false
	}
	return true

}

func (client *Client) publicChat() {
	var chatMsg string
	fmt.Println(">>>>请输入聊天内容，输入exit退出")

	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("失败：", err)
				return
			}
		}
		chatMsg = ""
		fmt.Println(">>>>请输入聊天内容，输入exit退出")
		fmt.Scanln(&chatMsg)

	}
}

func (client *Client) selectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("查询在线用户失败：", err)
		return
	}
}

func (client *Client) privateChat() {
	var remoteName string
	var chatMsg string
	client.selectUsers()

	fmt.Println(">>>>输入聊天对象的用户名，exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>请输入聊天内容,exit退出")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("Error writing", err)
				break
			}
			chatMsg = ""
			fmt.Println(">>>请继续输入聊天内容,exit退出")
			fmt.Scanln(&chatMsg)
		}
		client.selectUsers()
		fmt.Println(">>>>输入聊天对象的用户名，exit退出")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("\n=========================")
	fmt.Println("1. 公聊频道")
	fmt.Println("2. 私聊频道")
	fmt.Println("3.更新用户名")
	fmt.Println("0. 退出")
	fmt.Println("请选择：")
	fmt.Scan(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("输入有误，请重新输入范围内的数字")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}
		switch client.flag {
		case 1:
			client.publicChat()
			break
		case 2:
			client.privateChat()
			break
		case 3:
			client.updateName()
			break
		}

	}
	fmt.Println("退出聊天室！")
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "服务器IP")
	flag.IntVar(&serverPort, "port", 8888, "服务器端口")
}

func main() {
	flag.Parse()
	client := NewClient("127.0.0.1", serverPort)
	if client == nil {
		fmt.Println(">>>>>连接失败")
		return
	}

	go client.DealResponse()
	fmt.Println(">>>>连接服务器成功。。。")
	client.Run()
}
