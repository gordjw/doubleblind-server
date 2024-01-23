package models

import(
	// "log"
	"database/sql"
)

type Vote struct {
	ExperimentId int
	OptionId int
	ParticpantId int
}

type VoteModel struct {
	DB *sql.DB
}

func (v VoteModel) VoteFor(exp_id int, opt_id int, par_id int) error {
	// err := v.DB.Query(`INSERT INTO Vote `)

	return nil
}

func (v VoteModel) AttachedToExperiment(exp_id int) (*[]Vote, error) {
	var votes []Vote
	// err := v.DB.Query(`INSERT INTO Vote `)

	return &votes, nil
}

func (v VoteModel) Setup() error {
	_, err := v.DB.Exec(`
		DROP TABLE IF EXISTS Vote;
		CREATE TABLE Vote (
			id				INTEGER PRIMARY KEY AUTOINCREMENT,
			experiment_id	INTEGER NOT NULL,
			option_id		INTEGER NOT NULL,
			participant_id	INTEGER NOT NULL,
			FOREIGN KEY (experiment_id)
				REFERENCES Experiment(id)
				ON DELETE CASCADE,
			FOREIGN KEY (option_id)
				REFERENCES Option(id)
				ON DELETE CASCADE,
			FOREIGN KEY (participant_id)
				REFERENCES Participant(id)
				ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	_, err = v.DB.Exec(`
		INSERT INTO Vote (experiment_id, option_id, participant_id) VALUES ('1', '1', '1'), ('1', '1', '2'), ('1', '3', '3')
	`)
	if err != nil {
		return err
	}

	return nil
}