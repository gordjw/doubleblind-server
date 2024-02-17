package models

import (
	"database/sql"
)

/**
 * Vote models a many-to-many relationship between participants and the options that they're voting for
 */

type Vote struct {
	OptionId     string
	ParticpantId string
}

type VoteModel struct {
	DB *sql.DB
}

func (v VoteModel) VoteFor(experimentId, optionId, participantId int) error {
	v.DB.Exec(`INSERT INTO Vote 
		(experiment_id, option_id, participant_id)
		VALUES (?, ?, ?)`,
		experimentId,
		optionId,
		participantId,
	)
	return nil
}
