package database

import (
	"fmt"
)

func CountRepositories(posted *bool) (int, error) {
	var count int64
	var query string
	var args []interface{}

	if posted != nil {
		var postedValue int
		if *posted {
			postedValue = 1
		} else {
			postedValue = 0
		}
		query = "SELECT COUNT(*) FROM github_repositories WHERE posted = ?"
		args = []interface{}{postedValue}
	} else {
		query = "SELECT COUNT(*) FROM github_repositories"
		args = []interface{}{}
	}

	err := DBThinkRoot.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting repositories: %v", err)
	}

	return int(count), nil
}
