package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"business-mgmt-backend/pkg/database"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	cfg := database.LoadConfigFromEnv()
	fmt.Println("==================================================================")
	fmt.Println("🐘 Business Management Software — Database Migration Runner")
	fmt.Printf("🎯 Target Database: %s:%s/%s\n", cfg.Host, cfg.Port, cfg.DBName)
	fmt.Println("==================================================================")

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Could not connect to database: %v", err)
	}
	defer db.Close()

	switch command {
	case "up":
		fmt.Println("🚀 Executing pending migrations...")
		if err := database.MigrateUp(db); err != nil {
			log.Fatalf("❌ Migration up failed: %v", err)
		}

	case "down":
		fmt.Println("⏪ Rolling back last migration...")
		if err := database.MigrateDown(db); err != nil {
			log.Fatalf("❌ Migration down failed: %v", err)
		}

	case "reset":
		fmt.Println("⚠️ Resetting all migrations and re-applying...")
		if err := database.MigrateDownAll(db); err != nil {
			log.Fatalf("❌ Rollback all failed: %v", err)
		}
		if err := database.MigrateUp(db); err != nil {
			log.Fatalf("❌ Re-migration failed: %v", err)
		}

	case "status":
		fmt.Println("📊 Migration Status:")
		statuses, err := database.Status(db)
		if err != nil {
			log.Fatalf("❌ Could not fetch status: %v", err)
		}
		fmt.Printf("%-35s %-12s %s\n", "MIGRATION", "STATUS", "APPLIED AT")
		fmt.Println("------------------------------------------------------------------")
		for _, s := range statuses {
			state := "⏳ Pending"
			appliedTime := "-"
			if s.IsApplied {
				state = "✅ Applied"
				if s.AppliedAt != nil {
					appliedTime = s.AppliedAt.Local().Format(time.RFC3339)
				}
			}
			fmt.Printf("%-35s %-12s %s\n", s.Name, state, appliedTime)
		}

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	fmt.Println("==================================================================")
	fmt.Println("✨ Done!")
}

func printUsage() {
	fmt.Println("Usage: go run ./cmd/migrate [command]")
	fmt.Println("")
	fmt.Println("Available Commands:")
	fmt.Println("  up      Apply all pending migrations (schema + seed)")
	fmt.Println("  down    Roll back the latest applied migration")
	fmt.Println("  reset   Roll back all migrations and re-apply from scratch")
	fmt.Println("  status  Check current status of all migrations")
	fmt.Println("")
	fmt.Println("Environment Variables:")
	fmt.Println("  DATABASE_URL  Full Postgres connection string (optional)")
	fmt.Println("  DB_HOST       Postgres host (default: localhost)")
	fmt.Println("  DB_PORT       Postgres port (default: 5432)")
	fmt.Println("  DB_USER       Postgres user (default: postgres)")
	fmt.Println("  DB_PASSWORD   Postgres password (default: postgres)")
	fmt.Println("  DB_NAME       Database name (default: business_mgmt)")
	fmt.Println("  DB_SSLMODE    SSL mode (default: disable)")
}
