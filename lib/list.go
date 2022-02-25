package lib

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

const (
	UP    = "UP"
	DOWN  = "DOWN"
	ERROR = "ERROR"
)

func getStatusMessage(status bool, err error) string {
	if err != nil {
		msg := err.Error()

		if strings.Contains(msg, "i/o timeout") {
			return UP
		}

		if strings.Contains(msg, "read: connection refused") {
			return DOWN
		}

		return ERROR
	}

	if status {
		return UP
	} else {
		return DOWN
	}
}

func GetStatuses(sl ServerList) []ServerListStatus {
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

	return statuses
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

	statuses := GetStatuses(sl)

	asjson, _ := json.MarshalIndent(statuses, "", "  ")
	fmt.Println(string(asjson))

	return len(sl.Servers), nil
}
