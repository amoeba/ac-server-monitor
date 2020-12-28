package main

import (
	"fmt"
	"net/http"
	"text/template"
)

// TemplateDataServer ...
type TemplateDataServer struct {
	Info ServerItem
	IsUp bool
}

// TemplateData ...
type TemplateData struct {
	Servers []TemplateDataServer
}

func main() {
	sl := fetch()
	fmt.Println("Fetched", len(sl.Servers), "servers")

	statuses := check(sl)

	// Cron
	// c := cron.New()
	// c.AddFunc("@every 5s", func() { check(sl) })
	// c.Start()

	tmpl := template.Must(template.ParseFiles("resources/template.html"))

	//
	tmplData := new(TemplateData)

	for _, server := range sl.Servers {
		s := new(TemplateDataServer)

		s.Info = server
		s.IsUp = statuses[server.ID]

		news := *s
		tmplData.Servers = append(tmplData.Servers, news)
	}
	//

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
		fmt.Println(isup(connectionstring))
		statuses[server.ID] = isup(connectionstring)
	}

	return statuses
}
