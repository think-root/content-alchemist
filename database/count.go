package database

import (
	"fmt"
)

func CountRepositories(posted *bool) (int, error) {
	var count int64
	query := DBThinkRoot.Model(&GithubRepositories{})

	if posted != nil {
		query = query.Where("posted = ?", *posted)
	}

	result := query.Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("error counting repositories: %v", result.Error)
	}

	return int(count), nil
}
