package lib

import "encoding/xml"

type TemplateDataServer struct {
	Info ServerItem
	IsUp bool
}

type TemplateData struct {
	Servers []TemplateDataServer
}

type ServerList struct {
	XMLName xml.Name     `xml:"ArrayOfServerItem"`
	Servers []ServerItem `xml:"ServerItem"`
}

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
