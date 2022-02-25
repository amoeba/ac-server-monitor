package api

import (
	"database/sql"
	"log"
	"time"
)

type StatusesRow struct {
	GUID      string
	Name      string
	Host      string
	Port      string
	Status    bool
	IsListed  bool
	UpdatedAt int
	LastSeen  int
	FirstSeen int
	Count     int
}

type ServerAPIResponse struct {
	ID      string                   `json:"id"`
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
	IsOnline    bool   `json:"online"`
	FirstSeen   string `json:"first_seen"`
	LastSeen    string `json:"last_seen"`
	LastChecked string `json:"last_checked"`
}

func Servers(db *sql.DB) []ServerAPIResponse {
	log.Println("API.Servers")

	stmt := `
	SELECT
		servers.guid,
		servers.name,
		servers.host,
		servers.port,
		statuses.status,
		servers.is_listed,
		servers.updated_at,
		MAX(statuses.created_at) AS last_seen,
		MIN(statuses.created_at) AS first_seen,
		COUNT(statuses.created_at) as count
	FROM
		servers
	LEFT JOIN
		statuses ON servers.id = statuses.server_id
	GROUP BY servers.id
	ORDER BY lower(servers.name);
	`

	rows, err := db.Query(stmt)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var statuses []StatusesRow
	var status StatusesRow

	for rows.Next() {
		err := rows.Scan(
			&status.GUID,
			&status.Name,
			&status.Host,
			&status.Port,
			&status.Status,
			&status.IsListed,
			&status.UpdatedAt,
			&status.LastSeen,
			&status.FirstSeen,
			&status.Count,
		)

		if err != nil {
			log.Println(err)
		} else {
			statuses = append(statuses, status)
		}
	}

	// Add in relative times
	var finalResponse []ServerAPIResponse

	var item ServerAPIResponse

	var firstSeen time.Time
	var lastSeen time.Time
	var lastChecked time.Time

	for i := range statuses {
		firstSeen = time.Unix(int64(statuses[i].FirstSeen), 0)
		firstSeenTime, err := firstSeen.UTC().MarshalText()

		if err != nil {
			log.Fatal(err)
		}

		lastSeen = time.Unix(int64(statuses[i].LastSeen), 0)
		lastSeenTime, err := lastSeen.UTC().MarshalText()

		if err != nil {
			log.Fatal(err)
		}

		lastChecked = time.Unix(int64(statuses[i].UpdatedAt), 0)
		lastCheckedTime, err := lastChecked.UTC().MarshalText()

		if err != nil {
			log.Fatal(err)
		}

		item = ServerAPIResponse{
			ID:     statuses[i].GUID,
			Name:   statuses[i].Name,
			Active: statuses[i].IsListed,
			Address: ServerAPIResponseAddress{
				Host: statuses[i].Host,
				Port: statuses[i].Port,
			},
			Status: ServerAPIResponseStatus{
				IsOnline:    statuses[i].Status,
				FirstSeen:   string(firstSeenTime),
				LastSeen:    string(lastSeenTime),
				LastChecked: string(lastCheckedTime),
			},
		}

		finalResponse = append(finalResponse, item)
	}

	return finalResponse
}

type UptimeRow struct {
	Date  string  `json:"date"`
	Ratio float64 `json:"uptime"`
}

// TODO: Handle the param properly
// TODO: Handle not found
func Uptime(db *sql.DB, server_id int) []UptimeRow {
	log.Println("API.Uptime")

	stmt := `
		SELECT
			DATE(created_at, "unixepoch") AS created_datetime,
			(SUM(status) * 1.0 / COUNT(status)) AS ratio
		FROM statuses
		WHERE
			statuses.server_id = ?
		GROUP BY
			created_datetime
		ORDER BY
			created_datetime ASC;
	`

	rows, err := db.Query(stmt, server_id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var uptimes []UptimeRow
	var uptime UptimeRow

	for rows.Next() {
		err := rows.Scan(
			&uptime.Date,
			&uptime.Ratio,
		)

		if err != nil {
			log.Println(err)
		} else {
			uptimes = append(uptimes, uptime)
		}
	}

	return uptimes
}
