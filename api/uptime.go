package api

import (
	"database/sql"
	"log"
)

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
