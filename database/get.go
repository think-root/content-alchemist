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
		if sortOrder == "ASC" {
			orderBy = "date_posted IS NULL DESC, date_posted ASC"
		} else if sortOrder == "DESC" {
			orderBy = "date_posted IS NULL ASC, date_posted DESC"
		}
	case "date_added":
		if sortOrder == "ASC" {
			orderBy = "date_added ASC"
		} else if sortOrder == "DESC" {
			orderBy = "date_added DESC"
		}
	case "id":
		if sortOrder == "ASC" {
			orderBy = "id ASC"
		} else if sortOrder == "DESC" {
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
