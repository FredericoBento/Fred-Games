package sqlite

import "database/sql"

func CreateTables(db *sql.DB) error {
	var err error

	if err = createUserTable(db); err != nil {
		return err
	}

	return nil
}

func createUserTable(db *sql.DB) error {
	query := `
	    CREATE TABLE IF NOT EXISTS users (
	        id INTEGER PRIMARY KEY AUTOINCREMENT,
	        username TEXT NOT NULL,
	        password TEXT NOT NULL
	    );`

	_, err := db.Exec(query)

	return err

}
