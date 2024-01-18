package models

import(
	"database/sql"
)

type Participant struct {
	id int
	name string
	email string
}

type ParticipantModel struct {
	DB *sql.DB
}