package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"monitor/api"
	"monitor/db"
	"monitor/lib"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron"
)

type App struct {
	Port     string
	Database *sql.DB
	T        *template.Template
}

func (a App) Start() {
	// migrate
	migrate_error := db.AutoMigrate(a.Database)

	if migrate_error != nil {
		log.Fatalf("Error in AutoMigrate: %s", migrate_error)
	}

	// cron
	c := cron.New()

	c.AddFunc("@every 10m", func() {
		lib.Update(a.Database)
	})

	c.Start()

	// web

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/api/servers/", lib.LogReq(a.ApiServers))
	http.Handle("/api/uptime/", lib.LogReq(a.ApiUptimes))
	http.Handle("/api/", lib.LogReq(a.Api))
	http.Handle("/about/", lib.LogReq(a.About))
	http.Handle("/static/", lib.LogReq(lib.StaticHandler("static")))
	http.Handle("/", lib.LogReq(a.Index))

	addr := fmt.Sprintf(":%s", a.Port)

	log.Printf("Starting app on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (a App) About(w http.ResponseWriter, r *http.Request) {
	lib.RenderTemplate(w, "about.html", nil)
}

func (a App) Api(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := struct {
		Routes []string `json:"routes"`
	}{
		Routes: []string{"/api/servers"},
	}

	output, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	w.Write(output)
}

func (a App) ApiServers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data []api.ServerAPIResponse = api.Servers(a.Database)

	output, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	w.Write(output)
}

func (a App) ApiUptimes(w http.ResponseWriter, r *http.Request) {
	log.Println("ApiUptime")
	w.Header().Set("Content-Type", "application/json")

	var data []api.UptimeRow = api.Uptime(a.Database, 1)

	for i, s := range data {
		log.Println(i)
		log.Println(s.Date)
		log.Println(s.Ratio)
	}

	output, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	w.Write(output)
}

func (a App) Index(w http.ResponseWriter, r *http.Request) {
	var servers []api.ServerAPIResponse = api.Servers(a.Database)

	data := struct {
		Servers []api.ServerAPIResponse
	}{
		Servers: servers,
	}

	lib.RenderTemplate(w, "index.html", data)
}

func main() {
	// Logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// DB
	database, err := sql.Open("sqlite3", lib.Env("DB_PATH", "./monitor.db"))

	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()

	// Serve (default) or handle args
	args := os.Args[1:]

	if len(args) == 1 && args[0] == "update" {
		lib.Update(database)

		return
	}

	// Serve

	app := App{
		Port:     lib.Env("PORT", "8080"),
		Database: database,
	}

	app.Start()
}
