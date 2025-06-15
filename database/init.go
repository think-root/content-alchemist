package database

import (
	"content-alchemist/config"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
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
	return "alchemist_github_repositories"
}

func buildPostgreSQLDSN(host, port, user, password, dbname string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		host, port, user, password, dbname)
}

func createDatabaseIfNotExists() error {
	// Connect to PostgreSQL without specifying a database
	defaultDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		config.POSTGRES_HOST, config.POSTGRES_PORT, config.POSTGRES_USER, config.POSTGRES_PASSWORD)
	
	db, err := sql.Open("postgres", defaultDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Check if database exists
	var exists bool
	query := "SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)"
	err = db.QueryRow(query, config.POSTGRES_DB).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %v", err)
	}

	// Create database if it doesn't exist
	if !exists {
		createQuery := fmt.Sprintf("CREATE DATABASE \"%s\" WITH ENCODING='UTF8' LC_COLLATE='en_US.UTF-8' LC_CTYPE='en_US.UTF-8'", config.POSTGRES_DB)
		_, err = db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("failed to create database: %v", err)
		}
		log.Printf("Created database: %s", config.POSTGRES_DB)
	} else {
		log.Printf("Database %s already exists", config.POSTGRES_DB)
	}

	return nil
}

func createTableIfNotExists() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS alchemist_github_repositories (
		id BIGSERIAL PRIMARY KEY,
		url TEXT NOT NULL,
		text TEXT NOT NULL,
		posted INTEGER NOT NULL DEFAULT 0,
		date_added TIMESTAMP,
		date_posted TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_alchemist_github_repositories_url ON alchemist_github_repositories(url);
	CREATE INDEX IF NOT EXISTS idx_alchemist_github_repositories_posted ON alchemist_github_repositories(posted);
	CREATE INDEX IF NOT EXISTS idx_alchemist_github_repositories_date_added ON alchemist_github_repositories(date_added);
	CREATE INDEX IF NOT EXISTS idx_alchemist_github_repositories_date_posted ON alchemist_github_repositories(date_posted);
	`

	_, err := DBThinkRoot.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	log.Println("Table alchemist_github_repositories is ready")
	return nil
}

func synchronizeSequence() error {
	// Get the maximum ID from the table
	var maxID int64
	err := DBThinkRoot.QueryRow("SELECT COALESCE(MAX(id), 0) FROM alchemist_github_repositories").Scan(&maxID)
	if err != nil {
		return fmt.Errorf("failed to get max ID: %v", err)
	}

	// If there are records, synchronize the sequence
	if maxID > 0 {
		var currentSeqVal int64
		err = DBThinkRoot.QueryRow("SELECT last_value FROM alchemist_github_repositories_id_seq").Scan(&currentSeqVal)
		if err != nil {
			return fmt.Errorf("failed to get sequence value: %v", err)
		}

		// Only update if sequence is behind
		if currentSeqVal < maxID {
			_, err = DBThinkRoot.Exec("SELECT setval('alchemist_github_repositories_id_seq', $1)", maxID)
			if err != nil {
				return fmt.Errorf("failed to set sequence value: %v", err)
			}
			log.Printf("Synchronized sequence: updated from %d to %d", currentSeqVal, maxID)
		} else {
			log.Printf("Sequence is already synchronized (current: %d, max_id: %d)", currentSeqVal, maxID)
		}
	}

	return nil
}

func init() {
	var err error

	// Create database if it doesn't exist
	err = createDatabaseIfNotExists()
	if err != nil {
		log.Printf("Error creating database: %v", err)
		return
	}

	// Build PostgreSQL DSN
	dsn := buildPostgreSQLDSN(
		config.POSTGRES_HOST,
		config.POSTGRES_PORT,
		config.POSTGRES_USER,
		config.POSTGRES_PASSWORD,
		config.POSTGRES_DB,
	)

	// Connect to the specific database
	DBThinkRoot, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Error opening PostgreSQL connection: %v", err)
		return
	}

	// Test the connection
	err = DBThinkRoot.Ping()
	if err != nil {
		log.Printf("Error pinging PostgreSQL database: %v", err)
		return
	}

	// Configure connection pool
	DBThinkRoot.SetMaxOpenConns(25)
	DBThinkRoot.SetMaxIdleConns(10)
	DBThinkRoot.SetConnMaxLifetime(time.Hour)
	DBThinkRoot.SetConnMaxIdleTime(30 * time.Minute)

	log.Printf("Successfully connected to PostgreSQL database: %s", config.POSTGRES_DB)

	// Create table if it doesn't exist
	err = createTableIfNotExists()
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	// Ensure sequence is synchronized with existing data (prevents duplicate key errors after migration)
	err = synchronizeSequence()
	if err != nil {
		log.Printf("Warning: Could not synchronize sequence: %v", err)
	}
}
