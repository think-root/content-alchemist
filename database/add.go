package database

import "fmt"

func AddRepositoryToDB(url, text string) error {
	repository := GithubRepositories{
		URL:    url,
		Text:   text,
	}
	result := DBThinkRoot.Create(&repository)
	if result.Error != nil {
		return fmt.Errorf("error adding post to DB: %v", result.Error)
	}
	return nil
}
