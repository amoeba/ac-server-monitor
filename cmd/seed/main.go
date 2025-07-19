package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"monitor/lib"
	"os"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"
	_ "github.com/mattn/go-sqlite3"
)

// TwoPartName generates realistic two-part names for servers
type TwoPartName struct {
	FirstPart  string
	SecondPart string
}

// GenerateTwoPartServerName creates realistic server names with proper capitalization
func GenerateTwoPartServerName() string {
	// Define lists of realistic server name components
	firstParts := []string{
		"Dragon", "Phoenix", "Thunder", "Shadow", "Crystal", "Storm", "Fire", "Ice",
		"Dark", "Light", "Steel", "Golden", "Silver", "Blood", "Iron", "Stone",
		"Wild", "Swift", "Fierce", "Noble", "Royal", "Sacred", "Ancient", "Mystic",
		"Crimson", "Azure", "Emerald", "Obsidian", "Titanium", "Platinum",
	}

	secondParts := []string{
		"Fortress", "Castle", "Stronghold", "Citadel", "Keep", "Tower", "Sanctuary", "Haven",
		"Arena", "Colosseum", "Stadium", "Battleground", "Fields", "Plains", "Valley", "Ridge",
		"Peak", "Summit", "Crater", "Cavern", "Dungeon", "Realm", "Domain", "Empire",
		"Kingdom", "Republic", "Federation", "Alliance", "Coalition", "Legion",
	}

	// Randomly select one from each list
	firstPart := firstParts[rand.Intn(len(firstParts))]
	secondPart := secondParts[rand.Intn(len(secondParts))]

	return fmt.Sprintf("%s %s", firstPart, secondPart)
}

// CreateFakeServer generates a fake server record
func CreateFakeServer() map[string]interface{} {
	// Generate a unique GUID
	guid := faker.UUIDDigit()
	
	// Generate server name
	name := GenerateTwoPartServerName()
	
	// Generate realistic server description
	description := fmt.Sprintf("A popular server hosting %s. Join for epic battles and adventures!", name)
	
	// Common emulators for AC servers
	emus := []string{"AC", "ACE", "GDLE", "ThwargLauncher"}
	emu := emus[rand.Intn(len(emus))]
	
	// Generate realistic IP addresses (using private IP ranges for safety)
	host := fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255))
	
	// Common AC server ports
	ports := []string{"9000", "9001", "9002", "9003", "9004", "9005", "9006", "9007", "9008", "9009"}
	port := ports[rand.Intn(len(ports))]
	
	// Server types
	types := []string{"PVP", "PVE", "RP", "Custom", "Retail"}
	serverType := types[rand.Intn(len(types))]
	
	// Generate website URL (optional)
	var websiteURL string
	if rand.Float32() < 0.7 { // 70% chance of having a website
		websiteURL = fmt.Sprintf("https://%s.com", strings.ToLower(strings.ReplaceAll(name, " ", "")))
	}
	
	// Generate Discord URL (optional)
	var discordURL string
	if rand.Float32() < 0.8 { // 80% chance of having Discord
		discordURL = "https://discord.gg/" + faker.Username()
	}
	
	now := time.Now().Unix()
	
	return map[string]interface{}{
		"guid":        guid,
		"name":        name,
		"description": description,
		"emu":         emu,
		"host":        host,
		"port":        port,
		"type":        serverType,
		"status":      "up", // Initial status
		"website_url": websiteURL,
		"discord_url": discordURL,
		"is_listed":   1, // All seeded servers are listed
		"created_at":  now,
		"updated_at":  now,
		"last_seen":   now,
		"is_online":   1, // Initially online
	}
}

// InsertServer inserts a server record into the database
func InsertServer(db *sql.DB, server map[string]interface{}) (int64, error) {
	query := `
		INSERT INTO servers (
			guid, name, description, emu, host, port, type, status, 
			website_url, discord_url, is_listed, created_at, updated_at, last_seen, is_online
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := db.Exec(query,
		server["guid"], server["name"], server["description"], server["emu"],
		server["host"], server["port"], server["type"], server["status"],
		server["website_url"], server["discord_url"], server["is_listed"],
		server["created_at"], server["updated_at"], server["last_seen"], server["is_online"],
	)
	
	if err != nil {
		return 0, err
	}
	
	return result.LastInsertId()
}

// CreateFakeUptimeData generates 2 weeks of fake uptime data for a server
func CreateFakeUptimeData(db *sql.DB, serverID int64, serverName string) error {
	log.Printf("Creating fake uptime data for server: %s (ID: %d)", serverName, serverID)
	
	// Generate data for the last 2 weeks (14 days)
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -14)
	
	// Check every 10 minutes (similar to the real cron job)
	interval := 10 * time.Minute
	
	// Server uptime characteristics (some servers are more reliable than others)
	uptimeReliability := 0.85 + rand.Float64()*0.14 // Between 85% and 99% uptime
	
	currentTime := startTime
	isCurrentlyUp := true // Start with server being up
	consecutiveDowns := 0
	
	for currentTime.Before(endTime) {
		var status int
		var rtt *int
		var message *string
		
		// Determine if server should be up or down
		if isCurrentlyUp {
			// If server is up, small chance of going down
			if rand.Float64() > uptimeReliability {
				isCurrentlyUp = false
				status = 0 // Down
				consecutiveDowns = 1
				
				// Random down reason
				downReasons := []string{
					"Connection refused",
					"Timeout",
					"Server maintenance",
					"Network error",
				}
				reason := downReasons[rand.Intn(len(downReasons))]
				message = &reason
			} else {
				status = 1 // Up
				// Generate realistic RTT (20-200ms)
				rttValue := 20 + rand.Intn(180)
				rtt = &rttValue
				upMessage := "Server online"
				message = &upMessage
			}
		} else {
			// If server is down, gradual chance of coming back up
			// Servers don't stay down forever
			consecutiveDowns++
			comeBackUpChance := float64(consecutiveDowns) * 0.1 // Increasing chance over time
			if comeBackUpChance > 0.8 {
				comeBackUpChance = 0.8 // Cap at 80%
			}
			
			if rand.Float64() < comeBackUpChance {
				isCurrentlyUp = true
				status = 1 // Up
				consecutiveDowns = 0
				rttValue := 20 + rand.Intn(180)
				rtt = &rttValue
				upMessage := "Server back online"
				message = &upMessage
			} else {
				status = 0 // Still down
				downReasons := []string{
					"Connection refused",
					"Timeout",
					"Server maintenance",
					"Network error",
					"Host unreachable",
				}
				reason := downReasons[rand.Intn(len(downReasons))]
				message = &reason
			}
		}
		
		// Insert status record
		query := `
			INSERT INTO statuses (server_id, created_at, status, rtt, message)
			VALUES (?, ?, ?, ?, ?)
		`
		
		_, err := db.Exec(query, serverID, currentTime.Unix(), status, rtt, message)
		if err != nil {
			return fmt.Errorf("failed to insert status for server %d: %w", serverID, err)
		}
		
		currentTime = currentTime.Add(interval)
	}
	
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting database seeding...")
	
	// Use the same DB path as the main application
	dbPath := lib.Env("DB_PATH", "./monitor.db")
	
	// Check if database already exists and warn user
	if _, err := os.Stat(dbPath); err == nil {
		log.Printf("Warning: Database %s already exists. This will add data to the existing database.", dbPath)
		log.Println("If you want a fresh database, please delete the existing file first.")
		
		// Give user a chance to cancel
		fmt.Print("Continue? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			log.Println("Seeding cancelled.")
			return
		}
	}
	
	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	// Run AutoMigrate to ensure database schema is up to date
	log.Println("Running database migrations...")
	if err := lib.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrations completed.")
	
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
	
	// Create 10 fake servers
	log.Println("Creating 10 fake servers...")
	var serverIDs []int64
	
	for i := 0; i < 10; i++ {
		server := CreateFakeServer()
		serverID, err := InsertServer(db, server)
		if err != nil {
			log.Fatalf("Failed to insert server %d: %v", i+1, err)
		}
		serverIDs = append(serverIDs, serverID)
		log.Printf("Created server: %s (ID: %d)", server["name"], serverID)
	}
	
	log.Println("All servers created successfully.")
	
	// Create 2 weeks of fake uptime data for each server
	log.Println("Generating 2 weeks of fake uptime data for each server...")
	
	for i, serverID := range serverIDs {
		// Get server name for logging
		var serverName string
		query := "SELECT name FROM servers WHERE id = ?"
		err := db.QueryRow(query, serverID).Scan(&serverName)
		if err != nil {
			log.Printf("Warning: Could not retrieve name for server ID %d: %v", serverID, err)
			serverName = fmt.Sprintf("Server-%d", i+1)
		}
		
		if err := CreateFakeUptimeData(db, serverID, serverName); err != nil {
			log.Fatalf("Failed to create uptime data for server %d: %v", serverID, err)
		}
	}
	
	log.Println("Database seeding completed successfully!")
	log.Printf("- Created 10 servers with realistic two-part names")
	log.Printf("- Generated 2 weeks of uptime data for each server")
	log.Printf("- Database saved to: %s", dbPath)
	log.Println("You can now run the main application to view the seeded data.")
}