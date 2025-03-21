package database

import (
	"fmt"
)

func GetRepository(limit int, posted *bool, sortBy string, sortOrder string) ([]GithubRepositories, error) {
	var repositories []GithubRepositories

	query := DBThinkRoot.Model(&GithubRepositories{})
	
	if posted != nil {
		query = query.Where("posted = ?", *posted)
	}
	
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
	
	query = query.Order(orderBy)
	
	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&repositories)
	if result.Error != nil {
		return nil, fmt.Errorf("error fetching repositories from DB: %v", result.Error)
	}

	return repositories, nil
}
