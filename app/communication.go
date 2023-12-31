package main

import (
	"fmt"
	"net"
	"strconv"
)

func buildMessage(msg string) string {
	// create a messsage with UserName as custom Header
	header := "USER_NAME=" + UserName
	message := header + "\n" + msg + "\n"
	return message
}

func connectToClient() {
	fmt.Println("Function to connect to each client")

	// loop through each client and check i'm connected or not
	// else, connect if it fails, leave him 😏
	if len(All_users) == 0 {
		println("No clients to connect ")
		return
	}

	for _, eachIp := range All_users {
		println("eachIp from comm ", eachIp)
		clientAddr := eachIp + ":" + strconv.Itoa(Port)
		tcpAddr, err := net.ResolveTCPAddr("tcp", clientAddr)
		if err != nil {
			println("ResolveTCPAddr failed:", err.Error())
			continue
		}

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			println("Dial failed:", err.Error())
			continue
		}
		_, err = conn.Write([]byte(buildMessage("hi there")))
		if err != nil {
			println("Write to server failed:", err.Error())
			continue
		}

		println("write to server = ", "Hii there \n")
	}
}
