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

func UpdateRepositoryTextByIDOrURL(identifier, text string, isID bool) (*GithubRepositories, error) {
	var repo GithubRepositories
	var query string
	
	if isID {
		query = `
			UPDATE alchemist_github_repositories
			SET text = $1
			WHERE id = $2
			RETURNING id, url, text, posted, date_added, date_posted
		`
	} else {
		query = `
			UPDATE alchemist_github_repositories
			SET text = $1
			WHERE url = $2
			RETURNING id, url, text, posted, date_added, date_posted
		`
	}
	
	err := DBThinkRoot.QueryRow(query, text, identifier).Scan(
		&repo.ID,
		&repo.URL,
		&repo.Text,
		&repo.Posted,
		&repo.DateAdded,
		&repo.DatePosted,
	)
	
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			if isID {
				return nil, fmt.Errorf("repository with ID %s not found", identifier)
			} else {
				return nil, fmt.Errorf("repository with URL %s not found", identifier)
			}
		}
		return nil, fmt.Errorf("error updating repository: %v", err)
	}
	
	return &repo, nil
}

func DeleteRepositoryByIDOrURL(identifier string, isID bool) error {
	var query string
	
	if isID {
		query = "DELETE FROM alchemist_github_repositories WHERE id = $1"
	} else {
		query = "DELETE FROM alchemist_github_repositories WHERE url = $1"
	}
	
	result, err := DBThinkRoot.Exec(query, identifier)
	if err != nil {
		return fmt.Errorf("error deleting repository: %v", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %v", err)
	}
	
	if rowsAffected == 0 {
		if isID {
			return fmt.Errorf("repository with ID %s not found", identifier)
		} else {
			return fmt.Errorf("repository with URL %s not found", identifier)
		}
	}
	
	return nil
}
