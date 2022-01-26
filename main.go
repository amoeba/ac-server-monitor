// Monitor is a command-line utility that can check whether private Asheron's
// Call servers are online. A single server can be checked,
//
//   ./monitor check play.coldeve.online:9000
//
// Or the entire public list at https://github.com/acresources/serverslist can
// be checked,
//
//   ./monitor list
package main

import (
	"flag"
	"fmt"
	"log"
	"monitor/lib"
	"os"
	"strings"
)

const name = "monitor"

// Test
func printUsage() {
	fmt.Printf("Usage: %s <command> [<args>]\n\n", name)
	fmt.Print("Available commands:\n\n")
	fmt.Print("  check <connection-info>: Check a single server\n\n")
	fmt.Printf("    Example: ./%s play.coldeve.online:9000\n\n", name)
	fmt.Print("  list: Check all servers in the public server list\n\n")
	fmt.Printf("    Example: ./%s list\n", name)
}

func ParseServerInfo(arg string) (lib.Server, error) {
	tokens := strings.Split(arg, ":")

	if len(tokens) != 2 {
		return lib.Server{}, fmt.Errorf("failed to parse '%s'. Try $HOST:$PORT", arg)
	}

	return lib.Server{Host: tokens[0], Port: tokens[1]}, nil
}

func main() {
	args := os.Args[1:]

	if len(args) <= 0 {
		printUsage()

		return
	}

	switch args[0] {
	case "check":
		checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
		checkCmd.Parse(args[1:])
		checkCmdArgs := checkCmd.Args()

		if len(checkCmdArgs) <= 0 {
			printUsage()
			os.Exit(1)
		}

		srv, err := ParseServerInfo(checkCmdArgs[0])

		if err != nil {
			log.Fatal(err)
		}

		_, checkErr := lib.Check(srv)

		if checkErr != nil {
			log.Fatal(checkErr)
		}
	case "list":
		_, listErr := lib.ListServers()

		if listErr != nil {
			log.Fatal(listErr)
		}
	case "serve":
		log.Fatal("Not implemented")
	default:
		printUsage()
		os.Exit(1)
	}
}
