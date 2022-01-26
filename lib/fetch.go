package lib

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
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

	return sl, nil
}
