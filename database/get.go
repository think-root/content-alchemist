package database

import (
	"fmt"
)

func GetRepository(limit int, offset int, posted *bool, sortBy string, sortOrder string) ([]GithubRepositories, int, error) {
	var repositories []GithubRepositories
	var totalCount int64

	query := DBThinkRoot.Model(&GithubRepositories{})

	if posted != nil {
		query = query.Where("posted = ?", *posted)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting total repositories: %v", err)
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

	if offset > 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&repositories)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("error fetching repositories from DB: %v", result.Error)
	}

	return repositories, int(totalCount), nil
}
