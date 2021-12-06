package main

import (
	"flag"
	"fmt"
	"monitor/lib"
	"strings"
)

func main() {
	wordPtr := flag.String("server", "ALL", "TODO")
	flag.Parse()

	fmt.Println("word:", *wordPtr)

	var s = lib.Fetch()

	var found *lib.ServerItem

	for _, server := range s.Servers {
		fmt.Println(server.Name)

		if strings.ToLower(server.Name) == strings.ToLower(*wordPtr) {
			fmt.Println(("FOUND"))
			found = &server
		}
	}

	if found == nil {
		fmt.Println("NOTFOUND")
	} else {
		fmt.Println("FOUND")

		fmt.Println(found.Description)
	}
}
