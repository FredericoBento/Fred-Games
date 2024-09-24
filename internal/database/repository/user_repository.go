package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/models"
)

var (
	ErrCouldNotStartTransaction = errors.New("could not start transaction")
	ErrCouldNotInsertUser       = errors.New("could not insert user")
	ErrCouldNotRollback         = errors.New("could not rollback transaction")
	ErrCouldNotCreateLogger     = errors.New("could not create logger for sqlite user repository")
	ErrCouldNotGetByUsername    = errors.New("could not get user by username")
)

type SQLiteUserRepository struct {
	DB  *sql.DB
	log *slog.Logger
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	lo, err := logger.NewRepositoryLogger("sqlite", "users", "", false)
	if err != nil {
		slog.Error(ErrCouldNotCreateLogger.Error() + " " + err.Error())
		lo = slog.Default()
	}
	return &SQLiteUserRepository{
		DB:  db,
		log: lo,
	}
}

func (r *SQLiteUserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	rows, err := r.DB.Query("SELECT * from users")
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()
	var users []models.User

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			r.log.Error(err.Error())
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *SQLiteUserRepository) Create(ctx context.Context, user *models.User) error {
	t, err := r.DB.Begin()
	if err != nil {
		r.log.Error(err.Error())
		return ErrCouldNotStartTransaction
	}
	query := "INSERT INTO users(username, password) VALUES(?, ?)"
	_, err = t.Exec(query, user.Username, user.Password)
	if err != nil {
		if err = t.Rollback(); err != nil {
			r.log.Error(err.Error())
			return ErrCouldNotRollback
		}
		r.log.Error(err.Error())
		return ErrCouldNotInsertUser
	}

	t.Commit()

	return nil
}

func (r *SQLiteUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := "SELECT * FROM users WHERE username = ?"
	rows, err := r.DB.Query(query, username)
	if err != nil {
		r.log.Error(err.Error())
		return nil, ErrCouldNotGetByUsername
	}
	defer rows.Close()
	var user models.User
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	err = rows.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	return &user, nil
}
