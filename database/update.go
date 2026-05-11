package database

import (
	"database/sql"
	"fmt"
	"time"
)

func UpdatePostedStatusByURL(url string, posted bool) error {
	// First check if repository exists
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM github_repositories WHERE url = ?)"
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

	updateQuery := "UPDATE github_repositories SET posted = ?, date_posted = ?, publish_priority = NULL WHERE url = ?"
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
			UPDATE github_repositories
			SET text = ?
			WHERE id = ?
			RETURNING id, url, text, posted, date_added, date_posted, publish_priority
		`
	} else {
		query = `
			UPDATE github_repositories
			SET text = ?
			WHERE url = ?
			RETURNING id, url, text, posted, date_added, date_posted, publish_priority
		`
	}

	err := DBThinkRoot.QueryRow(query, text, identifier).Scan(
		&repo.ID,
		&repo.URL,
		&repo.Text,
		&repo.Posted,
		&repo.DateAdded,
		&repo.DatePosted,
		&repo.PublishPriority,
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

func PromoteRepositoryToNextByIDOrURL(identifier string, isID bool) (*GithubRepositories, error) {
	existingRepo, err := GetRepositoryByIDOrURL(identifier, isID)
	if err != nil {
		return nil, err
	}

	if existingRepo.Posted != 0 {
		return nil, fmt.Errorf("repository is already posted")
	}

	var repo GithubRepositories
	var query string

	if isID {
		query = `
			UPDATE github_repositories
			SET publish_priority = (
				SELECT COALESCE(MAX(publish_priority), 0) + 1
				FROM github_repositories
				WHERE posted = 0
			)
			WHERE id = ? AND posted = 0
			RETURNING id, url, text, posted, date_added, date_posted, publish_priority
		`
	} else {
		query = `
			UPDATE github_repositories
			SET publish_priority = (
				SELECT COALESCE(MAX(publish_priority), 0) + 1
				FROM github_repositories
				WHERE posted = 0
			)
			WHERE url = ? AND posted = 0
			RETURNING id, url, text, posted, date_added, date_posted, publish_priority
		`
	}

	err = DBThinkRoot.QueryRow(query, identifier).Scan(
		&repo.ID,
		&repo.URL,
		&repo.Text,
		&repo.Posted,
		&repo.DateAdded,
		&repo.DatePosted,
		&repo.PublishPriority,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("repository is already posted")
		}
		return nil, fmt.Errorf("error promoting repository: %v", err)
	}

	return &repo, nil
}

func DeleteRepositoryByIDOrURL(identifier string, isID bool) error {
	var query string

	if isID {
		query = "DELETE FROM github_repositories WHERE id = ?"
	} else {
		query = "DELETE FROM github_repositories WHERE url = ?"
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
