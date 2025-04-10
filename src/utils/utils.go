package utils

import (
	"net"
)

func GetClientIP(address net.Addr) (ip string, err error) {
	ip, _, err = net.SplitHostPort(address.String())
	return
}
