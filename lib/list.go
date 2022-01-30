package lib

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
)

func getStatusMessage(status bool, err error) string {
	if err != nil {
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

	statuses := []ServerListStatus{}

	var wg sync.WaitGroup

	for _, item := range sl.Servers {
		wg.Add(1)

		go func(item ServerListItem) {
			defer wg.Done()
			srv := Server{Host: item.Host, Port: item.Port}
			checkResult, checkError := Check(srv)

			// Servers that error are assumed down (e.g., timeouts)
			if checkError != nil {
				checkResult = false
			}

			statusMessage := getStatusMessage(checkResult, checkError)
			sls := ServerListStatus{Name: item.Name, Status: statusMessage}
			statuses = append(statuses, sls)
		}(item)
	}

	wg.Wait()

	// Write out a column-aligned table of statuses
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATUS\t")

	for _, s := range statuses {
		fmt.Fprintf(w, "%s\t%s\t\n", s.Name, s.Status)
	}

	w.Flush()

	return len(sl.Servers), nil
}
