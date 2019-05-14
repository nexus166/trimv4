package processlists

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nexus166/trimv4/ipv4"
)

func HTTPGetBanlist(listurl string) {
	resp, err := http.Get(listurl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	merged, e := ipv4.MergeIPNets(ParseIPv4AndCIDR(string(data)))
	if e == nil {
		for ip := range merged {
			fmt.Println(merged[ip])
		}
	}
}
