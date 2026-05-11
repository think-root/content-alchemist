package database

import (
	"database/sql"
	"strconv"
	"testing"
	"time"
)

func setupPublicationQueueTestDB(t *testing.T) {
	t.Helper()

	originalDB := DBThinkRoot
	db, err := sql.Open("sqlite3", t.TempDir()+"/content-alchemist-test.db")
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
}

func insertTestRepository(t *testing.T, url string, dateAdded string, posted int) int64 {
	t.Helper()

	parsedDate, err := time.Parse(time.RFC3339, dateAdded)
	if err != nil {
		t.Fatalf("failed to parse date: %v", err)
	}

	result, err := DBThinkRoot.Exec(
		"INSERT INTO github_repositories (url, text, posted, date_added) VALUES (?, ?, ?, ?)",
		url,
		"text for "+url,
		posted,
		parsedDate,
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

func TestPublicationQueueSortingPromotedItemsFirst(t *testing.T) {
	setupPublicationQueueTestDB(t)

	oldID := insertTestRepository(t, "https://github.com/example/old", "2025-01-01T00:00:00Z", 0)
	middleID := insertTestRepository(t, "https://github.com/example/middle", "2025-01-02T00:00:00Z", 0)
	newID := insertTestRepository(t, "https://github.com/example/new", "2025-01-03T00:00:00Z", 0)

	if _, err := PromoteRepositoryToNextByIDOrURL(strconv.FormatInt(middleID, 10), true); err != nil {
		t.Fatalf("failed to promote middle repository: %v", err)
	}
	if _, err := PromoteRepositoryToNextByIDOrURL(strconv.FormatInt(oldID, 10), true); err != nil {
		t.Fatalf("failed to promote old repository: %v", err)
	}

	unposted := false
	repos, _, err := GetRepository(0, 0, &unposted, "publication_queue", "ASC")
	if err != nil {
		t.Fatalf("failed to get repositories: %v", err)
	}

	got := []int64{repos[0].ID, repos[1].ID, repos[2].ID}
	want := []int64{oldID, middleID, newID}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("queue order mismatch: got %v, want %v", got, want)
		}
	}
}

func TestUpdatePostedStatusClearsPublishPriority(t *testing.T) {
	setupPublicationQueueTestDB(t)

	id := insertTestRepository(t, "https://github.com/example/promoted", "2025-01-01T00:00:00Z", 0)
	repo, err := PromoteRepositoryToNextByIDOrURL(strconv.FormatInt(id, 10), true)
	if err != nil {
		t.Fatalf("failed to promote repository: %v", err)
	}
	if repo.ID != id || repo.PublishPriority == nil {
		t.Fatalf("repository was not promoted: %#v", repo)
	}

	if err := UpdatePostedStatusByURL("https://github.com/example/promoted", true); err != nil {
		t.Fatalf("failed to update posted status: %v", err)
	}

	updatedRepo, err := GetRepositoryByIDOrURL(strconv.FormatInt(id, 10), true)
	if err != nil {
		t.Fatalf("failed to fetch updated repository: %v", err)
	}
	if updatedRepo.Posted != 1 {
		t.Fatalf("expected repository to be posted, got %d", updatedRepo.Posted)
	}
	if updatedRepo.PublishPriority != nil {
		t.Fatalf("expected publish priority to be cleared, got %d", *updatedRepo.PublishPriority)
	}
}

func TestPromoteRepositoryRejectsPostedRepository(t *testing.T) {
	setupPublicationQueueTestDB(t)

	id := insertTestRepository(t, "https://github.com/example/posted", "2025-01-01T00:00:00Z", 1)

	if _, err := PromoteRepositoryToNextByIDOrURL(strconv.FormatInt(id, 10), true); err == nil {
		t.Fatal("expected error when promoting posted repository")
	}
}
