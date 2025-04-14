package protocols

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/boPopov/textprotocol/src/security"
)

func UserProtocolConnectionHandler(connection net.Conn, rateLimit *security.RateLimit) {
	defer connection.Close()

	connection.Write([]byte("220 localhost\n"))
	reader := bufio.NewReader(connection)
	quit := false
	ehloName := ""

	for {
		inputLine, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error while reading the input", err)
			break
		}

		if allowed := rateLimit.CommandRateLimit.Allow(); !allowed {
			connection.Write([]byte("Please slow Down!\n"))
			rateLimit.Release()
			break
		}

		inputLine = utils.GetValidCommand(inputLine)

		enteredProtocol := strings.TrimSpace(inputLine)

		if strings.Contains(enteredProtocol, "EHLO") {
			if len(enteredProtocol) <= 4 { //Add or section which checks if there is a second command for QUIT and DATE After EHLO
				connection.Write([]byte("550 Invalid EHLO command. The name is missing!\n"))
				continue
			}
			splitedProtocol := strings.Split(enteredProtocol, " ")
			//Test Later
			ehloName = strings.TrimSpace(splitedProtocol[len(splitedProtocol) - 1])
			enteredProtocol = strings.TrimSpace(splitedProtocol[0])

			// if len(splitedProtocol) == 2 {
			// 	ehloName = strings.TrimSpace(splitedProtocol[1])
			// 	enteredProtocol = strings.TrimSpace(splitedProtocol[0])
			// } else {
			// 	for _, value := range splitedProtocol {
			// 		if len(value) > 0 && value != "EHLO" {
			// 			ehloName = value
			// 		}
			// 	}
			// 	enteredProtocol = "EHLO"
			// }
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
