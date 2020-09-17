package main

import (
	"fmt"
	"net/http"
	"text/template"
)

// TemplateData ...
type TemplateData struct {
	ServerList ServerList
	Statuses   map[string]bool
}

func main() {
	sl := fetch()
	fmt.Println("Fetched", len(sl.Servers), "servers")

	statuses := check(sl)
	fmt.Println(statuses)

	// Cron
	// c := cron.New()
	// c.AddFunc("@every 5s", func() { check(sl) })
	// c.Start()

	tmpl := template.Must(template.ParseFiles("template.html"))

	tmplData := TemplateData{}
	tmplData.Statuses = statuses
	tmplData.ServerList = sl

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, tmplData)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func check(sl ServerList) map[string]bool {
	statuses := make(map[string]bool)

	for _, server := range sl.Servers {
		connectionstring := server.Host + ":" + server.Port

		fmt.Println(connectionstring)
		// fmt.Println(isup(connectionstring))
		statuses[server.ID] = true // temporary
	}

	return statuses
}
