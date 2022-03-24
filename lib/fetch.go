package lib

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const serverListURL = "https://raw.githubusercontent.com/acresources/serverslist/master/Servers.xml"

// Fetch fetches the public server list and returns a ServerList
func Fetch() (ServerList, error) {
	resp, err := http.Get(serverListURL)

	if err != nil {
		return ServerList{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return ServerList{}, err
	}

	sl := ServerList{}

	if err := xml.Unmarshal(body, &sl); err != nil {
		return ServerList{}, err
	}

	var b strings.Builder

	for i := range sl.Servers {
		b.WriteString(sl.Servers[i].Name)
		b.WriteString(", ")
	}

	log.Printf("Fetched server list: %s", b.String())

	return sl, err
}
