package models

import(
	"database/sql"
)

type Experiment struct {
	Id int
	Prompt string
	Options []Option
	Participants []Participant
	Organiser Participant
}

type ExperimentModel struct {
	DB *sql.DB
}

// Returns the winning Option and the number of votes it received
func (e Experiment) winner() (Option, int) {
	max := 0
	var winner Option

	for i := 0; i < len(e.Options); i++ {
		if e.Options[i].Votes > max {
			max = e.Options[i].Votes
			winner = e.Options[i]
		}
	}

	return winner, max
}


// Returns true if the Experiment is waiting on 1 or more Participants to vote
func (e Experiment) isOpen() (bool) {
	totalVotes := 0

	for i := 0; i < len(e.Options); i++ {
		totalVotes += e.Options[i].Votes
	}

	if totalVotes < len(e.Participants) {
		return true
	}

	return false
}

func (e ExperimentModel) All() ([]Experiment, error) {
	rows, err := e.DB.Query("SELECT * FROM Experiment")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var experiments []Experiment

	for rows.Next() {
		var experiment Experiment

		err := rows.Scan(&experiment.Id, &experiment.Prompt)
		if err != nil {
			return nil, err
		}

		experiments = append(experiments, experiment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return experiments, nil
}


func (e ExperimentModel) Setup() error {
	_, err := e.DB.Exec(`
		DROP TABLE IF EXISTS Experiment;
		CREATE TABLE Experiment (
			id				INTEGER PRIMARY KEY AUTOINCREMENT,
			prompt			VARCHAR(128) NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	_, err = e.DB.Exec(`
		INSERT INTO Experiment (prompt) VALUES('Where do you want to go to dinner?')
	`)
	if err != nil {
		return err
	}

	return nil

}