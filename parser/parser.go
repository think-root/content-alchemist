package parser

import (
	"content-alchemist/database"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// httpClient is a shared HTTP client with timeout for all requests
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// browserHeaders adds common browser headers to avoid being blocked by GitHub
func setBrowserHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
}

// doRequestWithRetry performs an HTTP request with retry logic and exponential backoff
func doRequestWithRetry(req *http.Request, maxRetries int) (*http.Response, error) {
	var res *http.Response
	var err error

	for i := 0; i < maxRetries; i++ {
		res, err = httpClient.Do(req)
		if err == nil && res.StatusCode == http.StatusOK {
			return res, nil
		}

		if res != nil {
			res.Body.Close()
		}

		if i < maxRetries-1 {
			backoff := time.Duration(i+1) * 2 * time.Second
			log.Printf("Request failed (attempt %d/%d), retrying in %v: %v", i+1, maxRetries, backoff, err)
			time.Sleep(backoff)

			// Recreate the request for retry (body might be consumed)
			req, _ = http.NewRequest(req.Method, req.URL.String(), nil)
			setBrowserHeaders(req)
		}
	}

	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("request failed after %d retries, last status: %d", maxRetries, res.StatusCode)
}

type Repository struct {
	URL      string `json:"url"`
	Language string `json:"language"`
	Stars    string `json:"stars"`
	Forks    string `json:"forks"`
}

func GetTrendingRepos(maxRepos int, since, spokenLanguageCode string) ([]Repository, error) {
	url := fmt.Sprintf("https://github.com/trending?since=%s&spoken_language_code=%s", since, spokenLanguageCode)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("failed to create request: " + err.Error())
	}
	setBrowserHeaders(req)

	res, err := doRequestWithRetry(req, 3)
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

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %w", err)
	}
	// For GitHub API, we use slightly different headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := doRequestWithRetry(req, 3)
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
