package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// ServerList ...
type ServerList struct {
	XMLName xml.Name     `xml:"ArrayOfServerItem"`
	Servers []ServerItem `xml:"ServerItem"`
}

// ServerItem ...
type ServerItem struct {
	XMLname     xml.Name `xml:"ServerItem"`
	ID          string   `xml:"id"`
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Emu         string   `xml:"emu"`
	Host        string   `xml:"server_host"`
	Port        string   `xml:"server_port"`
	Type        string   `xml:"type"`
	Status      string   `xml:"status"`
	Website     string   `xml:"website_url"`
	Discord     string   `xml:"discord_url"`
}

func fetch() ServerList {
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
