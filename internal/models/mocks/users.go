package mocks

import (
	"github.com/manny-e1/snippetbox/internal/models"
	"time"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dagim@gmail.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "manny@gmail.com" && password == "12345678" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}
func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (um *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:

		return &models.User{
			ID:      1,
			Name:    "Amanuel",
			Email:   "manny@gmail.com",
			Created: time.Date(2022, 12, 20, 9, 23, 0, 0, time.UTC),
		}, nil
	default:
		return nil, models.ErrNoRecord
	}
}
