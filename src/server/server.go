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

type Serverer interface {
	Setup()
	HandleConnections()
	Close()
	CheckIPPresence()
}

type Server struct {
	Listener       net.Listener
	Port           string
	RateLimitPerIp map[string]*security.RateLimit
	Config   *ServerConfig
	Serverer
}

func (server *Server) Setup() {
	var err error
	server.Listener, err = net.Listen("tcp", fmt.Sprintf(":%s", server.Config.Port))
	if err != nil {
		log.Fatalf("Failed to bind to port: %v", err)
		panic() // Stopping the program if we can not bind the port specified in the config.json
	}
	log.Printf("Server is listening on port: %s...", server.Config.Port)
	server.RateLimitPerIp = make(map[string]*security.RateLimit)
}

/**
 * This function handles all of the incoming Connections into the application.
 */
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
		connection.SetReadDeadline(time.Now().Add(server.Config.SessionActiveInterval * time.Hour)) //Limiting connection to X Hours (X is specified in the config.json).

		clientIp, errClientIp := utils.GetClientIP(connection.RemoteAddr())
		if errClientIp != nil {
			fmt.Println("Could not extract the ip from the connection")
			continue
		}

		server.CheckIPPresence(clientIp) 
		if canAllocate := server.RateLimitPerIp[clientIp].Allocate(); !canAllocate {
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
func (server *Server) CheckIPPresence(ip string) {
	if _, exists := server.RateLimitPerIp[ip]; !exists {
		server.RateLimitPerIp[ip] = new(security.RateLimit)
		server.RateLimitPerIp[ip].CreateRateLimiter(server.Config.RateLimitMaxSessions, server.RateLimitMaxInputPerInterval, server.RateLimitRefillDuration)
	}
}

func (server *Server) Close() {
	server.Listener.Close()
}
