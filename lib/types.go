package lib

import (
	"encoding/xml"
)

// Represents a complete server list XML document
type ServerList struct {
	XMLName xml.Name         `xml:"ArrayOfServerItem"`
	Servers []ServerListItem `xml:"ServerItem"`
}

// Represents an entry in the complete server list XML document
type ServerListItem struct {
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

// Represents the minimal information about a server
type Server struct {
	Host string
	Port string
}

// Represents the status of a server
type ServerListStatus struct {
	Name   string
	Status string
}
