package db

import (
	"database/sql"
	"log"
)

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

	log.Println("Done...")

	return nil
}
