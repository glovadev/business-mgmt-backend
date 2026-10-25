package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"
	"time"

	"business-mgmt-backend/migrations"
)

type MigrationRecord struct {
	Version   string
	AppliedAt time.Time
}

type MigrationStatus struct {
	Name      string
	IsApplied bool
	AppliedAt *time.Time
}

// EnsureMigrationTable initializes the tracking table if not present.
func EnsureMigrationTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to ensure schema_migrations table: %w", err)
	}
	return nil
}

// GetAppliedMigrations returns a map of already applied migration versions.
func GetAppliedMigrations(db *sql.DB) (map[string]time.Time, error) {
	if err := EnsureMigrationTable(db); err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT version, applied_at FROM schema_migrations ORDER BY version ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]time.Time)
	for rows.Next() {
		var version string
		var appliedAt time.Time
		if err := rows.Scan(&version, &appliedAt); err != nil {
			return nil, err
		}
		applied[version] = appliedAt
	}
	return applied, nil
}

// MigrateUp applies all pending *.up.sql migrations in ascending order.
func MigrateUp(db *sql.DB) error {
	if err := EnsureMigrationTable(db); err != nil {
		return err
	}

	applied, err := GetAppliedMigrations(db)
	if err != nil {
		return err
	}

	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to read embedded migrations: %w", err)
	}

	var upFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			upFiles = append(upFiles, entry.Name())
		}
	}
	sort.Strings(upFiles)

	var appliedCount int
	for _, filename := range upFiles {
		version := strings.TrimSuffix(filename, ".up.sql")
		if _, ok := applied[version]; ok {
			continue // Already applied
		}

		log.Printf("⏳ Applying migration: %s ...", filename)
		content, err := fs.ReadFile(migrations.FS, filename)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed executing migration %s: %w", filename, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version, applied_at) VALUES ($1, NOW())", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed recording migration %s: %w", filename, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed committing migration %s: %w", filename, err)
		}

		log.Printf("✅ Successfully applied: %s", filename)
		appliedCount++
	}

	if appliedCount == 0 {
		log.Println("✨ Database schema is up to date (no pending migrations).")
	} else {
		log.Printf("🎉 Successfully applied %d migration(s).", appliedCount)
	}
	return nil
}

// MigrateDown rolls back the single most recently applied migration.
func MigrateDown(db *sql.DB) error {
	if err := EnsureMigrationTable(db); err != nil {
		return err
	}

	applied, err := GetAppliedMigrations(db)
	if err != nil {
		return err
	}

	if len(applied) == 0 {
		log.Println("ℹ️ No applied migrations to roll back.")
		return nil
	}

	// Find the latest applied version
	var versions []string
	for v := range applied {
		versions = append(versions, v)
	}
	sort.Strings(versions)
	latestVersion := versions[len(versions)-1]
	downFilename := latestVersion + ".down.sql"

	log.Printf("⏳ Rolling back migration: %s ...", downFilename)
	content, err := fs.ReadFile(migrations.FS, downFilename)
	if err != nil {
		return fmt.Errorf("failed to read rollback file %s: %w", downFilename, err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if _, err := tx.Exec(string(content)); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed executing rollback %s: %w", downFilename, err)
	}

	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", latestVersion); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed removing migration record %s: %w", latestVersion, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed committing rollback %s: %w", downFilename, err)
	}

	log.Printf("✅ Successfully rolled back: %s", latestVersion)
	return nil
}

// MigrateDownAll rolls back all applied migrations in reverse order.
func MigrateDownAll(db *sql.DB) error {
	applied, err := GetAppliedMigrations(db)
	if err != nil {
		return err
	}

	var versions []string
	for v := range applied {
		versions = append(versions, v)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))

	for _, v := range versions {
		downFilename := v + ".down.sql"
		log.Printf("⏳ Rolling back: %s ...", downFilename)
		content, err := fs.ReadFile(migrations.FS, downFilename)
		if err != nil {
			return fmt.Errorf("failed reading %s: %w", downFilename, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("error in %s: %w", downFilename, err)
		}

		if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", v); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
		log.Printf("✅ Rolled back: %s", v)
	}
	return nil
}

// Status returns current migration status of all available migration files.
func Status(db *sql.DB) ([]MigrationStatus, error) {
	applied, err := GetAppliedMigrations(db)
	if err != nil {
		return nil, err
	}

	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return nil, err
	}

	var upFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			upFiles = append(upFiles, entry.Name())
		}
	}
	sort.Strings(upFiles)

	var statuses []MigrationStatus
	for _, filename := range upFiles {
		version := strings.TrimSuffix(filename, ".up.sql")
		appliedAt, isApplied := applied[version]
		var at *time.Time
		if isApplied {
			at = &appliedAt
		}
		statuses = append(statuses, MigrationStatus{
			Name:      version,
			IsApplied: isApplied,
			AppliedAt: at,
		})
	}
	return statuses, nil
}
