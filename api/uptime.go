package api

import (
	"database/sql"
	"log"
)

type UptimeRow struct {
	Date   string  `json:"date"`
	Uptime float64 `json:"uptime"`
	N      int     `json:"n"`
}

type UptimeTemplateItem struct {
	Date      string  `json:"date"`
	Uptime    float64 `json:"uptime"`
	UptimeFmt string  `json:"-"`
	N         int     `json:"n"`
}

var QUERY_UPTIME = `
	WITH ts(day, level)
	AS
	(
	SELECT  date('now') AS day, 0 AS level
		UNION ALL
	SELECT date('now', '-' || level || ' day') AS day, level + 1 AS level FROM ts WHERE level < 14
	)
	SELECT
		day,
		COALESCE((sum(status) * 1.0 / COUNT(status)) * 100, 0) AS uptime,
		COUNT(status) AS n
	FROM ts
	LEFT JOIN statuses
	ON
		date(statuses.created_at, 'unixepoch') = ts.day
	AND
		statuses.server_id = ?
	GROUP BY day;
`

// TODO: Handle the param properly
// TODO: Handle not found
func Uptime(db *sql.DB, server_id int) []UptimeRow {
	log.Println("API.Uptime")

	rows, err := db.Query(QUERY_UPTIME, server_id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var uptimes []UptimeRow
	var uptime UptimeRow

	for rows.Next() {
		err := rows.Scan(
			&uptime.Date,
			&uptime.Uptime,
			&uptime.N,
		)

		if err != nil {
			log.Println(err)
		} else {
			uptimes = append(uptimes, uptime)
		}
	}

	return uptimes
}
