package database

import (
	"fmt"
	"time"
)

func AddRepositoryToDB(url, text string) error {
	now := time.Now()
	query := `
		INSERT INTO alchemist_github_repositories (url, text, date_added)
		VALUES ($1, $2, $3)
		ON CONFLICT (url) DO NOTHING
	`
	
	_, err := DBThinkRoot.Exec(query, url, text, now)
	if err != nil {
		return fmt.Errorf("error adding post to DB: %v", err)
	}
	return nil
}
