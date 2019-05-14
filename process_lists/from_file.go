package processlists

import (
	"fmt"
	"io/ioutil"

	"github.com/nexus166/trimv4/ipv4"
)

func FileGetBanlist(fpath string) {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Println(err)
		return
	}
	merged, e := ipv4.MergeIPNets(ParseIPv4AndCIDR(string(b)))
	if e == nil {
		for ip := range merged {
			fmt.Println(merged[ip])
		}
	}
}
