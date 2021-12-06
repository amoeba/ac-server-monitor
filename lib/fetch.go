package lib

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Fetch() ServerList {
	resp, err := http.Get("https://raw.githubusercontent.com/acresources/serverslist/master/Servers.xml")

	if err != nil {
		fmt.Println("error")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	sl := ServerList{}

	if err := xml.Unmarshal(body, &sl); err != nil {
		log.Fatal(err)
	}

	return sl
}
