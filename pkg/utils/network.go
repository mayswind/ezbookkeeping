package utils

import (
	"bytes"
	"net"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// GetLocalIPAddressesString returns all local ip address, every ip split by comma
func GetLocalIPAddressesString() (string, error) {
	localAddrs, err := GetLocalIPAddresses()

	if err != nil {
		return "", err
	}

	if len(localAddrs) < 1 {
		return "", errs.ErrGettingLocalAddress
	}

	buff := &bytes.Buffer{}

	for i := 0; i < len(localAddrs); i++ {
		if i > 0 {
			buff.WriteString(",")
		}

		buff.WriteString(localAddrs[i].String())
	}

	return string(buff.Bytes()), nil
}

// GetLocalIPAddresses returns all local ip address object array
func GetLocalIPAddresses() ([]net.IP, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return nil, err
	}

	var localAddrs []net.IP

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.IsLoopback() {
				continue
			}

			if ipnet.IP.IsLinkLocalUnicast() {
				continue
			}

			ip := ipnet.IP.To16()

			if ip != nil {
				localAddrs = append(localAddrs, ip)
			}
		}
	}

	return localAddrs, nil
}
