package utils

import (
	"errors"
	"net"
	"strconv"
)

var ErrNoIpAddress error = errors.New("Cannot retrieve node ip address.")

func GetPort() (string, error) {
	port := "5000"
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return port, err
	}

	l, err := net.ListenTCP("tcp", addr)
	defer l.Close()

	if err != nil {
		return port, err
	}

	port = strconv.Itoa(l.Addr().(*net.TCPAddr).Port)

	return port, err
}

func GetHostIp() (string, error) {
	addrs, _ := net.InterfaceAddrs()
	for _, address := range addrs {

		// check the address type and if it is not a loopback then display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}
	}

	return "", ErrNoIpAddress
}
