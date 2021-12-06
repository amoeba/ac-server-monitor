package monitor

import (
	"fmt"
	"monitor/lib"
	"net/http"
	"text/template"
)

func main() {
	sl := lib.Fetch()
	fmt.Println("Fetched", len(sl.Servers), "servers")

	statuses := lib.Check(sl)

	// Cron
	// c := cron.New()
	// c.AddFunc("@every 5s", func() { check(sl) })
	// c.Start()

	tmpl := template.Must(template.ParseFiles("resources/template.html"))
	tmplData := new(lib.TemplateData)

	for _, server := range sl.Servers {
		s := new(lib.TemplateDataServer)

		s.Info = server
		s.IsUp = statuses[server.ID]

		news := *s
		tmplData.Servers = append(tmplData.Servers, news)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, tmplData)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
