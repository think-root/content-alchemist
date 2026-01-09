package database

import (
	"content-alchemist/config"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DBThinkRoot *sql.DB

type GithubRepositories struct {
	ID         int64      `json:"id"`
	URL        string     `json:"url"`
	Text       string     `json:"text"`
	Posted     int        `json:"posted"`
	DateAdded  *time.Time `json:"date_added"`
	DatePosted *time.Time `json:"date_posted"`
}

func (GithubRepositories) TableName() string {
	return "github_repositories"
}

func init() {
	var err error

	// Ensure data directory exists
	dbPath := config.SQLITE_DB_PATH
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Error creating database directory: %v", err)
	}

	// Connect to SQLite
	DBThinkRoot, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening SQLite connection: %v", err)
	}

	// Enable WAL mode
	if _, err := DBThinkRoot.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		log.Fatalf("Error enabling WAL mode: %v", err)
	}

	if err := createTableIfNotExists(); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	// Check migration status
	var userVersion int
	if err := DBThinkRoot.QueryRow("PRAGMA user_version;").Scan(&userVersion); err != nil {
		log.Fatalf("Error checking user_version: %v", err)
	}

	if userVersion == 0 {
		log.Println("Starting migration from PostgreSQL...")
		if err := migrateFromPostgres(); err != nil {
			log.Printf("Migration failed: %v", err)
			// Decide if we want to fail hard or just log. Failing hard is probably safer for data consistency.
			// However, if Postgres connects fails (e.g. credentials missing), maybe we should just continue?
			// The user said "migrate automatically", implying if it's possible. All existing logic relies on the data.
			// Let's log error but allow startup, assuming empty DB is better than crash if PG is down.
			// BUT, the requirements say "migrate... then lift version". So if migration fails, we don't lift version.
		} else {
			log.Println("Migration successful.")
			if _, err := DBThinkRoot.Exec("PRAGMA user_version = 1;"); err != nil {
				log.Printf("Error updating user_version: %v", err)
			}
		}
	} else {
		log.Println("Database already migrated (version >= 1).")
	}

	log.Printf("Successfully connected to SQLite database: %s", config.SQLITE_DB_PATH)
}

func createTableIfNotExists() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS github_repositories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL UNIQUE,
		text TEXT NOT NULL,
		posted INTEGER NOT NULL DEFAULT 0,
		date_added DATETIME,
		date_posted DATETIME
	);
	CREATE INDEX IF NOT EXISTS idx_github_repositories_url ON github_repositories(url);
	CREATE INDEX IF NOT EXISTS idx_github_repositories_posted ON github_repositories(posted);
	`
	_, err := DBThinkRoot.Exec(createTableQuery)
	return err
}

func buildPostgreSQLDSN(host, port, user, password, dbname string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		host, port, user, password, dbname)
}

func migrateFromPostgres() error {
	// Check if Postgres config is available
	if config.POSTGRES_HOST == "" || config.POSTGRES_DB == "" {
		log.Println("PostgreSQL configuration missing, skipping migration.")
		return fmt.Errorf("postgres config missing")
	}

	dsn := buildPostgreSQLDSN(
		config.POSTGRES_HOST,
		config.POSTGRES_PORT,
		config.POSTGRES_USER,
		config.POSTGRES_PASSWORD,
		config.POSTGRES_DB,
	)

	pgDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open postgres connection: %v", err)
	}
	defer pgDB.Close()

	if err := pgDB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to postgres: %v", err)
	}

	rows, err := pgDB.Query("SELECT url, text, posted, date_added, date_posted FROM alchemist_github_repositories")
	if err != nil {
		return fmt.Errorf("failed to query postgres data: %v", err)
	}
	defer rows.Close()

	tx, err := DBThinkRoot.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin sqlite transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO github_repositories (url, text, posted, date_added, date_posted) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmt.Close()

	count := 0
	for rows.Next() {
		var url, text string
		var posted int
		var dateAdded, datePosted *time.Time

		if err := rows.Scan(&url, &text, &posted, &dateAdded, &datePosted); err != nil {
			return fmt.Errorf("failed to scan postgres row: %v", err)
		}

		if _, err := stmt.Exec(url, text, posted, dateAdded, datePosted); err != nil {
			return fmt.Errorf("failed to insert row into sqlite: %v", err)
		}
		count++
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Migrated %d records from PostgreSQL.", count)
	return nil
}
