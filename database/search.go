package database

import "fmt"

func SearchPostInDB(link string) (bool, error) {
	var count int64
	query := "SELECT COUNT(*) FROM alchemist_github_repositories WHERE url = $1"
	
	err := DBThinkRoot.QueryRow(query, link).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error searching post in DB: %v", err)
	}
	return count > 0, nil
}
