package api

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"gopkg.in/guregu/null.v4"
)

type ServerTableRow struct {
	ID          int
	GUID        string
	Name        string
	Description string
	Emulator    string
	Host        string
	Port        string
	Type        string
	Status      string
	WebsiteURL  string
	DiscordURL  string
	IsListed    int
	CreatedAt   int
	UpdatedAt   int
}

type ServerStatusRow struct {
	ID           int
	GUID         string
	Name         string
	Host         string
	Port         string
	IsListed     bool
	IsOnline     sql.NullBool
	UpdatedAt    int
	LastSeen     sql.NullInt64
}

type ServerAPIResponse struct {
	LastChecked string `json:"last_checked"`
	Count int `json:"count"`
	Servers  []ServerAPIResponseServer `json:"servers"`
}

type ServerAPIResponseServer struct {
	ID      int                      `json:"-"`
	GUID    string                   `json:"guid"`
	Name    string                   `json:"name"`
	Active  bool                     `json:"active"`
	Address ServerAPIResponseAddress `json:"address"`
	Status  ServerAPIResponseStatus  `json:"status"`
}

type ServerAPIResponseAddress struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type ServerAPIResponseStatus struct {
	IsOnline    null.Bool `json:"online"`
	LastSeen    null.String  `json:"last_seen"`
	LastChecked string       `json:"last_checked"`
}

type ServerAPIResponseWithUptime struct {
	ID      int                      `json:"id"`
	GUID    string                   `json:"guid"`
	Name    string                   `json:"name"`
	Active  bool                     `json:"active"`
	Address ServerAPIResponseAddress `json:"address"`
	Status  ServerAPIResponseStatus  `json:"status"`
	Uptime  []UptimeTemplateItem     `json:"uptime"`
}

func Server(db *sql.DB, id int) ServerTableRow {
	var response ServerTableRow

	stmt := `
	SELECT
		id,
		guid,
		name,
		description,
		emu,
		host,
		port,
		type,
		status,
		website_url,
		discord_url,
		is_listed,
		created_at,
		updated_at
	FROM
		servers
	WHERE
		servers.id = ?
	LIMIT 1
	`

	rows, err := db.Query(stmt, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&response.ID,
			&response.GUID,
			&response.Name,
			&response.Description,
			&response.Emulator,
			&response.Host,
			&response.Port,
			&response.Type,
			&response.Status,
			&response.WebsiteURL,
			&response.DiscordURL,
			&response.IsListed,
			&response.CreatedAt,
			&response.UpdatedAt,
		)

		if err != nil {
			log.Fatalf("error scanning %v", err)
		}
	}

	return response
}

func Servers(db *sql.DB) ServerAPIResponse {
	stmt := `
	SELECT
		servers.id,
		servers.guid,
		servers.name,
		servers.host,
		servers.port,
		servers.is_listed,
		servers.is_online,
		servers.updated_at,
		servers.last_seen
	FROM
		servers
	WHERE
		servers.is_listed IS TRUE
	GROUP BY servers.id
	ORDER BY lower(servers.name);
	`

	rows, err := db.Query(stmt)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var statuses []ServerStatusRow
	var status ServerStatusRow

	for rows.Next() {
		err := rows.Scan(
			&status.ID,
			&status.GUID,
			&status.Name,
			&status.Host,
			&status.Port,
			&status.IsListed,
			&status.IsOnline,
			&status.UpdatedAt,
			&status.LastSeen,
		)

		if err != nil {
			log.Println(err)
		} else {
			statuses = append(statuses, status)
		}
	}

	var finalResponse ServerAPIResponse
	var items []ServerAPIResponseServer
	var item ServerAPIResponseServer

	for i := range statuses {
		isOnline := BoolOrNull(statuses[i].IsOnline)
		lastSeenTime := PrettyTimeOrNullString(statuses[i].LastSeen)
		lastChecked := time.Unix(int64(statuses[i].UpdatedAt), 0)
		lastCheckedTime, err := lastChecked.UTC().MarshalText()

		if err != nil {
			log.Fatal(err)
		}

		item = ServerAPIResponseServer{
			ID:     statuses[i].ID,
			GUID:   statuses[i].GUID,
			Name:   statuses[i].Name,
			Active: statuses[i].IsListed,
			Address: ServerAPIResponseAddress{
				Host: statuses[i].Host,
				Port: statuses[i].Port,
			},
			Status: ServerAPIResponseStatus{
				IsOnline:    isOnline,
				LastSeen:    lastSeenTime,
				LastChecked: string(lastCheckedTime),
			},
		}

		items = append(items, item)
	}

	finalResponse.Servers = items
	finalResponse.Count = len(items)
	
	// Set LastChecked if we have servers
	if len(items) > 0 {
		finalResponse.LastChecked = items[0].Status.LastChecked
	} else {
		finalResponse.LastChecked = time.Now().UTC().Format(time.RFC3339)
	}

	return finalResponse
}

func ServersWithUptimes(db *sql.DB) []ServerAPIResponseWithUptime {
	servers := Servers(db)

	var response []ServerAPIResponseWithUptime

	for i := range servers.Servers {
		var server ServerAPIResponseWithUptime

		server.ID = servers.Servers[i].ID
		server.GUID = servers.Servers[i].GUID
		server.Name = servers.Servers[i].Name
		server.Active = servers.Servers[i].Active
		server.Address = servers.Servers[i].Address
		server.Status = servers.Servers[i].Status

		// Add in uptime info
		rows, err := db.Query(QUERY_UPTIME, server.ID)

		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		var uptimes []UptimeTemplateItem

		for rows.Next() {
			var uptime UptimeRow
			var uptimeTmplItem UptimeTemplateItem

			err := rows.Scan(
				&uptime.Date,
				&uptime.Uptime,
				&uptime.N,
				&uptime.RTTMin,
				&uptime.RTTMax,
				&uptime.RTTMean,
			)

			if err != nil {
				log.Fatal(err)
			}

			uptimeTmplItem.Date = uptime.Date
			uptimeTmplItem.Uptime = uptime.Uptime
			uptimeTmplItem.UptimeFmt = fmt.Sprintf("%.3g", uptime.Uptime)
			uptimeTmplItem.UptimeClass = GetUptimeClass(uptime.Uptime)
			uptimeTmplItem.N = uptime.N
			uptimeTmplItem.RTTMin = SQLNullInt64ToString(uptime.RTTMin)
			uptimeTmplItem.RTTMax = SQLNullInt64ToString(uptime.RTTMax)
			uptimeTmplItem.RTTMean = SQLFloat64ToIntString(uptime.RTTMean)

			uptimes = append(uptimes, uptimeTmplItem)
		}

		server.Uptime = uptimes
		response = append(response, server)
	}

	return response
}

func SQLNullInt64ToString(input sql.NullInt64) string {
	if input.Valid {
		return fmt.Sprintf("%d", input.Int64)
	} else {
		return "n/a"
	}
}

func SQLFloat64ToIntString(input sql.NullFloat64) string {
	if input.Valid {
		return fmt.Sprintf("%d", int(math.Round(input.Float64)))
	} else {
		return "n/a"
	}
}

// Convert an sql.NullInt64 into either a pretty datetime string, "n/a", or
// "err"
func PrettyTimeOrNAString(value sql.NullInt64) string {
	if value.Valid {
		t := time.Unix(int64(value.Int64), 0)
		result, err := t.UTC().MarshalText()

		if err != nil {
			return "err"
		}

		return string(result)
	} else {
		return "n/a"
	}
}

// Convert an sql.NullInt64 into either a pretty datetime string, null
// via null.String
func PrettyTimeOrNullString(value sql.NullInt64) null.String {
	if value.Valid {
		t := time.Unix(int64(value.Int64), 0)
		result, err := t.UTC().MarshalText()

		if err != nil {
			return null.StringFrom("err")
		}

		return null.StringFrom(string(result))
	} else {
		return null.String{}
	}
}

func BoolOrNull(value sql.NullBool) null.Bool {
	if value.Valid {
		return null.BoolFrom(value.Bool)
	} else {
		return null.Bool{}
	}
}
