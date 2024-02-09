package models

import "database/sql"

type Model interface{}

type dbModel struct {
	DB *sql.DB
}
