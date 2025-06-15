package database

import (
	"fmt"
	"time"
)

func UpdatePostedStatusByURL(url string, posted bool) error {
	// First check if repository exists
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM alchemist_github_repositories WHERE url = $1)"
	err := DBThinkRoot.QueryRow(checkQuery, url).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking repository existence: %v", err)
	}
	
	if !exists {
		return fmt.Errorf("repository with URL %s not found", url)
	}

	// Update the repository
	var postedValue int
	var datePosted *time.Time
	
	if posted {
		postedValue = 1
		now := time.Now()
		datePosted = &now
	} else {
		postedValue = 0
		datePosted = nil
	}

	updateQuery := "UPDATE alchemist_github_repositories SET posted = $1, date_posted = $2 WHERE url = $3"
	_, err = DBThinkRoot.Exec(updateQuery, postedValue, datePosted, url)
	if err != nil {
		return fmt.Errorf("error updating repository: %v", err)
	}

	return nil
}
