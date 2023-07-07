package lib

import (
	"database/sql"
	"fmt"
	"log"
	"monitor/api"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func DropOrFail(t *testing.T, db *sql.DB, table_name string) {
	_, err := db.Query("DROP TABLE IF EXISTS '?'", table_name)

	if err != nil {
		t.Fatalf("Failed to drop servers table: %s", err)
	}
}

func GenerateTestServerList() ServerList {
	fmt.Println("GenerateTestServerList")

	list := ServerList{}

	list.Servers = append(
		list.Servers,
		ServerListItem{ID: "UpServer", Name: "UpServer"},
		ServerListItem{ID: "DownServer", Name: "DownServer"},
	)

	fmt.Println(len(list.Servers))
	return list
}

func AssertNRows(t *testing.T, db *sql.DB, table_name string, count int) {
	fmt.Printf("Asserting '%s' table is of size %d\n", table_name, count)

	rows, err := db.Query(fmt.Sprintf(`
		SELECT
			count(*)
		FROM '%s'
		`, table_name))

	if err != nil {
		t.Fatal(err)
	}

	defer rows.Close()

	var n int

	for rows.Next() {
		rows.Scan(&n)
	}

	assert.Equal(t, count, n)
}

func SetLastSeen(t *testing.T, db *sql.DB, server_name string, last_seen int64) {
	query := `
		UPDATE servers
		SET last_seen = ?
		WHERE name = ?
	`

	_, err := db.Exec(query, last_seen, server_name)

	if err != nil {
		t.Fatal(err)
	}
}

func TestLastSeen(t *testing.T) {
	// TODO: Factor this out
	db, err := sql.Open("sqlite3", "../monitor_test.db")

	if err != nil {
		log.Fatal(err)
	}

	// Empty out and set up the database
	DropOrFail(t, db, "servers")
	AutoMigrate(db)

	// Populate the database with a mocked Check
	list := GenerateTestServerList()
	UpdateServersTable(db, list)

	// Check number of rows we inserted
	AssertNRows(t, db, "servers", 2)

	// Simulate both servers being up at some point
	// and then DownServer being down at some future point
	now := time.Now().UTC().Unix()
	future := now + 3600
	SetLastSeen(t, db, "UpServer", now)
	SetLastSeen(t, db, "DownServer", now)
	SetLastSeen(t, db, "UpServer", future)

	// Verify API response result
	response := api.Servers(db)
	assert.Equal(t, response[0].Status.LastSeen, api.PrettyTimeOrNullString(sql.NullInt64{now, true}))
	assert.Equal(t, response[1].Status.LastSeen, api.PrettyTimeOrNullString(sql.NullInt64{future, true}))
}
