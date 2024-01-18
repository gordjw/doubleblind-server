package models

import(
	"database/sql"
)

type Option struct {
	value string
	votes int
}

type OptionModel struct {
	DB *sql.DB
}