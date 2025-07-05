package database

import (
	"fmt"
)

func GetRepository(limit int, offset int, posted *bool, sortBy string, sortOrder string) ([]GithubRepositories, int, error) {
	var repositories []GithubRepositories
	var totalCount int64

	// Build WHERE clause
	var whereClause string
	var args []interface{}
	argIndex := 1

	if posted != nil {
		var postedValue int
		if *posted {
			postedValue = 1
		} else {
			postedValue = 0
		}
		whereClause = fmt.Sprintf("WHERE posted = $%d", argIndex)
		args = append(args, postedValue)
		argIndex++
	}

	// Build ORDER BY clause
	orderBy := "id DESC"
	switch sortBy {
	case "date_posted":
		switch sortOrder {
		case "ASC":
			orderBy = "date_posted IS NULL DESC, date_posted ASC"
		case "DESC":
			orderBy = "date_posted IS NULL ASC, date_posted DESC"
		}
	case "date_added":
		switch sortOrder {
		case "ASC":
			orderBy = "date_added ASC"
		case "DESC":
			orderBy = "date_added DESC"
		}
	case "id":
		switch sortOrder {
		case "ASC":
			orderBy = "id ASC"
		case "DESC":
			orderBy = "id DESC"
		}
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM alchemist_github_repositories %s", whereClause)
	err := DBThinkRoot.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting total repositories: %v", err)
	}

	// Build data query
	dataQuery := fmt.Sprintf("SELECT id, url, text, posted, date_added, date_posted FROM alchemist_github_repositories %s ORDER BY %s", whereClause, orderBy)

	// Add LIMIT and OFFSET
	if limit > 0 {
		dataQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++
	}

	if offset > 0 {
		dataQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
	}

	// Execute data query
	rows, err := DBThinkRoot.Query(dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching repositories from DB: %v", err)
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		var repo GithubRepositories
		err := rows.Scan(&repo.ID, &repo.URL, &repo.Text, &repo.Posted, &repo.DateAdded, &repo.DatePosted)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning repository: %v", err)
		}
		repositories = append(repositories, repo)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating over rows: %v", err)
	}

	return repositories, int(totalCount), nil
}

// GetRepositoryByIDOrURL retrieves a single repository by ID or URL
func GetRepositoryByIDOrURL(identifier string, isID bool) (*GithubRepositories, error) {
	var repo GithubRepositories
	var query string

	if isID {
		query = `
			SELECT id, url, text, posted, date_added, date_posted
			FROM alchemist_github_repositories
			WHERE id = $1
		`
	} else {
		query = `
			SELECT id, url, text, posted, date_added, date_posted
			FROM alchemist_github_repositories
			WHERE url = $1
		`
	}

	err := DBThinkRoot.QueryRow(query, identifier).Scan(
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
		return nil, fmt.Errorf("error fetching repository: %v", err)
	}

	return &repo, nil
}
