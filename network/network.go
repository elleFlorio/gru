package network

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
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

func DoRequest(method string, path string, body []byte) ([]byte, error) {
	b := bytes.NewBuffer(body)

	req, err := http.NewRequest(method, path, b)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)

	return data, nil
}
