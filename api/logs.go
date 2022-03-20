package api

import (
	"database/sql"
	"log"
	"time"
)

type LogRow struct {
	Message   string
	CreatedAt int
}

type LogApiItem struct {
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

var QUERY_LOGS = `
SELECT message, created_at
FROM logs
ORDER BY created_at DESC;
`

func Logs(db *sql.DB) []LogApiItem {
	rows, err := db.Query(QUERY_LOGS)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var logs []LogApiItem

	for rows.Next() {
		var row LogRow
		var item LogApiItem

		err := rows.Scan(
			&row.Message,
			&row.CreatedAt,
		)

		if err != nil {
			log.Fatal(err)
		}

		// Song and dance to convert the database row to JSON
		createdAt := time.Unix(int64(row.CreatedAt), 0)
		createdAtTime, err := createdAt.UTC().MarshalText()

		item.Message = row.Message
		item.CreatedAt = string(createdAtTime)

		if err != nil {
			log.Println(err)
		} else {
			logs = append(logs, item)
		}
	}

	return logs
}
