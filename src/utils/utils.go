package utils

import (
	"net"
	"strings"
)

func GetClientIP(address net.Addr) (ip string, err error) {
	ip, _, err = net.SplitHostPort(address.String())
	return
}

func GetValidCommand(inputStream string) string {
	if strings.Contains(inputStream, "QUIT") {
		return "QUIT"
	} else if strings.Contains(inputStream, "DATE"){
		return "DATE"
	} else {
		return inputStream
	}
}