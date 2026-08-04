package database

import (
	"database/sql"
	"strconv"
	"testing"
	"time"
)

func setupArchiveTestDB(t *testing.T) {
	t.Helper()

	originalDB := DBThinkRoot
	// Must be the app driver: GetArchivedRepositories relies on unicode_lower(),
	// which only exists on connections opened through it.
	db, err := sql.Open(SQLiteDriverName, t.TempDir()+"/content-alchemist-test.db")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	DBThinkRoot = db
	t.Cleanup(func() {
		DBThinkRoot.Close()
		DBThinkRoot = originalDB
	})

	if err := createTableIfNotExists(); err != nil {
		t.Fatalf("failed to create test schema: %v", err)
	}
	if err := ensurePublicationQueueSchema(); err != nil {
		t.Fatalf("failed to ensure publication queue schema: %v", err)
	}
	if err := ensureArchiveSchema(); err != nil {
		t.Fatalf("failed to ensure archive schema: %v", err)
	}
}

// insertPostedTestRepository inserts a published repository with an explicit
// publication date. dateAdded/datePosted are RFC3339 strings.
func insertPostedTestRepository(t *testing.T, url, dateAdded, datePosted string) int64 {
	t.Helper()

	parsedAdded, err := time.Parse(time.RFC3339, dateAdded)
	if err != nil {
		t.Fatalf("failed to parse date_added: %v", err)
	}
	parsedPosted, err := time.Parse(time.RFC3339, datePosted)
	if err != nil {
		t.Fatalf("failed to parse date_posted: %v", err)
	}

	result, err := DBThinkRoot.Exec(
		"INSERT INTO github_repositories (url, text, posted, date_added, date_posted) VALUES (?, ?, 1, ?, ?)",
		url,
		"text for "+url,
		parsedAdded,
		parsedPosted,
	)
	if err != nil {
		t.Fatalf("failed to insert test repository: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("failed to read inserted id: %v", err)
	}

	return id
}

func countMainRepositories(t *testing.T) int {
	t.Helper()

	var count int
	if err := DBThinkRoot.QueryRow("SELECT COUNT(*) FROM github_repositories").Scan(&count); err != nil {
		t.Fatalf("failed to count repositories: %v", err)
	}
	return count
}

func TestArchiveRepositoriesMovesPublishedRepository(t *testing.T) {
	setupArchiveTestDB(t)

	id := insertPostedTestRepository(t, "https://github.com/a/b", "2026-01-01T10:00:00Z", "2026-01-05T10:00:00Z")

	archived, failures, err := ArchiveRepositories([]string{"1"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(failures) != 0 {
		t.Fatalf("expected no failures, got %+v", failures)
	}
	if len(archived) != 1 {
		t.Fatalf("expected 1 archived repository, got %d", len(archived))
	}

	got := archived[0]
	if got.URL != "https://github.com/a/b" {
		t.Errorf("got url %q, want %q", got.URL, "https://github.com/a/b")
	}
	if got.OriginalID == nil || *got.OriginalID != id {
		t.Errorf("got original_id %v, want %d", got.OriginalID, id)
	}
	if got.DateArchived.IsZero() {
		t.Error("date_archived was not set")
	}
	if got.DatePosted == nil {
		t.Error("date_posted was not carried over")
	}

	if count := countMainRepositories(t); count != 0 {
		t.Errorf("repository was not removed from the main table, %d rows left", count)
	}

	total, err := CountArchivedRepositories()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Errorf("got %d archived rows, want 1", total)
	}
}

func TestArchiveRepositoriesRejectsUnpublished(t *testing.T) {
	setupArchiveTestDB(t)

	insertTestRepository(t, "https://github.com/a/unposted", "2026-01-01T10:00:00Z", 0)

	archived, failures, err := ArchiveRepositories([]string{"https://github.com/a/unposted"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(archived) != 0 {
		t.Fatalf("expected nothing to be archived, got %d", len(archived))
	}
	if len(failures) != 1 || failures[0].Reason != ArchiveReasonNotPosted {
		t.Fatalf("expected a %s failure, got %+v", ArchiveReasonNotPosted, failures)
	}

	if count := countMainRepositories(t); count != 1 {
		t.Errorf("unpublished repository must stay in the main table, got %d rows", count)
	}
}

func TestArchiveRepositoriesReportsMissingAndDuplicateIdentifiers(t *testing.T) {
	setupArchiveTestDB(t)

	insertPostedTestRepository(t, "https://github.com/a/b", "2026-01-01T10:00:00Z", "2026-01-05T10:00:00Z")

	archived, failures, err := ArchiveRepositories([]string{"1", "1", "999"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(archived) != 1 {
		t.Fatalf("expected 1 archived repository, got %d", len(archived))
	}
	if len(failures) != 2 {
		t.Fatalf("expected 2 failures, got %+v", failures)
	}
	if failures[0].Reason != ArchiveReasonAlreadyProcessed {
		t.Errorf("got reason %q, want %q", failures[0].Reason, ArchiveReasonAlreadyProcessed)
	}
	if failures[1].Reason != ArchiveReasonNotFound {
		t.Errorf("got reason %q, want %q", failures[1].Reason, ArchiveReasonNotFound)
	}
}

func TestArchiveAllowsDuplicateURLs(t *testing.T) {
	setupArchiveTestDB(t)

	const url = "https://github.com/a/b"

	for i := 0; i < 2; i++ {
		insertPostedTestRepository(t, url, "2026-01-01T10:00:00Z", "2026-01-05T10:00:00Z")
		if _, failures, err := ArchiveRepositories([]string{url}, false); err != nil || len(failures) != 0 {
			t.Fatalf("archiving round %d failed: err=%v failures=%+v", i, err, failures)
		}
	}

	total, err := CountArchivedRepositories()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Errorf("got %d archived rows, want 2 (the archive must allow duplicate urls)", total)
	}
}

func TestArchivedRepositoryCanBeCollectedAgain(t *testing.T) {
	setupArchiveTestDB(t)

	const url = "https://github.com/a/b"
	insertPostedTestRepository(t, url, "2026-01-01T10:00:00Z", "2026-01-05T10:00:00Z")

	if _, _, err := ArchiveRepositories([]string{url}, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	exists, err := SearchPostInDB(url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Fatal("an archived url must not be treated as an existing repository")
	}

	if err := AddRepositoryToDB(url, "fresh description"); err != nil {
		t.Fatalf("failed to re-add archived repository: %v", err)
	}
	if count := countMainRepositories(t); count != 1 {
		t.Errorf("expected the repository to be re-added, got %d rows", count)
	}
}

func TestArchiveRepositoriesOlderThanUsesDatePosted(t *testing.T) {
	setupArchiveTestDB(t)

	now := time.Now().UTC()
	oldPost := now.AddDate(0, 0, -40).Format(time.RFC3339)
	recentPost := now.AddDate(0, 0, -5).Format(time.RFC3339)
	oldAdded := now.AddDate(0, 0, -100).Format(time.RFC3339)

	insertPostedTestRepository(t, "https://github.com/a/old", oldAdded, oldPost)
	insertPostedTestRepository(t, "https://github.com/a/recent", oldAdded, recentPost)
	// Added long ago but never published: must be left alone.
	insertTestRepository(t, "https://github.com/a/unposted", oldAdded, 0)

	archived, err := ArchiveRepositoriesOlderThan(30, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(archived) != 1 {
		t.Fatalf("expected 1 archived repository, got %d", len(archived))
	}
	if archived[0].URL != "https://github.com/a/old" {
		t.Errorf("got url %q, want %q", archived[0].URL, "https://github.com/a/old")
	}
	if archived[0].ID == 0 || archived[0].DateArchived.IsZero() {
		t.Error("archived row is missing its id or date_archived")
	}

	if count := countMainRepositories(t); count != 2 {
		t.Errorf("got %d rows left in the main table, want 2", count)
	}
}

func TestArchiveRepositoriesOlderThanDryRunChangesNothing(t *testing.T) {
	setupArchiveTestDB(t)

	oldPost := time.Now().UTC().AddDate(0, 0, -40).Format(time.RFC3339)
	insertPostedTestRepository(t, "https://github.com/a/old", "2026-01-01T10:00:00Z", oldPost)

	preview, err := ArchiveRepositoriesOlderThan(30, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(preview) != 1 {
		t.Fatalf("expected 1 previewed repository, got %d", len(preview))
	}
	if !preview[0].DateArchived.IsZero() {
		t.Error("a dry run must not set date_archived")
	}

	if count := countMainRepositories(t); count != 1 {
		t.Errorf("a dry run must not remove rows, got %d rows", count)
	}
	total, err := CountArchivedRepositories()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 0 {
		t.Errorf("a dry run must not write to the archive, got %d rows", total)
	}
}

func TestArchiveRepositoriesOlderThanRejectsNonPositiveDays(t *testing.T) {
	setupArchiveTestDB(t)

	if _, err := ArchiveRepositoriesOlderThan(0, false); err == nil {
		t.Fatal("expected an error for days = 0")
	}
}

// seedArchive inserts archive rows directly so the read path can be tested
// without going through the move logic.
func seedArchive(t *testing.T, url, text, dateAdded, datePosted, dateArchived string) {
	t.Helper()

	parse := func(value string) time.Time {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			t.Fatalf("failed to parse %q: %v", value, err)
		}
		return parsed
	}

	_, err := DBThinkRoot.Exec(
		"INSERT INTO archived_repositories (original_id, url, text, date_added, date_posted, date_archived) VALUES (?, ?, ?, ?, ?, ?)",
		nil, url, text, parse(dateAdded), parse(datePosted), parse(dateArchived),
	)
	if err != nil {
		t.Fatalf("failed to seed archive: %v", err)
	}
}

func TestGetArchivedRepositoriesFiltersAndPaginates(t *testing.T) {
	setupArchiveTestDB(t)

	seedArchive(t, "https://github.com/alpha/one", "alpha tool", "2026-01-01T10:00:00Z", "2026-01-05T10:00:00Z", "2026-02-01T10:00:00Z")
	seedArchive(t, "https://github.com/beta/two", "beta tool", "2026-02-01T10:00:00Z", "2026-02-10T10:00:00Z", "2026-03-01T10:00:00Z")
	seedArchive(t, "https://github.com/gamma/three", "gamma helper", "2026-03-01T10:00:00Z", "2026-03-15T10:00:00Z", "2026-04-01T10:00:00Z")

	mustParse := func(value string) *time.Time {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			t.Fatalf("failed to parse %q: %v", value, err)
		}
		return &parsed
	}

	tests := []struct {
		name      string
		filter    ArchivedFilter
		wantTotal int
		wantURLs  []string
	}{
		{
			name:      "no filter, newest archived first",
			filter:    ArchivedFilter{},
			wantTotal: 3,
			wantURLs:  []string{"https://github.com/gamma/three", "https://github.com/beta/two", "https://github.com/alpha/one"},
		},
		{
			name:      "url substring",
			filter:    ArchivedFilter{URL: "beta"},
			wantTotal: 1,
			wantURLs:  []string{"https://github.com/beta/two"},
		},
		{
			name:      "text substring",
			filter:    ArchivedFilter{Text: "helper"},
			wantTotal: 1,
			wantURLs:  []string{"https://github.com/gamma/three"},
		},
		{
			name:      "date_archived range",
			filter:    ArchivedFilter{DateArchivedFrom: mustParse("2026-02-15T00:00:00Z"), DateArchivedTo: mustParse("2026-03-15T00:00:00Z")},
			wantTotal: 1,
			wantURLs:  []string{"https://github.com/beta/two"},
		},
		{
			name:      "date_posted range",
			filter:    ArchivedFilter{DatePostedFrom: mustParse("2026-02-01T00:00:00Z")},
			wantTotal: 2,
			wantURLs:  []string{"https://github.com/gamma/three", "https://github.com/beta/two"},
		},
		{
			name:      "date_added range",
			filter:    ArchivedFilter{DateAddedTo: mustParse("2026-02-02T00:00:00Z")},
			wantTotal: 2,
			wantURLs:  []string{"https://github.com/beta/two", "https://github.com/alpha/one"},
		},
		{
			name:      "sorted by date_posted ascending",
			filter:    ArchivedFilter{SortBy: "date_posted", SortOrder: "asc"},
			wantTotal: 3,
			wantURLs:  []string{"https://github.com/alpha/one", "https://github.com/beta/two", "https://github.com/gamma/three"},
		},
		{
			name:      "second page",
			filter:    ArchivedFilter{Limit: 2, Offset: 2},
			wantTotal: 3,
			wantURLs:  []string{"https://github.com/alpha/one"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, total, err := GetArchivedRepositories(tt.filter)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if total != tt.wantTotal {
				t.Errorf("got total %d, want %d", total, tt.wantTotal)
			}
			if len(got) != len(tt.wantURLs) {
				t.Fatalf("got %d items, want %d", len(got), len(tt.wantURLs))
			}
			for i, want := range tt.wantURLs {
				if got[i].URL != want {
					t.Errorf("item %d: got url %q, want %q", i, got[i].URL, want)
				}
			}
		})
	}
}

func TestGetArchivedRepositoriesKeepsNullDatesLast(t *testing.T) {
	setupArchiveTestDB(t)

	seedArchive(t, "https://github.com/a/dated", "dated", "2026-01-01T10:00:00Z", "2026-01-05T10:00:00Z", "2026-02-01T10:00:00Z")

	// A row archived without a publication date — possible for a repository that
	// was marked posted without date_posted ever being set.
	dateArchived, err := time.Parse(time.RFC3339, "2026-02-02T10:00:00Z")
	if err != nil {
		t.Fatalf("failed to parse date: %v", err)
	}
	if _, err := DBThinkRoot.Exec(
		"INSERT INTO archived_repositories (original_id, url, text, date_added, date_posted, date_archived) VALUES (NULL, ?, ?, NULL, NULL, ?)",
		"https://github.com/a/undated", "undated", dateArchived,
	); err != nil {
		t.Fatalf("failed to seed archive: %v", err)
	}

	for _, order := range []string{"asc", "desc"} {
		t.Run(order, func(t *testing.T) {
			got, _, err := GetArchivedRepositories(ArchivedFilter{SortBy: "date_posted", SortOrder: order})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != 2 {
				t.Fatalf("got %d items, want 2", len(got))
			}
			if got[len(got)-1].URL != "https://github.com/a/undated" {
				t.Errorf("rows without date_posted must sort last, got %q last", got[len(got)-1].URL)
			}
		})
	}
}

func TestGetArchivedRepositoriesNegativeLimitDoesNotDisablePaging(t *testing.T) {
	setupArchiveTestDB(t)

	for i := 0; i < 3; i++ {
		seedArchive(t, "https://github.com/a/"+strconv.Itoa(i), "text", "2026-01-01T10:00:00Z", "2026-01-05T10:00:00Z", "2026-02-01T10:00:00Z")
	}

	// The handler rejects a negative limit; this pins the storage-layer guard so
	// a bypass cannot silently return the whole archive.
	got, total, err := GetArchivedRepositories(ArchivedFilter{Limit: -1, Offset: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 3 {
		t.Errorf("got total %d, want 3", total)
	}
	if len(got) != 2 {
		t.Errorf("got %d items, want 2 (offset 1 of 3)", len(got))
	}
}

func TestGetArchivedRepositoriesSearchIsCaseInsensitiveBeyondASCII(t *testing.T) {
	setupArchiveTestDB(t)

	seedArchive(t, "https://github.com/Alpha/Tool", "Старий Опис Проєкту", "2026-01-01T10:00:00Z", "2026-01-05T10:00:00Z", "2026-02-01T10:00:00Z")

	tests := []struct {
		name   string
		filter ArchivedFilter
	}{
		{"url as stored", ArchivedFilter{URL: "Alpha"}},
		{"url lowercased", ArchivedFilter{URL: "alpha"}},
		{"url uppercased", ArchivedFilter{URL: "ALPHA"}},
		{"cyrillic as stored", ArchivedFilter{Text: "Старий"}},
		{"cyrillic lowercased", ArchivedFilter{Text: "старий"}},
		{"cyrillic uppercased", ArchivedFilter{Text: "СТАРИЙ"}},
		{"cyrillic with ї", ArchivedFilter{Text: "проєкту"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, total, err := GetArchivedRepositories(tt.filter)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if total != 1 {
				t.Errorf("got %d matches, want 1", total)
			}
		})
	}
}

func TestGetArchivedRepositoriesTreatsSearchWildcardsLiterally(t *testing.T) {
	setupArchiveTestDB(t)

	seedArchive(t, "https://github.com/a/b", "plain", "2026-01-01T10:00:00Z", "2026-01-05T10:00:00Z", "2026-02-01T10:00:00Z")

	got, total, err := GetArchivedRepositories(ArchivedFilter{URL: "%"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 0 || len(got) != 0 {
		t.Errorf("a literal %% must not match everything, got %d rows", total)
	}
}
