package protocols

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/boPopov/textprotocol/src/security"
)

func UserProtocolConnectionHandler(connection net.Conn, rateLimit *security.RateLimit, sessionActiveInterval int, connectionLifeSpanMinutes int) {
	defer connection.Close()

	connection.Write([]byte("220 localhost\n"))
	reader := bufio.NewReader(connection)
	quit := false
	ehloName := ""

	go func() { //Setting the connection maximum lifespan.
		time.Sleep(time.Duration(int64(connectionLifeSpanMinutes) * int64(time.Minute)))
		connection.Close()
	}()

	for {
		connection.SetReadDeadline(time.Now().Add(time.Duration(int64(sessionActiveInterval) * int64(time.Second))))
		inputLine, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Error while reading the input", err)
			rateLimit.Release()
			break
		}

		if allowed := rateLimit.CommandRateLimit.Allow(); !allowed { //Checking if the user has reached the maximum number of commands entered in the dedicated interval.
			connection.Write([]byte("Please slow Down!\n"))
			rateLimit.Release() //Releasing the channel for the Allocated Connection.
			break
		}

		enteredProtocol := strings.TrimSpace(inputLine)

		if strings.Contains(enteredProtocol, "EHLO") {
			if len(enteredProtocol) <= 4 {
				connection.Write([]byte("550 Invalid EHLO command. The name is missing!\n"))
				continue
			}
			splitedProtocol := strings.Split(enteredProtocol, " ")
			ehloName = strings.TrimSpace(splitedProtocol[len(splitedProtocol)-1])
			enteredProtocol = strings.TrimSpace(splitedProtocol[0])
		}

		switch enteredProtocol {
		case "QUIT":
			connection.Write([]byte("221 Bye!\n"))
			quit = true
			rateLimit.Release()
			return
		case "EHLO":
			connection.Write([]byte(fmt.Sprintf("Please to meet you %s\n", ehloName)))
		case "DATE":
			dateResponse := "550 Bad state\n"
			if ehloName != "" {
				dateResponse = fmt.Sprintf("250 %s\n", time.Now().Format("02/01/2006T15:04:05"))
			}
			connection.Write([]byte(dateResponse))
		default:
			connection.Write([]byte("Wrong protocol!\n"))
		}

		if quit {
			break
		}
	}
}
