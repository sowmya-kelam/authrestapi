package models

import (
	"database/sql"
	"errors"

	"restapi/utils"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int64
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshtoken" binding:"required"`
}

func (u User) Save(db *sql.DB) error {
	if db == nil {
		return errors.New("database is nill")
	}

	query := `INSERT INTO users(email,password)
	 VALUES (?,?)`

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedpassword, err := utils.Hashpassword(u.Password)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.Email, hashedpassword)

	if err != nil {
		return err
	}

	// u.ID, err = result.LastInsertId()

	// if err != nil {
	// 	return err
	// }

	return nil
}

func (u *User) Validate(db *sql.DB) error {
	if db == nil {
		return errors.New("database is nill")
	}

	query := `select id,password from users where email=?`

	row := db.QueryRow(query, u.Email)
	var retrievedpassword string
	err := row.Scan(&u.ID, &retrievedpassword)
	if err != nil {
		return err
	}

	passwordvalid := utils.Checkpasswordhash(u.Password, retrievedpassword)

	if !passwordvalid {
		return errors.New("invalid credentials")
	}

	return nil
}


