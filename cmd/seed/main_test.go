package main

import (
	"database/sql"
	"monitor/lib"
	"os"
	"strings"
	"testing"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestSeedProgram(t *testing.T) {
	// Create a temporary database file
	tempDB := "./test_seed.db"
	defer os.Remove(tempDB)

	// Set environment variable for test database
	os.Setenv("DB_PATH", tempDB)
	defer os.Unsetenv("DB_PATH")

	// Open database connection
	db, err := sql.Open("sqlite3", tempDB)
	assert.NoError(t, err)
	defer db.Close()

	// Run AutoMigrate
	err = lib.AutoMigrate(db)
	assert.NoError(t, err)

	// Test fake server generation
	server := CreateFakeServer()
	assert.NotEmpty(t, server["name"])
	assert.NotEmpty(t, server["guid"])
	assert.NotEmpty(t, server["host"])
	assert.NotEmpty(t, server["port"])

	// Test server name format (should contain a space for two-part name)
	name, ok := server["name"].(string)
	assert.True(t, ok)
	assert.Contains(t, name, " ", "Server name should be two-part with a space")

	// Test server insertion
	serverID, err := InsertServer(db, server)
	assert.NoError(t, err)
	assert.Greater(t, serverID, int64(0))

	// Verify server was inserted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM servers").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Test uptime data creation
	err = CreateFakeUptimeData(db, serverID, name)
	assert.NoError(t, err)

	// Verify uptime data was created (should have many entries for 2 weeks)
	var statusCount int
	err = db.QueryRow("SELECT COUNT(*) FROM statuses WHERE server_id = ?", serverID).Scan(&statusCount)
	assert.NoError(t, err)
	assert.Greater(t, statusCount, 1000, "Should have many status entries for 2 weeks of data")
}

func TestTwoPartServerNames(t *testing.T) {
	// Test multiple generated names to ensure they follow the two-part format
	for range 10 {
		name := GenerateTwoPartServerName()
		assert.Contains(t, name, " ", "Server name should contain a space for two-part format")

		// Check that both parts are capitalized
		parts := strings.Split(name, " ")
		assert.Len(t, parts, 2, "Server name should have exactly two parts")

		// First letter of each part should be uppercase
		assert.True(t, unicode.IsUpper(rune(parts[0][0])), "First part should be capitalized")
		assert.True(t, unicode.IsUpper(rune(parts[1][0])), "Second part should be capitalized")
	}
}
