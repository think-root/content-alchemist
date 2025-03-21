package database

import (
	"fmt"
	"time"
)

func AddRepositoryToDB(url, text string) error {
	now := time.Now()
	repository := GithubRepositories{
		URL:       url,
		Text:      text,
		DateAdded: &now,
	}
	result := DBThinkRoot.Create(&repository)
	if result.Error != nil {
		return fmt.Errorf("error adding post to DB: %v", result.Error)
	}
	return nil
}
