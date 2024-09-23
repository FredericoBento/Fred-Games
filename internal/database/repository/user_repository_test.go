package repository

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/FredericoBento/HandGame/internal/database/sqlite"
	"github.com/FredericoBento/HandGame/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

var (
	testDB *sql.DB
)

func TestMain(m *testing.M) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	testDB = db
	err = sqlite.CreateTables(testDB)
	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}

func TestGetAll(t *testing.T) {

	t.Run("ReturnZeroUsers", func(t *testing.T) {
		repo := NewSQLiteUserRepository(testDB)

		users, err := repo.GetAll(context.TODO())
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}

		if users != nil {
			t.Errorf("expected users to be nil but got: %v", users)
		}

	})

	t.Run("ReturnUsers", func(t *testing.T) {
		testUser1 := models.User{
			Username: "randomdata",
			Password: "randomdata",
		}

		testUser2 := models.User{
			Username: "randomdata2",
			Password: "randomdata2",
		}

		_, err := testDB.Exec("INSERT INTO users(username, password) VALUES(?, ?)", testUser1.Username, testUser1.Password)
		if err != nil {
			log.Fatal(err)
		}

		_, err = testDB.Exec("INSERT INTO users(username, password) VALUES(?, ?)", testUser2.Username, testUser2.Password)
		if err != nil {
			log.Fatal(err)
		}

		repo := NewSQLiteUserRepository(testDB)

		users, err := repo.GetAll(context.TODO())
		if err != nil {
			t.Errorf(err.Error())
		}

		if len(users) < 2 || len(users) > 2 {
			t.Errorf("expected 2 users got: %d", len(users))
		}

		for i := range len(users) {
			if users[i].Username != testUser1.Username && users[i].Username != testUser2.Username {
				t.Errorf("expected the same username")
			} else {
				if users[i].Username == testUser1.Username {
					if users[i].Password != testUser1.Password {
						t.Errorf("expected the same password")
					}
				} else {
					if users[i].Username == testUser2.Username {
						if users[i].Password != testUser2.Password {
							t.Errorf("expected the same password")
						}
					}
				}
			}
		}
	})
}
