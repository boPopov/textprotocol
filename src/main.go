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
	fmt.Println("Starting Instance!")

	tcpServer.Config.Load(os.Args[1]) //Load Server Confg 
	tcpServer.Config.Print() //Print Server Config

	tcpServer.Setup()
	tcpServer.HandleConnections()
}
