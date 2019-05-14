package processlists

import (
	"bufio"
	"net"
	"regexp"
	"strings"

	"github.com/nexus166/trimv4/ipv4"
)

func ParseIPv4AndCIDR(data string) []*net.IPNet {
	var reIPv4 = regexp.MustCompile(`(((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])+)(\/(3[0-2]|[1-2][0-9]|[0-9]))?`)
	//var reIPv4 = regexp.MustCompile(`(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/(3[0-2]|[1-2][0-9]|[0-9]))`)
	scanner := bufio.NewScanner(strings.NewReader(data))

	addrs := make([]*net.IPNet, 0)
	for scanner.Scan() {
		x := reIPv4.FindString(scanner.Text())
		if !strings.Contains(x, "/") {
			if !strings.Contains(x, ":") {
				x = x + "/32"
			} else {
				x = x + "/128"
			}
		}
		if addr, cidr, e := net.ParseCIDR(x); e == nil {
			if !ipv4.IsRFC1918(addr) && !ipv4.IsRFC4193(addr) && !ipv4.IsLoopback(addr) && !ipv4.IsBogonIP(addr) {
				addrs = append(addrs, cidr)
			}
		}
	}
	return addrs
}
