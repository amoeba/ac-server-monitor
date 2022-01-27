package lib

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func getStatusMessage(status bool, err error) string {
	if err != nil {
		fmt.Println(err.Error())
		msg := err.Error()

		if strings.Contains(msg, "i/o timeout") {
			return "DOWN"
		}

		if strings.Contains(msg, "read: connection refused") {
			return "DOWN"
		}

		return "ERROR"
	}

	if status {
		return "UP"
	} else {
		return "DOWN"
	}
}

func ListServers() (int, error) {
	sl, err := Fetch()

	if err != nil {
		return -1, err
	}

	if len(sl.Servers) <= 0 {
		fmt.Println("Server list was empty. Exiting")

		return 0, nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATUS\t")

	for _, item := range sl.Servers {
		srv := Server{Host: item.Host, Port: item.Port}
		checkResult, checkError := Check(srv)

		// Servers that error are assumed down (e.g., timeouts)
		if checkError != nil {
			checkResult = false
		}

		statusMessage := getStatusMessage(checkResult, checkError)
		fmt.Fprintf(w, "%s\t%s\t\n", item.Name, statusMessage)
	}

	w.Flush()

	return len(sl.Servers), nil
}
