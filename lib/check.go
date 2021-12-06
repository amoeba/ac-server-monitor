package lib

import "fmt"

func Check(sl ServerList) map[string]bool {
	statuses := make(map[string]bool)

	for _, server := range sl.Servers {
		connectionstring := server.Host + ":" + server.Port

		fmt.Println(connectionstring)
		fmt.Println(IsUp(connectionstring))
		statuses[server.ID] = IsUp(connectionstring)
	}

	return statuses
}
