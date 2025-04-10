package server

import (
	"fmt"
	"log"
	"net"
	"time"

	connectionHandler "github.com/boPopov/textprotocol/src/protocols"
	"github.com/boPopov/textprotocol/src/security"
	"github.com/boPopov/textprotocol/src/utils"
)

// Defining the behavior of te Server Structure
type Serverer interface {
	Setup()
	HandleConnections()
	Close()
	CheckRateLimit()
}

type Server struct {
	Listener       net.Listener
	Port           string
	RateLimitPerIp map[string]*security.RateLimit
	Serverer
}

func (server *Server) Setup() {
	var err error
	server.Listener, err = net.Listen("tcp", fmt.Sprintf(":%s", server.Port))
	if err != nil {
		log.Fatalf("Failed to bind to port: %v", err)
	}
	log.Printf("Server is listening on port: %s...", server.Port) //Add new line here.
	server.RateLimitPerIp = make(map[string]*security.RateLimit)
}

func (server *Server) HandleConnections() {
	for {
		connection, err := server.Listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		connection.SetReadDeadline(time.Now().Add(2 * time.Hour)) //Limiting connection to 2 Hours|Might delete later.

		ip, errClientIp := utils.GetClientIP(connection.RemoteAddr())
		if errClientIp != nil {
			fmt.Println("Could not extract the ip from the connection")
			continue
		}

		server.CheckRateLimit(ip)
		if canAllocate := server.RateLimitPerIp[ip].Allocate(); !canAllocate {
			connection.Write([]byte("You have reached the maximum amount of connections"))
			connection.Close()
		} else {
			go connectionHandler.UserProtocolConnectionHandler(connection, server.RateLimitPerIp[ip]) //Add new package that will handle the logic behind the protocols
		}
	}
}

func (server *Server) CheckRateLimit(ip string) {
	if _, exists := server.RateLimitPerIp[ip]; !exists {
		server.RateLimitPerIp[ip] = new(security.RateLimit)
		server.RateLimitPerIp[ip].CreateRateLimiter()
	}
}

func (server *Server) Close() {
	server.Listener.Close()
}
