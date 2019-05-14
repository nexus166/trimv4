package ipv4

import (
	"net"
)

func IsRFC1918(ip net.IP) bool {
	for _, cidr := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"} {
		_, pvt, _ := net.ParseCIDR(cidr)
		if pvt.Contains(ip) {
			return true
		}
	}
	return false
}

func IsRFC4193(ip net.IP) bool {
	_, subnet, _ := net.ParseCIDR("fd00::/8")
	return subnet.Contains(ip)
}

func IsLoopback(ip net.IP) bool {
	for _, cidr := range []string{"127.0.0.0/8", "::1/128"} {
		_, lpbk, _ := net.ParseCIDR(cidr)
		if lpbk.Contains(ip) {
			return true
		}
	}
	return false
}

func IsBogonIP(ip net.IP) bool {
	for _, cidr := range []string{
		"100.64.0.0/10",
		"169.254.0.0/16",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"224.0.0.0/3",
		"0000::/8",
		"0200::/7",
		"3ffe::/16",
		"2001:db8::/32",
		"2002:e000::/20",
		"2002:7f00::/24",
		"2002:0000::/24",
		"2002:ff00::/24",
		"2002:0a00::/24",
		"2002:ac10::/28",
		"2002:c0a8::/32",
		"fc00::/7",
		"fe80::/10",
		"fec0::/10",
		"ff00::/8",
	} {
		_, bgn, _ := net.ParseCIDR(cidr)
		if bgn.Contains(ip) {
			return true
		}
	}
	return false
}
