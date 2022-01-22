package manager

import "database/sql"

type Folder struct {
	ID       int           `db:"id"`
	ParentID sql.NullInt32 `db:"parent"`
	Name     string        `db:"name"`
	Priority int           `db:"priority"`
	FullPath string        `db:"full_path"`
}
