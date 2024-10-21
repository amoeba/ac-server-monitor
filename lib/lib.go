package lib

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/text/message"
)

func Env(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func LogReq(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s", r.URL.Path)

		f(w, r)
	})
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	// This is inefficient - it reads the templates from the filesystem every
	// time. This makes it much easier to develop though, so we can edit our
	// templates and the changes will be reflected without having to restart
	// the app.
	t, err := template.ParseGlob("templates/*.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err.Error()), 500)
		return
	}

	err = t.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err.Error()), 500)
		return
	}
}

func StaticHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static", http.FileServer(http.Dir(dir))).ServeHTTP(w, r)
	}
}

func RelativeTime(ts int64) string {
	now := time.Now().Unix()
	diff := now - ts

	if diff == 0 {
		return "just now"
	} else if diff == 1 {
		return "a second ago"
	} else if diff >= 2 && diff <= 59 {
		return fmt.Sprintf("%d seconds ago", diff)
	} else if diff >= 60 && diff <= 119 {
		return "a minute ago"
	} else if diff >= 120 && diff <= 3540 {
		return fmt.Sprintf("%d minutes ago", diff/60)
	} else if diff >= 3541 && diff <= 7100 {
		return "an hour ago"
	} else if diff >= 7101 && diff <= 82800 {
		return fmt.Sprintf("%d hours ago", (diff+99)/3600)
	} else if diff >= 82801 && diff <= 172000 {
		return "a day ago"
	} else if diff >= 172001 && diff <= 518400 {
		return fmt.Sprintf("%d days ago", (diff+800)/(60*60*24))
	} else if diff >= 518400 && diff <= 1036800 {
		return "a week ago"
	} else {
		return fmt.Sprintf("%d weeks ago", (diff+180000)/(60*60*24*7))
	}
}

func CommafyNumber(number int64) string {
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprintf("%d", number)
}

func BufferToPrettyString(buf []byte) string {
	var builder strings.Builder
	// We know how much space we need ahead of time
	// `0x00 ` for each byte, -1 since no trailer
	builder.Grow(len(buf)*4 - 1)

	for i, b := range buf {
		builder.WriteString(fmt.Sprintf("0x%02X", b))

		// Print space unless we're at the last position
		if i+1 < len(buf) {
			builder.WriteByte(' ')
		}
	}

	return builder.String()
}
