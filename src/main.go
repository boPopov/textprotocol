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

func readJson() (port string, errEncountered error) {
	file, err := os.Open("../config.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		errEncountered = err
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	errEncountered = decoder.Decode(&port)
	if err != nil {
		fmt.Println("Could not decode the JSON file:", err)
		return
	}
	return
}

func main() {
	fmt.Println("Starting Instance!")

	port, err := readJson()
	if err != nil {
		fmt.Println("Could not extra the port, because of", err)
	}
	fmt.Println("Port is: ", port)

	tcpServer.Port = "4242"
	tcpServer.Setup()
	tcpServer.HandleConnections()
}
