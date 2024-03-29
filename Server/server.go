package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const BUFFERSIZE = 1024

func main() {
	server, err := net.Listen("tcp", ":29999")
	if err != nil {
		fmt.Println("Error listetning: ", err)
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("Server started! Waiting for connections...")
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("Client connected")
		go sendFileToClient(connection)
	}
}

func fillStr(retStr string, toLen int) string {
	for {
		lenStr := len(retStr)
		if lenStr < toLen {
			retStr = retStr + ":"
			continue
		}
		break
	}
	return retStr
}

func sendFileToClient(connection net.Conn) {
	fmt.Println("A client has connected!")
	defer connection.Close()
	bufferFileName := make([]byte, 64)
	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")
	fmt.Println("Received file name: ", fileName, "! Opening ...")
	
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fileSize := fillStr(strconv.FormatInt(fileInfo.Size(), 10), 10)
	connection.Write([]byte(fileSize))
	fmt.Println("File Size sent!")
	
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection!")
	return
}