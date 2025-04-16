package tests

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/boPopov/textprotocol/src/server"
)

// CONFIG_PATH=/home/bojan/Projects/textprotocol/config.json go test .
func startServer() {
	go func() {
		tcpServer := server.Server{}
		tcpServer.Config = new(server.ServerConfig)
		fmt.Println("CONFIG_PATH:", os.Getenv("CONFIG_PATH"))
		tcpServer.Config.Load(os.Getenv("CONFIG_PATH")) //Load Server Confg
		tcpServer.Setup()
		tcpServer.HandleConnections()
	}()
	time.Sleep(2000 * time.Millisecond)
}

func startConnection(t *testing.T) net.Conn {
	serverConnection, err := net.Dial("tcp", "localhost:4242")
	if err != nil {
		t.Logf("Failed to connect to server: %v", err)
	}
	return serverConnection
}

func readOutput(serverConnection net.Conn, t *testing.T) (reply string, err error) {
	reply, err = bufio.NewReader(serverConnection).ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}
	return
}

func TestServerConnection(t *testing.T) {
	startServer()

	serverConnection := startConnection(t)
	defer serverConnection.Close()

	time.Sleep(1 * time.Second)

	reply, _ := readOutput(serverConnection, t)

	if !strings.Contains(reply, "220 localhost") {
		t.Fail()
		t.Log("The reply from the Server is not '220 localhost'")
	}
}

func TestEhloCommand(t *testing.T) {
	serverConnection := startConnection(t)
	defer serverConnection.Close()

	time.Sleep(1 * time.Second)

	readOutput(serverConnection, t)

	time.Sleep(200 * time.Millisecond)

	name := "test"

	bytesWritten, errCommandEntered := fmt.Fprintf(serverConnection, fmt.Sprintf("EHLO %s\n", name))
	if errCommandEntered != nil {
		t.Fatalf("Failed to send command: %v", errCommandEntered)
	}
	t.Log("Bytes written", bytesWritten)

	reply, _ := readOutput(serverConnection, t)

	if !strings.Contains(reply, fmt.Sprintf("Please to meet you %s", name)) {
		t.Fail()
		t.Logf("Response is not 'Please to meet you %s'", name)
	}

	t.Log(reply)
}

func TestDateCommand(t *testing.T) {
	serverConnection := startConnection(t)
	defer serverConnection.Close()

	time.Sleep(1 * time.Second)

	readOutput(serverConnection, t)
	time.Sleep(200 * time.Millisecond)

	name := "test"

	bytesWritten, errCommandEntered := fmt.Fprintf(serverConnection, fmt.Sprintf("EHLO %s\n", name))
	if errCommandEntered != nil {
		t.Fatalf("Failed to send command: %v", errCommandEntered)
	}
	t.Log("Bytes written", bytesWritten)
	time.Sleep(5 * time.Second)

	readOutput(serverConnection, t)

	bytesWritten, errCommandEntered = fmt.Fprintf(serverConnection, "DATE\n")
	currentTime := time.Now().Format("02/01/2006T15:04")
	if errCommandEntered != nil {
		t.Fatalf("Failed to send command: %v", errCommandEntered)
	}
	t.Log("Bytes written", bytesWritten)
	time.Sleep(5 * time.Second)

	reply, _ := readOutput(serverConnection, t)

	if !strings.Contains(reply, fmt.Sprintf("250 %s", currentTime)) {
		t.Fail()
		t.Logf("Response is not '250 %s'", currentTime)
	}
}

func TestQuitCommand(t *testing.T) {
	serverConnection := startConnection(t)
	defer serverConnection.Close()

	time.Sleep(1 * time.Second)

	readOutput(serverConnection, t)

	time.Sleep(200 * time.Millisecond)

	bytesWritten, errCommandEntered := fmt.Fprintf(serverConnection, "QUIT\n")
	if errCommandEntered != nil {
		t.Fatalf("Failed to send command: %v", errCommandEntered)
	}
	t.Log("Bytes written", bytesWritten)

	reply, _ := readOutput(serverConnection, t)

	if !strings.Contains(reply, "221 Bye!") {
		t.Fail()
		t.Log("Response is not '221 Bye!'")
	}
}

func TestDateCommandWithoutEhlo(t *testing.T) {
	serverConnection := startConnection(t)
	defer serverConnection.Close()

	time.Sleep(1 * time.Second)

	readOutput(serverConnection, t)

	time.Sleep(200 * time.Millisecond)

	bytesWritten, errCommandEntered := fmt.Fprintf(serverConnection, "DATE\n")
	if errCommandEntered != nil {
		t.Fatalf("Failed to send command: %v", errCommandEntered)
	}
	t.Log("Bytes written", bytesWritten)

	reply, _ := readOutput(serverConnection, t)

	if !strings.Contains(reply, "550 Bad state") {
		t.Fail()
		t.Logf("Response is not '550 Bad state'")
	}
}

func TestEhloCommandWithoutName(t *testing.T) {
	serverConnection := startConnection(t)
	defer serverConnection.Close()

	time.Sleep(1 * time.Second)

	readOutput(serverConnection, t)

	time.Sleep(200 * time.Millisecond)

	bytesWritten, errCommandEntered := fmt.Fprintf(serverConnection, "EHLO\n")
	if errCommandEntered != nil {
		t.Fatalf("Failed to send command: %v", errCommandEntered)
	}
	t.Log("Bytes written", bytesWritten)

	t.Log("Before readOutput")
	reply, err := bufio.NewReader(serverConnection).ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}
	t.Log("Reply", reply)

	if !strings.Contains(reply, "550 Invalid EHLO command. The name is missing!") {
		t.Fail()
		t.Logf("Response is not '550 Invalid EHLO command. The name is missing!'")
	}
}

func TestRateLimitConnections(t *testing.T) {
	connectionNumberSlice := make([]net.Conn, 0)
	for connectionNumber := 0; connectionNumber < 6; connectionNumber++ {
		serverConnection := startConnection(t)
		time.Sleep(1 * time.Second)
		if connectionNumber == 5 {
			_, err := bufio.NewReader(serverConnection).ReadString('\n')
			if err == nil {
				t.Fatalf("Failed to read response: %v", err)
			}
		} else {
			connectionNumberSlice = append(connectionNumberSlice, serverConnection)
			time.Sleep(200 * time.Millisecond)
		}
	}

	for _, connection := range connectionNumberSlice {
		defer connection.Close()
		time.Sleep(1 * time.Second)
	}
}

func TestInvalidCommand(t *testing.T) {
	serverConnection := startConnection(t)
	defer serverConnection.Close()

	time.Sleep(1 * time.Second)
	fmt.Println("before reading connection")
	readOutput(serverConnection, t)
	fmt.Println("after reading connection")
	time.Sleep(200 * time.Millisecond)

	bytesWritten, errCommandEntered := fmt.Fprintf(serverConnection, " BANANA\n")
	if errCommandEntered != nil {
		t.Fatalf("Failed to send command: %v", errCommandEntered)
	}
	fmt.Println("Bytes written", bytesWritten)

	reply, errInvalidCommand := readOutput(serverConnection, t)

	if errInvalidCommand != nil {
		fmt.Println("ERROR IS: ", errInvalidCommand)
	}

	if !strings.Contains(reply, "Wrong protocol!") {
		t.Fail()
		t.Logf("Response is not 'Wrong protocol!'")
	}
}
