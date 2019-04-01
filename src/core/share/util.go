package share

import (
	"core/log"
	"errors"
	"net"
)

func GetIPByInterface(name string) (string, error) {
	ni, err := net.InterfaceByName(name)
	if err != nil {
		return "", err
	}
	if ni == nil {
		return "", errors.New("net interface not found:" + name)
	}
	addrs, err := ni.Addrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.To4().String(), nil
			}
		}
	}
	return "", errors.New("net interface not found:" + name)
}

func CheckFatalErr(msg string, err error) {
	if err != nil {
		log.Fatal(msg, log.NamedError("err", err))
	}
}
