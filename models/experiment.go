package models

import (
	"database/sql"
	"fmt"
	"log"
)

type Experiment struct {
	Id           string
	Prompt       string
	Options      map[string]Option
	Participants map[string]Participant
	Organiser    Participant
}

type ExperimentModel struct {
	DB               *sql.DB
	OptionModel      *OptionModel
	ParticipantModel *ParticipantModel
}

func (e ExperimentModel) All() ([]Experiment, error) {
	const organiserId = 1

	rows, err := e.DB.Query(
		`SELECT id, prompt 
		FROM Experiment
		WHERE Experiment.organiserId = ?`,
		organiserId,
	)
	if err != nil {
		log.Println(err)
		return []Experiment{}, err
	}
	defer rows.Close()

	var experiments []Experiment

	for rows.Next() {
		var experiment Experiment

		err := rows.Scan(&experiment.Id, &experiment.Prompt)
		if err != nil {
			log.Println(err)
			return []Experiment{}, err
		}

		experiments = append(experiments, experiment)
	}
	if err = rows.Err(); err != nil {
		fmt.Printf(err.Error())
		return []Experiment{}, err
	}

	return experiments, nil
}

func (e ExperimentModel) One(id string) (Experiment, error) {
	rows, err := e.DB.Query(
		`SELECT Experiment.id, Experiment.prompt,
		Option.id as option_id, Option.value as option_label,
		count(Vote.option_id) as votes
		FROM Experiment
		INNER JOIN Option
		ON Experiment.id = Option.experiment_id
		LEFT JOIN Vote
		ON Option.id = Vote.option_id
		WHERE Experiment.id = ?
		GROUP BY Option.id, Vote.option_id`, id,
	)
	if err != nil {
		fmt.Println(err)
		return Experiment{}, err
	}
	defer rows.Close()

	var experiment Experiment
	experiment.Options = make(map[string]Option)

	for rows.Next() {
		var option Option

		err := rows.Scan(
			&experiment.Id, &experiment.Prompt,
			&option.Id, &option.Value, &option.Votes,
		)
		if err != nil {
			log.Println(err)
			return Experiment{}, err
		}

		experiment.Options[option.Id] = option
	}
	if err = rows.Err(); err != nil {
		return Experiment{}, err
	}

	return experiment, nil
}

func (e ExperimentModel) Add(prompt string, organiserId int, options map[string]Option) error {
	rows, err := e.DB.Query(`
		INSERT INTO Experiment
		(prompt, organiserId)
		VALUES (?, ?)
		RETURNING id
	`,
		prompt,
		organiserId,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rows.Close()

	var experimentId int
	for rows.Next() {
		rows.Scan(&experimentId)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return err
	}

	for _, option := range options {
		_, err = e.DB.Exec(`
		INSERT INTO Option (experiment_id, value)
		VALUES (?, ?)
		`,
			experimentId,
			option.Value,
		)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (e ExperimentModel) Update(experimentId, prompt string) error {
	_, err := e.DB.Exec(`
		UPDATE Experiment
		SET prompt = ?
		WHERE Experiment.id = ?
	`,
		prompt,
		experimentId,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (e ExperimentModel) Delete(experimentId string) error {
	_, err := e.DB.Exec(`
		DELETE FROM Experiment
		WHERE Experiment.id = ?
	`,
		experimentId,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
