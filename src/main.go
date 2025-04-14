package main

import (
	"encoding/json"
	"fmt"
	"os"

	server "github.com/boPopov/textprotocol/src/server"
)

var tcpServer server.Server

func init() {
	tcpServer = server.Server{}
}

func main() {
	fmt.Println("Starting Program!")

	tcpServer.Config.Load(os.Args[1])
	tcpServer.Config.Print() //Remove line after testing

	tcpServer.Setup()
	tcpServer.HandleConnections()
}
