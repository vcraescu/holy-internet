package utils

import (
	"io"
	"net"
	"time"
)

const defaultPort = "80"

func resolveHost(host string) (string, string, error) {
	h, port, err := net.SplitHostPort(host)
	if err == nil {
		host = h
	} else {
		port = defaultPort
	}

	ip, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return "", "", err
	}

	return ip.String(), port, nil
}

func Ping(host string) (bool, error) {
	ip, port, err := resolveHost(host)
	if err != nil {
		return false, err
	}

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), time.Second*3)
	if err != nil {
		return false, err
	}

	if _, err := conn.Read([]byte{}); err == io.EOF {
		conn.Close()
		conn = nil

		return false, err
	}

	return true, nil
}
