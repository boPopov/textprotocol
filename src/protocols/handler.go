package protocols

import (
	"net"
)

func UserProtocolConnectionHandler(connection net.Conn) {
	defer connection.Close()

	conn.Write([]byte("220 localhost"))
	var ehloReceived bool
	reader := bufio.NewReader(connection)

	for {
		
	}
}