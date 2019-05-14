package ipv4

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sort"
)

type CIDRBlockIPv4 struct {
	first uint32
	last  uint32
}

type CIDRBlockIPv4s []*CIDRBlockIPv4

type IPNets []*net.IPNet

func BCast4(addr uint32, prefix uint) uint32 {
	return addr | ^Netmask(prefix)
}

func SetBit(addr uint32, bit uint, val uint) uint32 {
	if bit < 0 {
		panic("negative bit index")
	}

	if val == 0 {
		return addr & ^(1 << (32 - bit))
	} else if val == 1 {
		return addr | (1 << (32 - bit))
	} else {
		panic("set bit is not 0 or 1")
	}
}

func UI32ToIPv4(addr uint32) net.IP {
	ip := make([]byte, net.IPv4len)
	binary.BigEndian.PutUint32(ip, addr)
	return ip
}

func IPv4ToUI32(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip)
}

func IPV4Merge(blocks CIDRBlockIPv4s) ([]*net.IPNet, error) {
	sort.Sort(blocks)

	for i := len(blocks) - 1; i > 0; i-- {
		if blocks[i].first <= blocks[i-1].last+1 {
			blocks[i-1].last = blocks[i].last
			if blocks[i].first < blocks[i-1].first {
				blocks[i-1].first = blocks[i].first
			}
			blocks[i] = nil
		}
	}

	var merged []*net.IPNet
	for _, block := range blocks {
		if block == nil {
			continue
		}

		if err := IPv4RangeSplit(0, 0, block.first, block.last, &merged); err != nil {
			return nil, err
		}
	}

	return merged, nil
}

func NewBlockIPv4(ip net.IP, mask net.IPMask) *CIDRBlockIPv4 {
	var block CIDRBlockIPv4
	block.first = IPv4ToUI32(ip)
	prefix, _ := mask.Size()
	block.last = BCast4(block.first, uint(prefix))

	return &block
}

func Netmask(prefix uint) uint32 {
	if prefix == 0 {
		return 0
	}
	return ^uint32((1 << (32 - prefix)) - 1)
}

func (c CIDRBlockIPv4s) Len() int {
	return len(c)
}

func (c CIDRBlockIPv4s) Less(i, j int) bool {
	lhs := c[i]
	rhs := c[j]

	if lhs.last < rhs.last {
		return true
	} else if lhs.last > rhs.last {
		return false
	}

	if lhs.first < rhs.first {
		return true
	} else if lhs.first > rhs.first {
		return false
	}

	return false
}

func (c CIDRBlockIPv4s) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func MergeIPNets(nets []*net.IPNet) ([]*net.IPNet, error) {
	if nets == nil {
		fmt.Println("no IPs detected in file")
		return nil, nil
	}
	if len(nets) == 0 {
		return make([]*net.IPNet, 0), nil
	}

	var block4s CIDRBlockIPv4s
	for _, net := range nets {
		ip4 := net.IP.To4()
		if ip4 != nil {
			block4s = append(block4s, NewBlockIPv4(ip4, net.Mask))
		} else {
			return nil, errors.New("Not implemented")
		}
	}

	merged, err := IPV4Merge(block4s)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return merged, nil
}

func IPv4RangeSplit(addr uint32, prefix uint, lo, hi uint32, cidrs *[]*net.IPNet) error {
	if prefix > 32 {
		return fmt.Errorf("Invalid mask size: %d", prefix)
	}

	bc := BCast4(addr, prefix)
	if (lo < addr) || (hi > bc) {
		return fmt.Errorf("%d, %d out of range for network %d/%d, broadcast %d", lo, hi, addr, prefix, bc)
	}

	if (lo == addr) && (hi == bc) {
		cidr := net.IPNet{IP: UI32ToIPv4(addr), Mask: net.CIDRMask(int(prefix), 8*net.IPv4len)}
		*cidrs = append(*cidrs, &cidr)
		return nil
	}

	prefix++
	lowerHalf := addr
	upperHalf := SetBit(addr, prefix, 1)
	if hi < upperHalf {
		return IPv4RangeSplit(lowerHalf, prefix, lo, hi, cidrs)
	} else if lo >= upperHalf {
		return IPv4RangeSplit(upperHalf, prefix, lo, hi, cidrs)
	} else {
		err := IPv4RangeSplit(lowerHalf, prefix, lo, BCast4(lowerHalf, prefix), cidrs)
		if err != nil {
			return err
		}
		return IPv4RangeSplit(upperHalf, prefix, upperHalf, hi, cidrs)
	}
}
