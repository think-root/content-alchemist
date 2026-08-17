package routers

import (
	"content-alchemist/database"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupRepositoryTestDB(t *testing.T) {
	t.Helper()

	originalDB := database.DBThinkRoot
	db, err := sql.Open("sqlite3", t.TempDir()+"/content-alchemist-test.db")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	database.DBThinkRoot = db
	t.Cleanup(func() {
		db.Close()
		database.DBThinkRoot = originalDB
	})

	_, err = db.Exec(`
		CREATE TABLE github_repositories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL UNIQUE,
			text TEXT NOT NULL,
			posted INTEGER NOT NULL DEFAULT 0,
			date_added DATETIME,
			date_posted DATETIME,
			publish_priority INTEGER
		)`)
	if err != nil {
		t.Fatalf("failed to create test schema: %v", err)
	}
}

func postGetRepository(t *testing.T, body string) (*httptest.ResponseRecorder, getRepositoryResponse) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/think-root/api/get-repository/", strings.NewReader(body))
	rec := httptest.NewRecorder()

	GetRepository(rec, req)

	var envelope struct {
		Data    getRepositoryResponse `json:"data"`
		Message string                `json:"message"`
		Status  string                `json:"status"`
	}
	if rec.Code == http.StatusOK {
		if err := json.Unmarshal(rec.Body.Bytes(), &envelope); err != nil {
			t.Fatalf("failed to decode response %q: %v", rec.Body.String(), err)
		}
	}

	return rec, envelope.Data
}

// The url filter is what makes a manual retry possible: it must return the
// requested repository even though it is already posted, with the text resolved
// to the requested language.
func TestGetRepositoryByURLOrID(t *testing.T) {
	setupRepositoryTestDB(t)

	const wantURL = "https://github.com/resemble-ai/chatterbox"
	multilingual := "===(en)English description===(uk)Український опис==="

	if _, err := database.DBThinkRoot.Exec(
		"INSERT INTO github_repositories (url, text, posted, date_added) VALUES (?, ?, 1, '2026-08-17T16:00:00Z')",
		wantURL, multilingual,
	); err != nil {
		t.Fatalf("failed to insert test repository: %v", err)
	}
	if _, err := database.DBThinkRoot.Exec(
		"INSERT INTO github_repositories (url, text, posted, date_added) VALUES (?, ?, 0, '2026-08-16T16:00:00Z')",
		"https://github.com/open-webui/open-webui", "===(en)Queued head===(uk)Голова черги===",
	); err != nil {
		t.Fatalf("failed to insert queued repository: %v", err)
	}

	t.Run("by url with language", func(t *testing.T) {
		rec, data := postGetRepository(t, `{"url":"`+wantURL+`","text_language":"en"}`)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
		}
		if len(data.Items) != 1 {
			t.Fatalf("got %d items, want exactly the requested one", len(data.Items))
		}
		if data.Items[0].URL != wantURL {
			t.Errorf("url = %q, want %q", data.Items[0].URL, wantURL)
		}
		if data.Items[0].Text != "English description" {
			t.Errorf("text = %q, want the English text", data.Items[0].Text)
		}
		if !data.Items[0].Posted {
			t.Error("posted = false, want an already published repository to be returned as posted")
		}
		if data.TotalItems != 1 || data.Page != 1 || data.PageSize != 1 {
			t.Errorf("pagination = {items:%d page:%d size:%d}, want a single-item page", data.TotalItems, data.Page, data.PageSize)
		}
		if data.All != 2 || data.Posted != 1 || data.Unposted != 1 {
			t.Errorf("counts = {all:%d posted:%d unposted:%d}, want the global counts", data.All, data.Posted, data.Unposted)
		}
	})

	t.Run("by id falls back to available language", func(t *testing.T) {
		rec, data := postGetRepository(t, `{"id":1,"text_language":"de"}`)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
		}
		if len(data.Items) != 1 || data.Items[0].URL != wantURL {
			t.Fatalf("items = %+v, want the repository with id 1", data.Items)
		}
		if data.Items[0].Text != "Український опис" {
			t.Errorf("text = %q, want the Ukrainian fallback", data.Items[0].Text)
		}
	})

	t.Run("without language returns the raw text", func(t *testing.T) {
		_, data := postGetRepository(t, `{"url":"`+wantURL+`"}`)

		if len(data.Items) != 1 || data.Items[0].Text != multilingual {
			t.Fatalf("text = %q, want the stored multilingual text", data.Items[0].Text)
		}
	})

	t.Run("unknown url", func(t *testing.T) {
		rec, _ := postGetRepository(t, `{"url":"https://github.com/nope/nope"}`)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404 (body: %s)", rec.Code, rec.Body.String())
		}
	})

	t.Run("multiline text keeps the requested language", func(t *testing.T) {
		const multilineURL = "https://github.com/multiline/repo"
		if _, err := database.DBThinkRoot.Exec(
			"INSERT INTO github_repositories (url, text, posted, date_added) VALUES (?, ?, 1, '2026-08-15T16:00:00Z')",
			multilineURL, "===(en)Line one\nLine two===(uk)Рядок один\nРядок два===",
		); err != nil {
			t.Fatalf("failed to insert multiline repository: %v", err)
		}

		_, data := postGetRepository(t, `{"url":"`+multilineURL+`","text_language":"en"}`)

		if len(data.Items) != 1 || data.Items[0].Text != "Line one\nLine two" {
			t.Fatalf("text = %q, want the multi-line English text", data.Items[0].Text)
		}
	})

	// The single-item branch must be entered only on an explicit identifier:
	// omitted or null fields have to keep serving the queue.
	t.Run("omitted and null identifiers stay in queue mode", func(t *testing.T) {
		for _, body := range []string{`{}`, `{"id":null,"url":null}`} {
			rec, data := postGetRepository(t, body)

			if rec.Code != http.StatusOK {
				t.Fatalf("body %s: status = %d, want 200 (%s)", body, rec.Code, rec.Body.String())
			}
			if len(data.Items) != data.All {
				t.Errorf("body %s: got %d items, want all %d", body, len(data.Items), data.All)
			}
		}

		// A null identifier alongside real filters must not hijack the query.
		_, data := postGetRepository(t, `{"limit":1,"posted":false,"sort_by":"publication_queue","sort_order":"ASC","id":null}`)
		if len(data.Items) != 1 || data.Items[0].Posted {
			t.Fatalf("items = %+v, want the unposted head of the queue", data.Items)
		}
	})

	t.Run("queue mode is unaffected", func(t *testing.T) {
		rec, data := postGetRepository(t, `{"limit":1,"posted":false,"sort_by":"publication_queue","sort_order":"ASC","text_language":"en"}`)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
		}
		if len(data.Items) != 1 || data.Items[0].Text != "Queued head" {
			t.Fatalf("items = %+v, want the unposted head of the queue", data.Items)
		}
	})
}

// The identifier checks run before any database access, so they are safe to
// exercise without a live connection.
func TestGetRepositoryRejectsInvalidIdentifiers(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "id and url together", body: `{"id":1,"url":"https://github.com/foo/bar"}`},
		{name: "non-positive id", body: `{"id":0}`},
		{name: "blank url", body: `{"url":"   "}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/think-root/api/get-repository/", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			GetRepository(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, http.StatusBadRequest, rec.Body.String())
			}
		})
	}
}

func TestParseMultilingualTextFallsBackWhenRequestedLanguageMissing(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		language string
		want     string
	}{
		{
			name:     "plain text is returned for any requested language",
			text:     "Plain repository description",
			language: "uk",
			want:     "Plain repository description",
		},
		{
			name:     "single language text falls back to its available content",
			text:     "(en)English repository description",
			language: "uk",
			want:     "English repository description",
		},
		{
			name:     "multilingual text prefers requested language when available",
			text:     "===(en)English repository description===(uk)Український опис===",
			language: "en",
			want:     "English repository description",
		},
		{
			name:     "multilingual text falls back to Ukrainian when requested language is missing",
			text:     "===(en)English repository description===(uk)Український опис===",
			language: "pl",
			want:     "Український опис",
		},
		{
			name:     "multilingual text falls back to first available language without Ukrainian",
			text:     "===(en)English repository description===(de)Deutsche Beschreibung===",
			language: "uk",
			want:     "English repository description",
		},
		{
			// update-repository-text stores whatever it is given, newlines
			// included. A segment spanning lines used to fail the match, so its
			// language was dropped and the caller silently received another one.
			name:     "multiline segment is returned for its own language",
			text:     "===(en)Line one\nLine two===(uk)Рядок один\nРядок два===",
			language: "en",
			want:     "Line one\nLine two",
		},
		{
			name:     "multiline single language text is returned in full",
			text:     "(en)Line one\nLine two",
			language: "en",
			want:     "Line one\nLine two",
		},
		{
			name:     "multiline text falls back to Ukrainian when the language is missing",
			text:     "===(en)Line one\nLine two===(uk)Рядок один\nРядок два===",
			language: "pl",
			want:     "Рядок один\nРядок два",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMultilingualText(tt.text, tt.language)
			if err != nil {
				t.Fatalf("expected fallback text without error, got %v", err)
			}
			if got != tt.want {
				t.Fatalf("unexpected text: got %q, want %q", got, tt.want)
			}
		})
	}
}
