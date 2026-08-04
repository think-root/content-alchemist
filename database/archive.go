package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// ArchivedRepository is a published repository that was moved out of
// github_repositories into the append-only archive. Unlike the main table the
// url is NOT unique here: the same repository may be published and archived
// several times.
type ArchivedRepository struct {
	ID           int64      `json:"id"`
	OriginalID   *int64     `json:"original_id"`
	URL          string     `json:"url"`
	Text         string     `json:"text"`
	DateAdded    *time.Time `json:"date_added"`
	DatePosted   *time.Time `json:"date_posted"`
	DateArchived time.Time  `json:"date_archived"`
}

func (ArchivedRepository) TableName() string {
	return "archived_repositories"
}

// Reasons reported by ArchiveRepositories for identifiers it refused to archive.
const (
	ArchiveReasonNotFound         = "not_found"
	ArchiveReasonNotPosted        = "not_posted"
	ArchiveReasonAlreadyProcessed = "already_processed"
)

// ArchiveFailure describes a single identifier that could not be archived.
type ArchiveFailure struct {
	Identifier string `json:"identifier"`
	Reason     string `json:"reason"`
	Message    string `json:"message"`
}

// ArchivedFilter carries the search/pagination options for GetArchivedRepositories.
// All filters are combined with AND. The *To bounds are exclusive.
type ArchivedFilter struct {
	Limit  int
	Offset int

	URL  string
	Text string

	DateAddedFrom    *time.Time
	DateAddedTo      *time.Time
	DatePostedFrom   *time.Time
	DatePostedTo     *time.Time
	DateArchivedFrom *time.Time
	DateArchivedTo   *time.Time

	SortBy    string
	SortOrder string
}

const archivedColumns = "id, original_id, url, text, date_added, date_posted, date_archived"

// ArchiveRepositories moves the given repositories from github_repositories into
// archived_repositories. Only published repositories (posted = 1) can be
// archived; everything else is reported back as a failure without aborting the
// rest of the batch. The whole batch runs in a single transaction, which is
// rolled back only on an actual database error.
func ArchiveRepositories(identifiers []string, isID bool) ([]ArchivedRepository, []ArchiveFailure, error) {
	archived := make([]ArchivedRepository, 0, len(identifiers))
	failures := make([]ArchiveFailure, 0)

	tx, err := DBThinkRoot.Begin()
	if err != nil {
		return nil, nil, fmt.Errorf("error starting archive transaction: %v", err)
	}
	defer tx.Rollback()

	selectQuery := "SELECT id, url, text, posted, date_added, date_posted FROM github_repositories WHERE url = ?"
	deleteQuery := "DELETE FROM github_repositories WHERE url = ?"
	if isID {
		selectQuery = "SELECT id, url, text, posted, date_added, date_posted FROM github_repositories WHERE id = ?"
		deleteQuery = "DELETE FROM github_repositories WHERE id = ?"
	}

	seen := make(map[string]bool, len(identifiers))
	now := time.Now().UTC()

	for _, identifier := range identifiers {
		if seen[identifier] {
			failures = append(failures, ArchiveFailure{
				Identifier: identifier,
				Reason:     ArchiveReasonAlreadyProcessed,
				Message:    "duplicate identifier in the same request",
			})
			continue
		}
		seen[identifier] = true

		var (
			id         int64
			url        string
			text       string
			posted     int
			dateAdded  *time.Time
			datePosted *time.Time
		)
		err := tx.QueryRow(selectQuery, identifier).Scan(&id, &url, &text, &posted, &dateAdded, &datePosted)
		if err == sql.ErrNoRows {
			failures = append(failures, ArchiveFailure{
				Identifier: identifier,
				Reason:     ArchiveReasonNotFound,
				Message:    "repository not found",
			})
			continue
		}
		if err != nil {
			return nil, nil, fmt.Errorf("error fetching repository %s: %v", identifier, err)
		}

		if posted != 1 {
			failures = append(failures, ArchiveFailure{
				Identifier: identifier,
				Reason:     ArchiveReasonNotPosted,
				Message:    "only published repositories can be archived",
			})
			continue
		}

		result, err := tx.Exec(
			"INSERT INTO archived_repositories (original_id, url, text, date_added, date_posted, date_archived) VALUES (?, ?, ?, ?, ?, ?)",
			id, url, text, dateAdded, datePosted, now,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("error archiving repository %s: %v", identifier, err)
		}
		archiveID, err := result.LastInsertId()
		if err != nil {
			return nil, nil, fmt.Errorf("error reading archive id for repository %s: %v", identifier, err)
		}

		if _, err := tx.Exec(deleteQuery, identifier); err != nil {
			return nil, nil, fmt.Errorf("error removing archived repository %s: %v", identifier, err)
		}

		originalID := id
		archived = append(archived, ArchivedRepository{
			ID:           archiveID,
			OriginalID:   &originalID,
			URL:          url,
			Text:         text,
			DateAdded:    dateAdded,
			DatePosted:   datePosted,
			DateArchived: now,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("error committing archive transaction: %v", err)
	}

	return archived, failures, nil
}

// ArchiveRepositoriesOlderThan archives every published repository whose
// date_posted is older than the given number of days. The age is measured by
// publication date, not by the date the repository was added. When dryRun is
// true nothing is modified and the matching rows are returned as a preview.
func ArchiveRepositoriesOlderThan(days int, dryRun bool) ([]ArchivedRepository, error) {
	if days < 1 {
		return nil, fmt.Errorf("days must be a positive integer")
	}

	cutoff := time.Now().UTC().AddDate(0, 0, -days)

	tx, err := DBThinkRoot.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting archive transaction: %v", err)
	}
	defer tx.Rollback()

	rows, err := tx.Query(`
		SELECT id, url, text, date_added, date_posted
		FROM github_repositories
		WHERE posted = 1 AND date_posted IS NOT NULL AND julianday(date_posted) < julianday(?)
		ORDER BY date_posted ASC
	`, cutoff)
	if err != nil {
		return nil, fmt.Errorf("error selecting repositories to archive: %v", err)
	}

	candidates := make([]ArchivedRepository, 0)
	for rows.Next() {
		var repo ArchivedRepository
		var id int64
		if err := rows.Scan(&id, &repo.URL, &repo.Text, &repo.DateAdded, &repo.DatePosted); err != nil {
			rows.Close()
			return nil, fmt.Errorf("error scanning repository to archive: %v", err)
		}
		originalID := id
		repo.OriginalID = &originalID
		candidates = append(candidates, repo)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("error iterating repositories to archive: %v", err)
	}
	rows.Close()

	if dryRun || len(candidates) == 0 {
		return candidates, nil
	}

	now := time.Now().UTC()
	for i := range candidates {
		repo := &candidates[i]

		result, err := tx.Exec(
			"INSERT INTO archived_repositories (original_id, url, text, date_added, date_posted, date_archived) VALUES (?, ?, ?, ?, ?, ?)",
			*repo.OriginalID, repo.URL, repo.Text, repo.DateAdded, repo.DatePosted, now,
		)
		if err != nil {
			return nil, fmt.Errorf("error archiving repository %d: %v", *repo.OriginalID, err)
		}
		archiveID, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("error reading archive id for repository %d: %v", *repo.OriginalID, err)
		}

		if _, err := tx.Exec("DELETE FROM github_repositories WHERE id = ?", *repo.OriginalID); err != nil {
			return nil, fmt.Errorf("error removing archived repository %d: %v", *repo.OriginalID, err)
		}

		repo.ID = archiveID
		repo.DateArchived = now
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing archive transaction: %v", err)
	}

	return candidates, nil
}

// CountArchivedRepositories returns the total number of rows in the archive.
func CountArchivedRepositories() (int, error) {
	var count int
	if err := DBThinkRoot.QueryRow("SELECT COUNT(*) FROM archived_repositories").Scan(&count); err != nil {
		return 0, fmt.Errorf("error counting archived repositories: %v", err)
	}
	return count, nil
}

// GetArchivedRepositories returns a filtered, sorted and paginated slice of the
// archive together with the total number of rows matching the filter.
func GetArchivedRepositories(filter ArchivedFilter) ([]ArchivedRepository, int, error) {
	var conditions []string
	var args []interface{}

	// unicode_lower on both sides: SQLite's own LIKE only case-folds ASCII, so
	// without it a search for "старий" would miss a stored "Старий".
	if filter.URL != "" {
		conditions = append(conditions, `unicode_lower(url) LIKE unicode_lower(?) ESCAPE '\'`)
		args = append(args, "%"+escapeLikePattern(filter.URL)+"%")
	}
	if filter.Text != "" {
		conditions = append(conditions, `unicode_lower(text) LIKE unicode_lower(?) ESCAPE '\'`)
		args = append(args, "%"+escapeLikePattern(filter.Text)+"%")
	}

	dateFilters := []struct {
		column string
		from   *time.Time
		to     *time.Time
	}{
		{"date_added", filter.DateAddedFrom, filter.DateAddedTo},
		{"date_posted", filter.DatePostedFrom, filter.DatePostedTo},
		{"date_archived", filter.DateArchivedFrom, filter.DateArchivedTo},
	}
	for _, df := range dateFilters {
		if df.from != nil {
			conditions = append(conditions, fmt.Sprintf("julianday(%s) >= julianday(?)", df.column))
			args = append(args, df.from.UTC())
		}
		if df.to != nil {
			conditions = append(conditions, fmt.Sprintf("julianday(%s) < julianday(?)", df.column))
			args = append(args, df.to.UTC())
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var totalCount int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM archived_repositories %s", whereClause)
	if err := DBThinkRoot.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("error counting archived repositories: %v", err)
	}

	dataQuery := fmt.Sprintf(
		"SELECT %s FROM archived_repositories %s ORDER BY %s",
		archivedColumns, whereClause, archivedOrderBy(filter.SortBy, filter.SortOrder),
	)

	if filter.Limit > 0 || filter.Offset > 0 {
		limit := filter.Limit
		if limit <= 0 {
			limit = -1 // SQLite requires a LIMIT before OFFSET; -1 means "no limit"
		}
		dataQuery += " LIMIT ? OFFSET ?"
		args = append(args, limit, filter.Offset)
	}

	rows, err := DBThinkRoot.Query(dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching archived repositories: %v", err)
	}
	defer rows.Close()

	repositories := make([]ArchivedRepository, 0)
	for rows.Next() {
		var repo ArchivedRepository
		if err := rows.Scan(&repo.ID, &repo.OriginalID, &repo.URL, &repo.Text, &repo.DateAdded, &repo.DatePosted, &repo.DateArchived); err != nil {
			return nil, 0, fmt.Errorf("error scanning archived repository: %v", err)
		}
		repositories = append(repositories, repo)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating over archived repositories: %v", err)
	}

	return repositories, int(totalCount), nil
}

// IsValidArchivedSortBy reports whether the column may be used for sorting.
func IsValidArchivedSortBy(sortBy string) bool {
	switch sortBy {
	case "", "date_archived", "date_posted", "date_added", "id":
		return true
	}
	return false
}

// archivedOrderBy maps a whitelisted column/direction pair to an ORDER BY clause.
func archivedOrderBy(sortBy, sortOrder string) string {
	column := "date_archived"
	switch sortBy {
	case "date_posted", "date_added", "id":
		column = sortBy
	}

	direction := "DESC"
	if strings.EqualFold(sortOrder, "asc") {
		direction = "ASC"
	}

	if column == "id" {
		return "id " + direction
	}

	// Keep NULL dates at the end regardless of the direction, then break ties by
	// id. `col IS NULL` is 1 for NULL rows, so sorting it ASC puts them last.
	return fmt.Sprintf("%s IS NULL ASC, %s %s, id %s", column, column, direction, direction)
}

// escapeLikePattern escapes the LIKE wildcards so a user-supplied substring is
// matched literally. Must be paired with ESCAPE '\' in the query.
func escapeLikePattern(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(value)
}
