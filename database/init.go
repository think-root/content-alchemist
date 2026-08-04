package database

import (
	"content-alchemist/config"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
)

var DBThinkRoot *sql.DB

// SQLiteDriverName is the driver every connection to the app database must be
// opened with — including test helpers, which would otherwise miss the custom
// SQL functions registered below.
const SQLiteDriverName = "sqlite3_unicode"

var registerDriverOnce sync.Once

// registerSQLiteDriver installs a SQLite driver that adds unicode_lower().
// SQLite's built-in lower() and its case-insensitive LIKE only fold ASCII, so a
// search for "старий" would not match a stored "Старий". The descriptions in
// this database are Ukrainian, so substring search has to fold Unicode too.
func registerSQLiteDriver() {
	registerDriverOnce.Do(func() {
		sql.Register(SQLiteDriverName, &sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				return conn.RegisterFunc("unicode_lower", strings.ToLower, true)
			},
		})
	})
}

type GithubRepositories struct {
	ID              int64      `json:"id"`
	URL             string     `json:"url"`
	Text            string     `json:"text"`
	Posted          int        `json:"posted"`
	DateAdded       *time.Time `json:"date_added"`
	DatePosted      *time.Time `json:"date_posted"`
	PublishPriority *int64     `json:"publish_priority"`
}

func (GithubRepositories) TableName() string {
	return "github_repositories"
}

func init() {
	var err error

	registerSQLiteDriver()

	dbPath := config.SQLITE_DB_PATH
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Error creating database directory: %v", err)
	}

	DBThinkRoot, err = sql.Open(SQLiteDriverName, dbPath)
	if err != nil {
		log.Fatalf("Error opening SQLite connection: %v", err)
	}

	if _, err := DBThinkRoot.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		log.Fatalf("Error enabling WAL mode: %v", err)
	}

	if err := createTableIfNotExists(); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	if err := ensurePublicationQueueSchema(); err != nil {
		log.Fatalf("Error ensuring publication queue schema: %v", err)
	}
	if err := ensureArchiveSchema(); err != nil {
		log.Fatalf("Error ensuring archive schema: %v", err)
	}

	var userVersion int
	if err := DBThinkRoot.QueryRow("PRAGMA user_version;").Scan(&userVersion); err != nil {
		log.Fatalf("Error checking user_version: %v", err)
	}

	if userVersion == 0 {
		log.Println("Starting migration from PostgreSQL...")
		if err := migrateFromPostgres(); err != nil {
			log.Printf("Migration failed: %v", err)
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
		date_posted DATETIME,
		publish_priority INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_github_repositories_url ON github_repositories(url);
	CREATE INDEX IF NOT EXISTS idx_github_repositories_posted ON github_repositories(posted);
	`
	_, err := DBThinkRoot.Exec(createTableQuery)
	return err
}

func ensurePublicationQueueSchema() error {
	hasPublishPriority, err := hasColumn("github_repositories", "publish_priority")
	if err != nil {
		return err
	}
	if !hasPublishPriority {
		if _, err := DBThinkRoot.Exec("ALTER TABLE github_repositories ADD COLUMN publish_priority INTEGER"); err != nil {
			return fmt.Errorf("failed to add publish_priority column: %v", err)
		}
	}

	_, err = DBThinkRoot.Exec(`
		CREATE INDEX IF NOT EXISTS idx_github_repositories_publication_queue
			ON github_repositories(posted, publish_priority DESC, date_added ASC, id ASC);
	`)
	if err != nil {
		return fmt.Errorf("failed to create publication queue index: %v", err)
	}

	return nil
}

func ensureArchiveSchema() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS archived_repositories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_id INTEGER,
		url TEXT NOT NULL,
		text TEXT NOT NULL,
		date_added DATETIME,
		date_posted DATETIME,
		date_archived DATETIME NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_archived_repositories_url ON archived_repositories(url);
	CREATE INDEX IF NOT EXISTS idx_archived_repositories_date_archived ON archived_repositories(date_archived DESC);
	CREATE INDEX IF NOT EXISTS idx_archived_repositories_date_posted ON archived_repositories(date_posted);
	`
	if _, err := DBThinkRoot.Exec(createTableQuery); err != nil {
		return fmt.Errorf("failed to create archived_repositories table: %v", err)
	}

	return nil
}

func hasColumn(tableName, columnName string) (bool, error) {
	rows, err := DBThinkRoot.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return false, fmt.Errorf("failed to inspect table %s: %v", tableName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name string
		var columnType string
		var notNull int
		var defaultValue interface{}
		var primaryKey int
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return false, fmt.Errorf("failed to scan table info for %s: %v", tableName, err)
		}
		if name == columnName {
			return true, nil
		}
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("failed to iterate table info for %s: %v", tableName, err)
	}

	return false, nil
}

func buildPostgreSQLDSN(host, port, user, password, dbname string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		host, port, user, password, dbname)
}

func migrateFromPostgres() error {
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
