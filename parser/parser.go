package parser

import (
	"content-alchemist/config"
	"content-alchemist/database"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

// httpClient is a shared HTTP client with timeout for all requests
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// ReadmeNotFoundError indicates that the repository doesn't have a README file
type ReadmeNotFoundError struct {
	Repo string
}

func (e *ReadmeNotFoundError) Error() string {
	return fmt.Sprintf("repository %s does not have a README file", e.Repo)
}

// ReadmeHTTPError indicates an HTTP error occurred while fetching the README
type ReadmeHTTPError struct {
	Repo       string
	StatusCode int
	Status     string
}

func (e *ReadmeHTTPError) Error() string {
	return fmt.Sprintf("HTTP error fetching README for %s: %d %s", e.Repo, e.StatusCode, e.Status)
}

// Precompiled regexps used to strip non-meaningful markup from a README when
// estimating how much real textual content it carries.
var (
	reCodeFence  = regexp.MustCompile("(?s)```.*?```")
	reInlineCode = regexp.MustCompile("`[^`]*`")
	reMdImage    = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	reMdLink     = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	reHTMLTag    = regexp.MustCompile(`(?s)<[^>]*>`)
	reBareURL    = regexp.MustCompile(`https?://\S+`)
	reMdMarkup   = regexp.MustCompile("[#>*_~` \\t\\-=|!]+")
	reWhitespace = regexp.MustCompile(`\s+`)
)

// MeaningfulContentLength returns the length (in runes) of the "meaningful"
// text of a README after removing markdown/HTML markup, links, images, code
// blocks and bare URLs. It is used as a heuristic to decide whether a README
// carries enough substance to be worth sending to the LLM.
func MeaningfulContentLength(readme string) int {
	return utf8.RuneCountInString(meaningfulText(readme))
}

// meaningfulText strips markdown/HTML markup, links, images, code blocks and
// bare URLs from a README, leaving only its prose. Both the content-length and
// the language heuristics operate on this cleaned text.
func meaningfulText(readme string) string {
	text := readme

	// Remove code blocks first (they may contain URLs / markup we don't want).
	text = reCodeFence.ReplaceAllString(text, " ")
	text = reInlineCode.ReplaceAllString(text, " ")

	// Images carry no textual meaning; drop them entirely.
	text = reMdImage.ReplaceAllString(text, " ")

	// Links: keep the visible label, drop the URL target.
	text = reMdLink.ReplaceAllString(text, "$1")

	// HTML tags (badges, <img>, <a>, etc.) and any remaining bare URLs.
	text = reHTMLTag.ReplaceAllString(text, " ")
	text = reBareURL.ReplaceAllString(text, " ")

	// Strip leftover markdown punctuation/markup and collapse whitespace.
	text = reMdMarkup.ReplaceAllString(text, " ")
	text = reWhitespace.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}

// maxReadmeBytes caps how much of a README is kept before sending it to the LLM.
const maxReadmeBytes = 70000

// minLettersForLanguageCheck is the number of letters a README must contain
// before its language can be judged. Below that there is nothing to measure —
// the content-length heuristic is the relevant guard for such READMEs.
const minLettersForLanguageCheck = 30

// NonLatinLetterRatio reports the share of non-Latin letters (Han, Hiragana,
// Katakana, Hangul, Cyrillic, Arabic, ...) among all letters of the README's
// meaningful text, together with the total number of letters found.
func NonLatinLetterRatio(readme string) (float64, int) {
	var letters, nonLatin int

	for _, r := range meaningfulText(readme) {
		if !unicode.IsLetter(r) {
			continue
		}
		letters++
		if !unicode.Is(unicode.Latin, r) {
			nonLatin++
		}
	}

	if letters == 0 {
		return 0, 0
	}

	return float64(nonLatin) / float64(letters), letters
}

// IsEnglishReadme reports whether a README is predominantly written in a
// Latin-script (i.e. presumably English) language. It also returns the measured
// non-Latin ratio so callers can report it. READMEs with too few letters to
// judge are accepted.
func IsEnglishReadme(readme string) (bool, float64) {
	ratio, letters := NonLatinLetterRatio(readme)
	if letters < minLettersForLanguageCheck {
		return true, ratio
	}

	return ratio*100 <= float64(config.README_MAX_NON_LATIN_PERCENT), ratio
}

// reservedGitHubPaths are first path segments of github.com that can never be a
// repository owner.
var reservedGitHubPaths = map[string]bool{
	"orgs":          true,
	"settings":      true,
	"topics":        true,
	"sponsors":      true,
	"features":      true,
	"marketplace":   true,
	"apps":          true,
	"collections":   true,
	"trending":      true,
	"explore":       true,
	"notifications": true,
	"pulls":         true,
	"issues":        true,
	"search":        true,
}

// NormalizeRepoURL converts any form of GitHub repository reference into the
// canonical "https://github.com/<owner>/<repo>" form. It tolerates a missing
// scheme, "http://", a "www." host, tracking query parameters, fragments, a
// ".git" suffix, a trailing slash and deep paths such as "/tree/main/src".
func NormalizeRepoURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("empty repository URL")
	}

	// url.Parse treats a scheme-less value as a bare path, so add one.
	if !strings.Contains(trimmed, "://") {
		trimmed = "https://" + trimmed
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", fmt.Errorf("not a valid URL: %w", err)
	}

	switch parsed.Scheme {
	case "http", "https":
	default:
		return "", fmt.Errorf("unsupported URL scheme %q", parsed.Scheme)
	}

	host := strings.ToLower(parsed.Hostname())
	if host != "github.com" && host != "www.github.com" {
		return "", fmt.Errorf("not a github.com URL: %q", host)
	}

	segments := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(segments) < 2 || segments[0] == "" || segments[1] == "" {
		return "", fmt.Errorf("URL does not point to a repository")
	}

	owner := segments[0]
	repo := strings.TrimSuffix(segments[1], ".git")

	if reservedGitHubPaths[strings.ToLower(owner)] {
		return "", fmt.Errorf("%q is not a repository owner", owner)
	}
	if repo == "" {
		return "", fmt.Errorf("URL does not point to a repository")
	}

	return fmt.Sprintf("https://github.com/%s/%s", owner, repo), nil
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

		// Don't retry on 404 - the resource doesn't exist
		if err == nil && res.StatusCode == http.StatusNotFound {
			return res, nil
		}

		if res != nil {
			res.Body.Close()
		}

		if i < maxRetries-1 {
			backoff := time.Duration(3*(1<<i)) * time.Second
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
	if normalized, err := NormalizeRepoURL(repo); err == nil {
		repo = strings.TrimPrefix(normalized, "https://github.com/")
	} else {
		repo = strings.TrimPrefix(repo, "https://github.com/")
	}
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

	// Add GitHub token for authentication if available
	if config.GITHUB_TOKEN != "" {
		req.Header.Set("Authorization", "token "+config.GITHUB_TOKEN)
	}

	resp, err := doRequestWithRetry(req, 3)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Return specific error for 404
	if resp.StatusCode == http.StatusNotFound {
		return "", &ReadmeNotFoundError{Repo: repo}
	}

	// Return specific error for other non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return "", &ReadmeHTTPError{
			Repo:       repo,
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
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
	if len(content) > maxReadmeBytes {
		// Cut on a rune boundary so the truncated README stays valid UTF-8.
		content = content[:maxReadmeBytes]
		for len(content) > 0 && !utf8.ValidString(content) {
			content = content[:len(content)-1]
		}
	}

	return content, nil
}
