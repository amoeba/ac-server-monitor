package api

import (
	"database/sql"
	"log"
)

type LogRow struct {
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

var QUERY_LOGS = `
SELECT message, created_at
FROM logs
ORDER BY created_at DESC;
`

func Logs(db *sql.DB) []LogRow {
	log.Println("API.Logs")

	rows, err := db.Query(QUERY_LOGS)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var logs []LogRow

	for rows.Next() {
		var row LogRow

		err := rows.Scan(
			&row.Message,
			&row.CreatedAt,
		)

		if err != nil {
			log.Println(err)
		} else {
			logs = append(logs, row)
		}
	}

	return logs
}
