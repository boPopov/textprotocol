package server

import (
	"fmt"
	"log"
	"net"

	connectionHandler "github.com/boPopov/textprotocol/src/protocols"
)

// Defining the behavior of te Server Structure
type Serverer interface {
	Setup()
	HandleConnections()
	Close()
}

type Server struct {
	Listener net.Listener
	Port     string
	Serverer
}

func (server *Server) Setup() {
	var err error
	server.Listener, err = net.Listen("tcp", fmt.Sprintf(":%s", server.Port))
	if err != nil {
		log.Fatalf("Failed to bind to port: %v", err)
	}
	log.Printf("Server is listening on port: %s...", server.Port) //Add new line here.
}

func (server *Server) HandleConnections() {
	for {
		connection, err := server.Listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}

		connectionHandler.UserProtocolConnectionHandler(connection) //Add new package that will handle the logic behind the protocols
		break
	}
	server.Close()
}

func (server *Server) Close() {
	server.Listener.Close()
}
