package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"monitor/api"
	"monitor/db"
	"monitor/lib"
	"net/http"
	"os"
	"regexp"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron"
)

type App struct {
	Port     string
	Database *sql.DB
	T        *template.Template
}

func (a App) Start(no_cron bool) {
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

	if !no_cron {
		log.Println("Skipping cron.Start() due to getting --offline flag")
		c.Start()
	}

	// web
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/api/servers/", lib.LogReq(a.ApiServers))
	http.Handle("/api/uptime/", lib.LogReq(a.ApiUptimes))
	http.Handle("/api/", lib.LogReq(a.Api))
	http.Handle("/export/", lib.LogReq(a.Export))
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
		Routes: []string{"/api/servers", "/api/uptime/:id"},
	}

	output, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	w.Write(output)
}

func (a App) Export(w http.ResponseWriter, r *http.Request) {
	x, err := ioutil.ReadFile(lib.Env("DB_PATH", "./monitor.db"))

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(x)))
	w.Header().Set("Content-Disposition", "attachment; filename=\"monitor.sqlite3\"")

	w.WriteHeader(http.StatusOK)
	w.Write(x)
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

	// Pull out server id from URL
	// TODO: Put into subroutine
	re := regexp.MustCompile(`\/api\/uptime\/(\d+)`)
	m := re.FindStringSubmatch(r.URL.Path)

	if len(m) != 2 {
		log.Printf("Failed to extract server_id from %s. Returning HTTP 400.", r.URL.Path)

		w.WriteHeader(400)

		return
	}

	server_id, err := strconv.Atoi(m[1])

	if err != nil {
		log.Printf("Failed to convert %s to an int. Returning HTTP 500.", m[1])

		w.WriteHeader(500)

		return
	}

	var data []api.UptimeRow = api.Uptime(a.Database, server_id)

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

	// (Hackily) handle a --no-cron command line arg so we can start the
	// app with server fetchiing and checking off.
	no_cron := false

	if len(args) == 1 && args[0] == "--no-cron" {
		no_cron = true
	}

	// Serve
	app := App{
		Port:     lib.Env("PORT", "8080"),
		Database: database,
	}

	app.Start(no_cron)
}
