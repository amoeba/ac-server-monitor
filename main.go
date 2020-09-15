package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/revel/cron"
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
	// Cron
	c := cron.New()
	c.AddFunc("* * * * * *", func() { fmt.Println("hi") })
	c.Start()

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
