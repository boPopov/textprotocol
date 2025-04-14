package main

import (
	"fmt"
	"os"

	server "github.com/boPopov/textprotocol/src/server"
)

var tcpServer server.Server

func init() {
	tcpServer = server.Server{}
	tcpServer.Config = new(server.ServerConfig)
}

func main() {
	fmt.Println("Starting Program!")

	tcpServer.Config.Load(os.Args[1])
	tcpServer.Config.Print()

	tcpServer.Setup()
	tcpServer.HandleConnections()
}
