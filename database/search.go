package database

import "fmt"

func SearchPostInDB(link string) (bool, error) {
	var count int64
	query := "SELECT COUNT(*) FROM github_repositories WHERE url = ?"
	
	err := DBThinkRoot.QueryRow(query, link).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error searching post in DB: %v", err)
	}
	return count > 0, nil
}
