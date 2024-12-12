package iputil

import (
	"fmt"
	"net"
)

// IsValidIPv6 验证IPv6地址
func IsValidIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() == nil && parsedIP.To16() != nil
}

// GetLocalIPv6 获取本地IPv6地址
func GetLocalIPv6() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ip := ipnet.IP.To16(); ip != nil && ip.To4() == nil {
					ipStr := ip.String()
					if IsValidIPv6(ipStr) {
						return ipStr, nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("no valid IPv6 address found")
}
