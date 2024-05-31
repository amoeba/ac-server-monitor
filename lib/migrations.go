package lib

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

func CreateLogsTable(db *sql.DB) (sql.Result, error) {
	log.Println("CreateLogsTable")

	createTableStatement := `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER NOT NULL PRIMARY KEY,
		message TEXT NOT NULL,
		created_at INTEGER NOT NULL
	);
	`

	return db.Exec(createTableStatement)
}

func DropLogsTable(db *sql.DB) (sql.Result, error) {
	log.Println("DropLogsTable")

	createTableStatement := `
	DROP TABLE  logs;
	`

	return db.Exec(createTableStatement)
}

func AlterStatusesAddRTTAndMessage(db *sql.DB) (sql.Result, error) {
	log.Println("AlterStatusesAddRTTAndMessage")

	checkIfColExistsStatement := `
	SELECT message
	FROM statuses
	LIMIT 1;
	`

	checkRes, checkErr := db.Exec(checkIfColExistsStatement)

	if checkErr == nil {
		log.Print("Skipping migration")

		return checkRes, checkErr
	}

	alterTableStatement := `
	ALTER TABLE statuses
	ADD COLUMN rtt INTEGER;
	ALTER TABLE statuses
	ADD COLUMN message TEXT;
	`

	log.Printf("Running %s", alterTableStatement)

	alterRes, alterErr := db.Exec(alterTableStatement)

	if alterErr != nil {
		log.Fatal(alterErr)
	}

	return alterRes, alterErr
}

func CreateStatusesCreatedAtIndex(db *sql.DB) (sql.Result, error) {
	log.Println("CreateStatusesCreatedAtIndex")

	createIndexStatement := `
	CREATE INDEX IF NOT EXISTS statuses_date ON statuses (date(created_at, 'unixepoch'));
	`

	return db.Exec(createIndexStatement)
}

func AlterServersAddLastSeen(db *sql.DB) (sql.Result, error) {
	log.Println("AlterServersAddLastSeen")

	createIndexStatement := `
	ALTER TABLE servers ADD last_seen INTEGER;
	`

	return db.Exec(createIndexStatement)
}

func AlterServersAddIsOnline(db *sql.DB) (sql.Result, error) {
	log.Println("AlterServersAddIsOnline")

	createIndexStatement := `
	ALTER TABLE servers ADD is_online INTEGER;
	`

	return db.Exec(createIndexStatement)
}

func UpdateStatusesFixDownWithNullMessage(db *sql.DB) (sql.Result, error) {
	// Fixes data issue partially addressed by
	// https://github.com/amoeba/ac-server-monitor/pull/14 and
	// commit 7601422fc7fbd157f0f0b0194f55bf9db6cbbb7b.
	log.Println("UpdateStatusesFixDownWithNullMessage")

	statement := `
		UPDATE statuses
		SET status = 1
		WHERE status = 0 AND message IS NULL
	`

	return db.Exec(statement)
}

func UpdateStatusesFixResponseSizeFourtyFour(db *sql.DB) (sql.Result, error) {
	// Fixes data issue partially addressed by
	// https://github.com/amoeba/ac-server-monitor/pull/14 and
	// commit 7601422fc7fbd157f0f0b0194f55bf9db6cbbb7b.
	log.Println("UpdateStatusesFixResponseSizeFourtyFour")

	statement := `
		UPDATE statuses
		SET status = 1
		WHERE status = 0 AND message LIKE '%bytes read was 44%';
	`

	return db.Exec(statement)
}

func AutoMigrate(db *sql.DB) error {
	log.Println("AutoMigrating...")

	var err error

	_, err = CreateServersTable(db)

	if err != nil {
		return err
	}

	_, err = CreateStatusesTable(db)

	if err != nil {
		return err
	}

	_, err = CreateLogsTable(db)

	if err != nil {
		return err
	}

	_, err = AlterStatusesAddRTTAndMessage(db)

	if err != nil {
		return err
	}

	_, err = CreateStatusesCreatedAtIndex(db)

	if err != nil {
		return err
	}

	// Ignore errors here since SQLite lets this pass silently
	_, err = AlterServersAddLastSeen(db)
	_, err = AlterServersAddIsOnline(db)

	_, err = DropLogsTable(db)

	if err != nil {
		return err
	}

	_, err = UpdateStatusesFixDownWithNullMessage(db)

	if err != nil {
		return err
	}

	_, err = UpdateStatusesFixResponseSizeFourtyFour(db)

	if err != nil {
		return err
	}

	log.Println("...AutoMigration Done")

	return nil
}
