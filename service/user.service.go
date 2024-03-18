package service

import (
	"database/sql"

	"github.com/crowmw/risiti/model"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Create(user model.User) error
	CheckEmail(email string) (bool, error)
}

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

func (s *UserService) Create(user model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO user(email, password) VALUES($1, $2)`

	_, err = s.DB.Exec(stmt, user.Email, hashedPassword)

	return err
}

func (s *UserService) CheckEmail(email string) (bool, error) {
	query := `SELECT * FROM user WHERE email=$1`

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return false, nil
	}

	defer stmt.Close()

	user := model.User{}
	err = stmt.QueryRow(email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}

		return false, nil
	}

	return true, nil
}
