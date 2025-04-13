package test

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/boPopov/textprotocol/src/server"
)

func startServer() {
	go func() {
		tcpServer := server.Server{}
		tcpServer.Port = "4242"
		tcpServer.Setup()
		tcpServer.HandleConnections()
	}()
	time.Sleep(2000 * time.Millisecond)
}

func startConnection(t *testing.T) net.Conn {
	serverConnection, err := net.Dial("tcp", "localhost:4242")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	return serverConnection
}

func TestServerConnection(t *testing.T) {
	startServer()

	serverConnection := startConnection(t)
	defer serverConnection.Close()

	time.Sleep(10 * time.Second)

	reply, err := bufio.NewReader(serverConnection).ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	if !strings.Contains(reply, "220 localhost") {
		t.Fail()
		t.Log("The reply from the Server is not '220 localhost'")
	}
}

func TestEhloCommand(t *testing.T) {
	serverConnection := startConnection(t)
	defer serverConnection.Close()

	time.Sleep(10 * time.Second)

	reply, err := bufio.NewReader(serverConnection).ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	time.Sleep(1 * time.Second)

	name := "test"

	bytesWritten, errCommandEntered := fmt.Fprintf(serverConnection, fmt.Sprintf("EHLO %s\n", name))
	if errCommandEntered != nil {
		t.Fatalf("Failed to send command: %v", err)
	}
	t.Log("Bytes written", bytesWritten)
	time.Sleep(5 * time.Second)

	reply, err = bufio.NewReader(serverConnection).ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	if !strings.Contains(reply, fmt.Sprintf("Please to meet you %s", name)) {
		t.Fail()
		t.Logf("Response is not 'Please to meet you %s'", name)
	}

	t.Log(reply)
}

func TestDateCommand(t *testing.T) {
	serverConnection := startConnection(t)
	defer serverConnection.Close()

	time.Sleep(10 * time.Second)

	reply, err := bufio.NewReader(serverConnection).ReadString('\n') //moving from first response.
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	time.Sleep(1 * time.Second)

	name := "test"

	bytesWritten, errCommandEntered := fmt.Fprintf(serverConnection, fmt.Sprintf("EHLO %s\n", name))
	if errCommandEntered != nil {
		t.Fatalf("Failed to send command: %v", err)
	}
	t.Log("Bytes written", bytesWritten)
	time.Sleep(5 * time.Second)

	reply, err = bufio.NewReader(serverConnection).ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	bytesWritten, errCommandEntered = fmt.Fprintf(serverConnection, "DATE\n")
	currentTime := time.Now().Format("02/01/2006T15:04")
	if errCommandEntered != nil {
		t.Fatalf("Failed to send command: %v", err)
	}
	t.Log("Bytes written", bytesWritten)
	time.Sleep(5 * time.Second)

	reply, err = bufio.NewReader(serverConnection).ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	if !strings.Contains(reply, fmt.Sprintf("250 %s", currentTime)) {
		t.Fail()
		t.Logf("Response is not '250 %s'", currentTime)
	}
}

func TestQuitCommand(t *testing.T) {
	serverConnection := startConnection(t)
	defer serverConnection.Close()

	time.Sleep(10 * time.Second)

	reply, err := bufio.NewReader(serverConnection).ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	time.Sleep(1 * time.Second)

	bytesWritten, errCommandEntered := fmt.Fprintf(serverConnection, "QUIT")
	if errCommandEntered != nil {
		t.Fatalf("Failed to send command: %v", err)
	}
	t.Log("Bytes written", bytesWritten)
	time.Sleep(5 * time.Second)

	reply, err = bufio.NewReader(serverConnection).ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

}
