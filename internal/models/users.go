package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (um *UserModel) Insert(name, email, password string) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name,email,password,created) 
			VALUES (?,?,?,UTC_TIMESTAMP())`
	_, err = um.DB.Exec(stmt, name, email, hashedPwd)
	if err != nil {
		var mySQLErr *mysql.MySQLError
		if errors.As(err, &mySQLErr) {
			if mySQLErr.Number == 1862 && strings.Contains(mySQLErr.Message, "users_us_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}
func (um *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}
func (um *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
