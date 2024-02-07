package models

import (
	"database/sql"
)

type Vote struct {
	ExperimentId string
	OptionId     string
	ParticpantId string
}

type VoteModel struct {
	DB *sql.DB
}

func (v VoteModel) VoteFor(exp_id int, opt_id int, par_id int) error {

	return nil
}

func (v VoteModel) AttachedToExperiment(exp_id string) ([]Vote, error) {
	var votes []Vote

	return votes, nil
}
