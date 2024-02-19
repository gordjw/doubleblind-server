package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	Id            string
	OauthId       string `json:"id"`
	OauthProvider string
	Name          string
	Email         string `json:"email"`
}

type UserModel struct {
	DB *sql.DB
}

func NewUser(id, oauthId, oauthProvider, name, email string) User {
	return User{
		Id:            id,
		OauthId:       oauthId,
		OauthProvider: "google",
		Name:          name,
		Email:         email,
	}
}

func (u UserModel) GetOrCreateUser(user User) error {
	_, err := u.DB.Exec(
		`INSERT INTO User (oauth_id, oauth_provider, email)
		VALUES(?, ?, ?)
		`,
		user.OauthId, user.OauthProvider, user.Email,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
