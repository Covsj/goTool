package ip

import "net"

func GetLocalIP() string {
	netInterfaces, _ := net.Interfaces()
	AddressArray, _ := net.InterfaceAddrs()
	for k, v := range AddressArray {
		if netInterfaces[k].Flags&net.FlagUp == 0 {
			continue
		}
		if ip, ok := v.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}
	return ""
}
