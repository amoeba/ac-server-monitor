package api

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"gopkg.in/guregu/null.v4"
)

type ServerStatusRow struct {
	ID        int
	GUID      string
	Name      string
	Host      string
	Port      string
	Status    sql.NullBool
	IsListed  bool
	UpdatedAt int
	LastSeen  sql.NullInt64
}

type ServerAPIResponse struct {
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
	IsOnline    bool        `json:"online"`
	LastSeen    null.String `json:"last_seen"`
	LastChecked string      `json:"last_checked"`
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

func Servers(db *sql.DB) []ServerAPIResponse {
	stmt := `
	SELECT
		servers.id,
		servers.guid,
		servers.name,
		servers.host,
		servers.port,
		statuses.status,
		servers.is_listed,
		servers.updated_at,
		servers.last_seen
	FROM
		servers
	LEFT JOIN
		statuses ON servers.id = statuses.server_id
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
			&status.Status,
			&status.IsListed,
			&status.UpdatedAt,
			&status.LastSeen,
		)

		if err != nil {
			log.Println(err)
		} else {
			statuses = append(statuses, status)
		}
	}

	var finalResponse []ServerAPIResponse
	var item ServerAPIResponse

	for i := range statuses {
		lastSeenTime := PrettyTimeOrNullString(statuses[i].LastSeen)
		lastChecked := time.Unix(int64(statuses[i].UpdatedAt), 0)
		lastCheckedTime, err := lastChecked.UTC().MarshalText()

		if err != nil {
			log.Fatal(err)
		}

		item = ServerAPIResponse{
			ID:     statuses[i].ID,
			GUID:   statuses[i].GUID,
			Name:   statuses[i].Name,
			Active: statuses[i].IsListed,
			Address: ServerAPIResponseAddress{
				Host: statuses[i].Host,
				Port: statuses[i].Port,
			},
			Status: ServerAPIResponseStatus{
				IsOnline:    statuses[i].Status.Bool,
				LastSeen:    lastSeenTime,
				LastChecked: string(lastCheckedTime),
			},
		}

		finalResponse = append(finalResponse, item)
	}

	return finalResponse
}

func ServersWithUptimes(db *sql.DB) []ServerAPIResponseWithUptime {
	servers := Servers(db)

	var response []ServerAPIResponseWithUptime

	for i := range servers {
		var server ServerAPIResponseWithUptime

		server.ID = servers[i].ID
		server.GUID = servers[i].GUID
		server.Name = servers[i].Name
		server.Active = servers[i].Active
		server.Address = servers[i].Address
		server.Status = servers[i].Status

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
