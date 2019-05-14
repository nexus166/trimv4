package processlists

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nexus166/trimv4/ipv4"
)

func StdInGetBanlist() {
	console, _ := os.Stdin.Stat()
	if (console.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println("Failed to read standard input")
		return
	}

	data, err := ioutil.ReadAll(os.Stdin)
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
