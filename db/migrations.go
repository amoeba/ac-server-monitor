package db

import (
	"database/sql"
	"log"
)

func CreateServersTable(db *sql.DB) (sql.Result, error) {
	log.Println("CreateServersTable")

	createTableStatement := `
	CREATE TABLE IF NOT EXISTS servers (
		id INTEGER NOT NULL PRIMARY KEY,
		guid TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		emu TEXT NOT NULL,
		host TEXT NOT NULL,
		port TEXT NOT NULL,
		type TEXT NOT NULL,
		status TEXT,
		website_url TEXT,
		discord_url TEXT,
		is_listed INTEGER NOT NULL,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	);

	CREATE INDEX IF NOT EXISTS servers_server_id ON servers (id);
	CREATE INDEX IF NOT EXISTS servers_server_name ON servers (name);
	`

	return db.Exec(createTableStatement)
}

func CreateStatusesTable(db *sql.DB) (sql.Result, error) {
	log.Println("CreateStatusesTable")

	createTableStatement := `
	CREATE TABLE IF NOT EXISTS statuses (
		id INTEGER not null primary key NOT NULL,
		server_id INTEGER NOT NULL,
		created_at INTEGER NOT NULL,
		status INTEGER NOT NULL
	);

	CREATE INDEX IF NOT EXISTS statuses_server_id ON statuses (server_id);
	`

	return db.Exec(createTableStatement)
}
