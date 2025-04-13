package server

import (
	"errors"
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
	CheckRateLimitMap()
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

func (server *Server) HandleConnections() error {
	if server.Listener == nil {
		return errors.New("Server is not Initialized")
	}

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

		server.CheckRateLimitMap(ip)
		if canAllocate := server.RateLimitPerIp[ip].Allocate(); !canAllocate {
			connection.Write([]byte("You have reached the maximum amount of connections"))
			connection.Close()
		} else {
			go connectionHandler.UserProtocolConnectionHandler(connection, server.RateLimitPerIp[ip]) //Add new package that will handle the logic behind the protocols
		}
	}

	return nil
}

/**
 * The Purpose of the function is to check if the IP is present in the RateLimitPerIp map.
 * If the IP is not present in the Map, a new pointer is create of the RateLimit structure and that pointer is set with the default values.
 */
func (server *Server) CheckRateLimitMap(ip string) {
	if _, exists := server.RateLimitPerIp[ip]; !exists {
		server.RateLimitPerIp[ip] = new(security.RateLimit)
		server.RateLimitPerIp[ip].CreateRateLimiter()
	}
}

func (server *Server) Close() {
	server.Listener.Close()
}
