package database

import (
	"fmt"
)

func GetRepository(limit int, posted *bool) ([]GithubRepositories, error) {
	var repositories []GithubRepositories

	query := DBThinkRoot.Model(&GithubRepositories{}).Order("id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}

	if posted != nil {
		query = query.Where("posted = ?", *posted)
	}

	result := query.Find(&repositories)
	if result.Error != nil {
		return nil, fmt.Errorf("error fetching repositories from DB: %v", result.Error)
	}

	return repositories, nil
}
