package lib

import (
	"database/sql"
	"log"
)

func QueryLastUpdated(db *sql.DB) string {
	query := `
	SELECT updated_at
	FROM servers
	ORDER BY updated_at desc
	LIMIT 1
	`

	res, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
	}

	var updated_at int64

	for res.Next() {
		res.Scan(&updated_at)
	}

	return RelativeTime(int64(updated_at))
}

func QueryTotalNumStatuses(db *sql.DB) int64 {
	query := `
	SELECT MAX(ROWID) as count
	FROM statuses
	LIMIT 1
	`

	res, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
	}

	var count int64

	for res.Next() {
		res.Scan(&count)
	}

	return count
}
