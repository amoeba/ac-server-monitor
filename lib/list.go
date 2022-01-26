package lib

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func getStatusMessage(status bool, err error) string {
	if err != nil {
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
		statusMessage := getStatusMessage(Check(srv))
		fmt.Fprintf(w, "%s\t\t%s\t\n", item.Name, statusMessage)
	}

	w.Flush()

	return len(sl.Servers), nil
}
