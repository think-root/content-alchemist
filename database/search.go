package database

import "fmt"

func SearchPostInDB(link string) (bool, error) {
	var count int64
	result := DBThinkRoot.Model(&GithubRepositories{}).Where("url = ?", link).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("error searching post in DB: %v", result.Error)
	}
	return count > 0, nil
}
