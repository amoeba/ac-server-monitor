package api

import (
	"database/sql"
	"fmt"
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

type UptimeResult struct {
	Server  string          `json:"server"`
	Count   int             `json:"count"`
	Uptimes []UptimeApiItem `json:"uptimes"`
}

type UptimeApiItem struct {
	Date   string  `json:"date"`
	Uptime float64 `json:"uptime"`
	N      int     `json:"n"`
	RTT    RTT     `json:"rtt"`
}

type UptimeTemplateItem struct {
	Date        string
	Uptime      float64
	UptimeFmt   string
	UptimeClass string
	N           int
	RTTMin      string
	RTTMax      string
	RTTMean     string
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

var QUERY_UPTIME_3_MONTHS = `
	WITH RECURSIVE week_grid(week_start, week_num, day_offset, current_day) AS (
		-- Find the Monday that's roughly 5 months ago
		SELECT
			date('now', '-150 days', 'weekday 1') as week_start,
			0 as week_num,
			0 as day_offset,
			date('now', '-150 days', 'weekday 1') as current_day

		UNION ALL

		-- Generate days until we reach today
		SELECT
			CASE WHEN day_offset = 6
				THEN date(week_start, '+7 days')
				ELSE week_start
			END,
			CASE WHEN day_offset = 6
				THEN week_num + 1
				ELSE week_num
			END,
			(day_offset + 1) % 7,
			date(current_day, '+1 day')
		FROM week_grid
		WHERE current_day < date('now')
	),
	calendar_days AS (
		SELECT
			current_day as day,
			week_num,
			day_offset
		FROM week_grid
	)
	SELECT
		day,
		COALESCE((sum(status) * 1.0 / COUNT(status)) * 100, 0) AS uptime,
		COUNT(status) AS n,
		MIN(rtt) as rtt_min,
		MAX(rtt) as rtt_max,
		AVG(rtt) as rtt_mean
	FROM calendar_days
	LEFT JOIN statuses ON
		date(statuses.created_at, 'unixepoch') = calendar_days.day
		AND statuses.server_id = ?
	GROUP BY day, week_num, day_offset
	ORDER BY week_num, day_offset;
`

const (
	UPTIME_CLASS_HIGH string = "high"
	UPTIME_CLASS_MID  string = "mid"
	UPTIME_CLASS_LOW  string = "low"
)

func GetUptimeClass(uptime float64) string {
	if uptime >= 99 {
		return UPTIME_CLASS_HIGH
	} else if uptime < 99 && uptime >= 50 {
		return UPTIME_CLASS_MID
	} else {
		return UPTIME_CLASS_LOW
	}
}

// TODO: Handle the param properly
// TODO: Handle not found
func Uptime(db *sql.DB, server_id int, name string) UptimeResult {
	rows, err := db.Query(QUERY_UPTIME, server_id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var result UptimeResult
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

	result.Server = name
	result.Count = len(uptimes)
	result.Uptimes = uptimes

	return result
}

func UptimeThreeMonths(db *sql.DB, server_id int, name string) []UptimeTemplateItem {
	rows, err := db.Query(QUERY_UPTIME_3_MONTHS, server_id)

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

	// Data is already in chronological order from the query

	return uptimes
}
