package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func UpdatePostedStatusByURL(url string, posted bool) error {
	var repository GithubRepositories
	result := DBThinkRoot.Where("url = ?", url).First(&repository)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("repository with URL %s not found", url)
		}
		return fmt.Errorf("error finding repository: %v", result.Error)
	}

	repository.Posted = posted
	if posted {
		now := time.Now()
		repository.DatePosted = &now
	} else {
		repository.DatePosted = nil
	}
	result = DBThinkRoot.Save(&repository)
	if result.Error != nil {
		return fmt.Errorf("error updating repository: %v", result.Error)
	}

	return nil
}
