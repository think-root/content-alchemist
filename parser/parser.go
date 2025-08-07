package parser

import (
	"content-alchemist/database"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Repository struct {
	URL      string `json:"url"`
	Language string `json:"language"`
	Stars    string `json:"stars"`
	Forks    string `json:"forks"`
}

func GetTrendingRepos(maxRepos int, since, spokenLanguageCode string) ([]Repository, error) {
	url := fmt.Sprintf("https://github.com/trending?since=%s&spoken_language_code=%s", since, spokenLanguageCode)

	res, err := http.Get(url)
	if err != nil {
		return nil, errors.New("failed to retrieve trending repositories: " + err.Error())
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.New("failed to parse HTML document: " + err.Error())
	}

	var allRepos []Repository

	doc.Find(".Box-row").Each(func(i int, s *goquery.Selection) {
		repoURL := "https://github.com" + s.Find(".lh-condensed > a").AttrOr("href", "")
		allRepos = append(allRepos, Repository{
			URL: repoURL,
		})
	})

	filteredRepos, err := FilterExistingRepos(allRepos)
	if err != nil {
		return nil, fmt.Errorf("failed to filter existing repositories: %v", err)
	}

	if len(filteredRepos) > maxRepos {
		filteredRepos = filteredRepos[:maxRepos]
	}

	return filteredRepos, nil
}

func FilterExistingRepos(repos []Repository) ([]Repository, error) {
	var filteredRepos []Repository
	countAll := 0
	for _, repo := range repos {
		exists, err := database.SearchPostInDB(repo.URL)
		if err != nil {
			return nil, fmt.Errorf("error checking repository existence for URL %s: %v", repo.URL, err)
		}
		countAll += 1
		if !exists {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	return filteredRepos, nil
}

func FilterExistingURLs(urls []string) ([]string, error) {
	var filteredURLs []string

	for _, url := range urls {
		exists, err := database.SearchPostInDB(url)
		if err != nil {
			return nil, fmt.Errorf("error checking URL existence for %s: %v", url, err)
		}

		if !exists {
			filteredURLs = append(filteredURLs, url)
		}
	}

	return filteredURLs, nil
}

func GetRepoReadme(repo string) (string, error) {
	repo = strings.TrimPrefix(repo, "https://github.com/")
	url := fmt.Sprintf("https://api.github.com/repos/%s/readme", repo)

	type gitHubAPIResponse struct {
		Content string `json:"content"`
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var readme gitHubAPIResponse
	if err := json.Unmarshal(body, &readme); err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	decodedContent, err := base64.StdEncoding.DecodeString(readme.Content)
	if err != nil {
		return "", fmt.Errorf("error decoding Base64 content: %w", err)
	}

	content := string(decodedContent)
	if len(content) > 70000 {
		content = content[:70000]
	}

	return content, nil
}
