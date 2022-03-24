package api

import (
	"database/sql"
	"log"
	"math"
)

type UptimeRow struct {
	Date    string
	Uptime  float64
	N       int
	RTTMin  sql.NullInt64
	RTTMax  sql.NullInt64
	RTTMean sql.NullFloat64
}

type RTT struct {
	Min  int `json:"min"`
	Max  int `json:"max"`
	Mean int `json:"mean"`
}

type UptimeApiItem struct {
	Date   string  `json:"date"`
	Uptime float64 `json:"uptime"`
	N      int     `json:"n"`
	RTT    RTT     `json:"rtt"`
}

type UptimeTemplateItem struct {
	Date      string
	Uptime    float64
	UptimeFmt string
	N         int
	RTTMin    string
	RTTMax    string
	RTTMean   string
}

var QUERY_UPTIME = `
	WITH ts(day, level)
	AS
	(
	SELECT date('now') AS day, 0 AS level
		UNION ALL
	SELECT date('now', '-' || level || ' day') AS day, level + 1 AS level FROM ts WHERE level < 14
	)
	SELECT
		day,
		COALESCE((sum(status) * 1.0 / COUNT(status)) * 100, 0) AS uptime,
		COUNT(status) AS n,
		MIN(rtt) as rtt_min,
		MAX(rtt) as rtt_max,
		AVG(rtt) as rtt_mean
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
func Uptime(db *sql.DB, server_id int) []UptimeApiItem {
	rows, err := db.Query(QUERY_UPTIME, server_id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var uptimes []UptimeApiItem

	for rows.Next() {
		var uptime UptimeRow
		var uptimeItem UptimeApiItem

		err := rows.Scan(
			&uptime.Date,
			&uptime.Uptime,
			&uptime.N,
			&uptime.RTTMin,
			&uptime.RTTMax,
			&uptime.RTTMean,
		)

		// Coerce to UptimeAPIItem
		uptimeItem.Date = uptime.Date
		uptimeItem.Uptime = uptime.Uptime
		uptimeItem.N = uptime.N
		uptimeItem.RTT.Min = int(uptime.RTTMin.Int64)
		uptimeItem.RTT.Max = int(uptime.RTTMax.Int64)
		uptimeItem.RTT.Mean = int(math.Round(uptime.RTTMean.Float64))

		if err != nil {
			log.Println(err)
		} else {
			uptimes = append(uptimes, uptimeItem)
		}
	}

	return uptimes
}
