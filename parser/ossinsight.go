package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type OssInsightResponse struct {
	Data struct {
		Rows []struct {
			RepoName        string `json:"repo_name"`
			PrimaryLanguage string `json:"primary_language"`
			Stars           string `json:"stars"`
			Forks           string `json:"forks"`
		} `json:"rows"`
	} `json:"data"`
}

func GetTrendingReposFromOssInsight(maxRepos int, period, language string) ([]Repository, error) {
	if period == "" {
		period = "past_24_hours"
	}
	if language == "" {
		language = "All"
	}

	baseURL := "https://api.ossinsight.io/v1/trends/repos/"
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %v", err)
	}

	q := u.Query()
	q.Set("period", period)
	q.Set("language", language)
	q.Set("limit", fmt.Sprintf("%d", maxRepos))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	setBrowserHeaders(req)

	req.Header.Set("Accept", "application/json")

	res, err := doRequestWithRetry(req, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve trending repositories from OssInsight: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OssInsight API returned status: %s", res.Status)
	}

	var apiRes OssInsightResponse
	if err := json.NewDecoder(res.Body).Decode(&apiRes); err != nil {
		return nil, fmt.Errorf("failed to decode OssInsight response: %v", err)
	}

	var allRepos []Repository
	for _, row := range apiRes.Data.Rows {
		repoURL := "https://github.com/" + row.RepoName
		allRepos = append(allRepos, Repository{
			URL:      repoURL,
			Language: row.PrimaryLanguage,
			Stars:    row.Stars,
			Forks:    row.Forks,
		})
	}

	filteredRepos, err := FilterExistingRepos(allRepos)
	if err != nil {
		return nil, fmt.Errorf("failed to filter existing repositories: %v", err)
	}

	return filteredRepos, nil
}
