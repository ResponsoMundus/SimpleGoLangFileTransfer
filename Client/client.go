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
	fileName := os.Args[1]
	//fmt.Println(filename)
	bufferFileName := fillStr(fileName, 64)
	
	serverAddress := os.Args[2]
	
	connection, err := net.Dial("tcp", serverAddress + ":29999")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	connection.Write([]byte(bufferFileName))
	fmt.Println("Connected to server! File name sent!")
	
	bufferFileSize := make([]byte, 10)
	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	fmt.Println("File size received!")
	
	newFile, err := os.Create(fileName)
	
	if err != nil {
		panic(err)
	}
	
	defer newFile.Close()
	var receivedBytes int64
	
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes + BUFFERSIZE) - fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Received file completely!")
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