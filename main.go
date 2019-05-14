package main

import (
	"flag"
	pl "github.com/nexus166/trimv4/process_lists"
	"os"
)

func main() {
	fpath := flag.String("f", "", "file containing list of IPs")
	furl := flag.String("u", "", "link to list of IPs")
	flag.Parse()

	if *fpath != "" {
		pl.FileGetBanlist(*fpath)
	}
	if *furl != "" {
		pl.HTTPGetBanlist(*furl)
	}
	if os.Args[1] == "-" {
		pl.StdInGetBanlist()
	}
}
