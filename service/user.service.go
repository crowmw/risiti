package service

import (
	"database/sql"

	"github.com/crowmw/risiti/model"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Create(user model.User) (model.User, error)
	Read(email string) (model.User, error)
	AnyExists() (bool, error)
}

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

func (s *UserService) Create(user model.User) (model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	query := `INSERT INTO user(email, password) VALUES($1, $2) RETURNING *`

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return model.User{}, err
	}

	defer stmt.Close()

	newUser := model.User{}
	err = stmt.QueryRow(user.Email, string(hashedPassword)).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.Password,
	)
	if err != nil {
		return model.User{}, err
	}

	return newUser, nil
}

func (s *UserService) Read(email string) (model.User, error) {
	query := `SELECT * FROM user WHERE email=$1`

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return model.User{}, nil
	}

	defer stmt.Close()

	user := model.User{}
	err = stmt.QueryRow(email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *UserService) AnyExists() (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM user) AS isEmpty`

	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return false, nil
	}

	defer stmt.Close()

	var isEmpty bool
	err = stmt.QueryRow().Scan(&isEmpty)
	if err != nil {
		return false, err
	}

	return isEmpty, nil
}
