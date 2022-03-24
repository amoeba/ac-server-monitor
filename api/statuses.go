package api

import (
	"database/sql"
	"log"
	"time"
)

var QUERY_SERVER_BY_ID = `
SELECT name
FROM servers
WHERE id = ?
LIMIT 1
`

var QUERY_STATUSES = `
SELECT status, created_at, rtt, message
FROM statuses
WHERE server_id = ?
ORDER BY created_at DESC
LIMIT 100;
`

type StatusApiResponse struct {
	ServerName string `json:"server"`
	Statuses   []StatusApiStatusItem
}

type StatusApiStatusItem struct {
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	RTT       int    `json:"rtt"`
	Message   string `json:"message"`
}

type StatusesRow struct {
	Status    int
	CreatedAt int
	RTT       sql.NullInt64
	Message   sql.NullString
}

func GetServerNameById(db *sql.DB, id int) (string, error) {
	result, err := db.Query(QUERY_SERVER_BY_ID, id)

	if err != nil {
		return "", err
	}

	defer result.Close()

	var name string

	for result.Next() {
		scanErr := result.Scan(&name)

		if scanErr != nil {
			return "", scanErr
		}
	}

	return name, nil
}

func Statuses(db *sql.DB, server_id int) StatusApiResponse {
	var response StatusApiResponse

	// Find the server's name by ID first
	server_name, getErr := GetServerNameById(db, server_id)

	if getErr != nil {
		return response
	}

	response.ServerName = server_name

	// Then grab the statuses
	rows, err := db.Query(QUERY_STATUSES, server_id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var statuses []StatusApiStatusItem

	for rows.Next() {
		var status StatusesRow
		var statusItem StatusApiStatusItem

		err := rows.Scan(
			&status.Status,
			&status.CreatedAt,
			&status.RTT,
			&status.Message,
		)

		if err != nil {
			log.Fatal(err)
		}

		// Song and dance to convert the database row to JSON
		createdAt := time.Unix(int64(status.CreatedAt), 0)
		createdAtTime, err := createdAt.UTC().MarshalText()

		statusItem.Message = status.Message.String

		if status.Status == 1 {
			statusItem.Status = "UP"
		} else {
			statusItem.Status = "DOWN"
		}

		statusItem.CreatedAt = string(createdAtTime)
		statusItem.RTT = int(status.RTT.Int64)
		statusItem.Message = string(status.Message.String)

		if err != nil {
			log.Println(err)
		} else {
			statuses = append(statuses, statusItem)
		}
	}

	response.Statuses = statuses

	return response
}
