package main

import (
	"html/template"
	"net/http"
)

// Server ...
type Server struct {
	Name string
}

// ServerList ...
type ServerList struct {
	Servers []Server
}

func main() {
	tmpl := template.Must(template.ParseFiles("template.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := ServerList{
			Servers: []Server{
				{Name: "Server 1"},
				{Name: "Server 2"},
				{Name: "Server 3"},
			},
		}

		tmpl.Execute(w, data)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
