package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/randnull/banner-service/internal/config"
	"github.com/randnull/banner-service/internal/errors"
	"github.com/randnull/banner-service/pkg/models"
)


type UserRepository struct {
	db *sqlx.DB
}


func NewUserRepository(cfg *config.Config) *UserRepository {
	link := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.UserDB, cfg.PasswordDB, cfg.HostDB, cfg.PortDB, cfg.NameDB)

	db, err := sqlx.Open("postgres", link)

	if err != nil {
		log.Fatal(err)
	}

	err = db.PingContext(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	return &UserRepository{
		db: db,
	}
}


func (storage *UserRepository) AddUser(register_form *models.Register, is_admin bool) error {
	query := `SELECT COUNT(*) FROM users WHERE username = $1`

	var count int

	err := storage.db.Get(&count, query, register_form.Username)

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.UsernameAlreadyTaken
	}

	query = `INSERT INTO users (username, password, is_admin) VALUES ($1, $2, $3)`
	
	_, err = storage.db.Exec(query,
							  register_form.Username,
							  register_form.Password,
							  is_admin)

	if err != nil {
		return err
	}

	return nil
}


func (storage *UserRepository) GetUser(username string, password string) (*models.User, error) {
	query := `SELECT id, username, password, is_admin FROM users WHERE username = $1 AND password = $2`

	var user models.User
	
	err := storage.db.QueryRow(query, username, password).Scan(&user.ID, &user.Username, &user.Password, &user.IsAdmin)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.UserNotFound
		}
		return nil, err
	}

	return &user, nil
}
