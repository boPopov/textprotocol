package main

import (
	"fmt"
)

var server server.Server

func init() {
	server = new(server.Server)
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
		error.Printf("Could not extra the port, because of: %v", err)
	}

	server.Port = port
	server.Setup()
	server.HandleConnections()
}