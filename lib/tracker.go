package lib

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func CreateServerRecord(tx *sql.Tx, s *ServerListItem) error {
	log.Printf("CreateServerRecord %s", s.Name)

	now := time.Now().Unix()

	queryString := `
		INSERT INTO servers ( guid, name, description, emu, host, port, type, status, website_url, discord_url, is_listed, created_at, updated_at )
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?);`

	createdResult, err := tx.Exec(
		queryString,
		s.ID,
		s.Name,
		s.Description,
		s.Emu,
		s.Host,
		s.Port,
		s.Type,
		s.Status,
		s.Website,
		s.Discord,
		now,
		now,
	)

	if err != nil {
		log.Fatal(err)
	}

	createRows, createRowsErr := createdResult.RowsAffected()

	if createRowsErr != nil {
		log.Fatal(err)
	}

	log.Printf("Created %d row(s)", createRows)

	return err
}

func UpdateServerRecord(tx *sql.Tx, s *ServerListItem) error {
	now := time.Now().Unix()

	queryString := `
		UPDATE servers
		SET
			 guid = ?,
			 name = ?,
			 description = ?,
			 emu = ?,
			 host = ?,
			 port = ?,
			 type = ?,
			 status = ?,
			 website_url = ?,
			 discord_url = ?,
			 is_listed = 1,
			 updated_at = ?
		WHERE guid = ?;
	`

	_, err := tx.Exec(
		queryString,
		s.ID,
		s.Name,
		s.Description,
		s.Emu,
		s.Host,
		s.Port,
		s.Type,
		s.Status,
		s.Website,
		s.Discord,
		now,
		s.ID,
	)

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func CreateOrUpdateServer(tx *sql.Tx, s *ServerListItem) error {
	// Find
	res, err := tx.Query(`
		SELECT id as count
		FROM servers
		WHERE guid = ?
	`, s.ID)

	if err != nil {
		log.Fatal(err)
	}

	was_found := res.Next()

	res.Close()

	if was_found {
		err = UpdateServerRecord(tx, s)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = CreateServerRecord(tx, s)

		if err != nil {
			log.Fatal(err)
		}
	}

	return err
}

func UpdateStatusForServer(db *sql.DB, s *ServerListItem) error {
	now := time.Now().Unix()

	// Get the server's ID
	res, err := db.Query(`
		SELECT id
		FROM servers
		WHERE guid = ?
		LIMIT 1
	`, s.ID)

	if err != nil {
		log.Fatal(err)
	}

	var id int

	for res.Next() {
		err := res.Scan(&id)

		if err != nil {
			log.Fatal(err)
		}
	}

	if id <= 0 {
		log.Fatalf("Failed to find server for status with guid %s", s.ID)
	}

	// Add a new row
	tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	defer tx.Commit()

	// Actually check server
	server := Server{
		Host: s.Host,
		Port: s.Port,
	}

	rtt_start := time.Now().UnixMilli()
	up, err := Check(server)
	rtt := time.Now().UnixMilli() - rtt_start

	if err != nil {
		up = false
		message := fmt.Sprintf("Check for server %s failed with error message `%s`.", s.Name, err)
		log.Print(message)
		WriteLog(db, message)
	} else {
		log.Printf("Check for server %s succeeded in %d ms", s.Name, rtt)
	}

	query := `
	INSERT INTO statuses (server_id, created_at, status, rtt, message)
	VALUES (?, ?, ?, ?, ?)
	`

	var message string

	if err != nil {
		message = err.Error()
	}

	_, txErr := tx.Exec(query, id, now, up, rtt, message)

	if txErr != nil {
		log.Fatal(txErr)
	}

	return nil
}

func Update(db *sql.DB) error {
	log.Print("Beginning update...")

	// Fetch latest list
	lst, err := Fetch()

	if err != nil {
		log.Fatalf("Error fetching server list in update: %s", err)
	}

	// open a tx
	tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	// Set each item in the list to not-in-list
	_, err = tx.Exec(`
		UPDATE servers
		SET is_listed = 0
  `)

	if err != nil {
		log.Fatal(err)
	}

	// Go through each item in the list and
	for i := range lst.Servers {
		upsertErr := CreateOrUpdateServer(tx, &lst.Servers[i])

		if upsertErr != nil {
			log.Fatal(upsertErr)
		}
	}

	// Force a commit now before moving on to updating statuses
	tx.Commit()

	// Get statuses for each server in the list
	for i := range lst.Servers {
		updateStatusError := UpdateStatusForServer(db, &lst.Servers[i])

		if updateStatusError != nil {
			log.Fatal(updateStatusError)
		}
	}

	log.Print("Done with update.")

	return nil
}
