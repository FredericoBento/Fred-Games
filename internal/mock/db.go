package mock

import (
	"database/sql"
	"github.com/FredericoBento/HandGame/internal/database/sqlite"
)

type MockSQLiteDB struct {
	DB *sql.DB
}

func NewMockSQLiteDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	err = sqlite.CreateTables(db)
	if err != nil {
		return nil, err
	}

	return db, nil

}
