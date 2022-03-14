package api

import (
	"database/sql"
	"log"
	"time"
)

type StatusesRow struct {
	ID        int
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
	IsOnline    bool   `json:"online"`
	FirstSeen   string `json:"first_seen"`
	LastSeen    string `json:"last_seen"`
	LastChecked string `json:"last_checked"`
}

type ServerAPIResponseWithUptime struct {
	ID      int                      `json:"id"`
	GUID    string                   `json:"guid"`
	Name    string                   `json:"name"`
	Active  bool                     `json:"active"`
	Address ServerAPIResponseAddress `json:"address"`
	Status  ServerAPIResponseStatus  `json:"status"`
	Uptime  []UptimeRow              `json:"uptime"`
}

func Servers(db *sql.DB) []ServerAPIResponse {
	log.Println("API.Servers")

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
			&status.ID,
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
			ID:     statuses[i].ID,
			GUID:   statuses[i].GUID,
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
		log.Printf("Querying uptime for %s", server.ID)

		rows, err := db.Query(QUERY_UPTIME, server.ID)

		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		var uptimes []UptimeRow

		for rows.Next() {
			var uptime UptimeRow

			err := rows.Scan(
				&uptime.Date,
				&uptime.Uptime,
			)

			if err != nil {
				log.Fatal(err)
			}

			uptimes = append(uptimes, uptime)
		}

		server.Uptime = uptimes
		response = append(response, server)
	}

	return response
}
