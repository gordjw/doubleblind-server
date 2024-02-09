package models

import (
	"database/sql"
	"log"
)

type Experiment struct {
	Id           string
	Prompt       string
	Options      []Option
	Participants []Participant
	OrganiserId  string
	Organiser    Participant
	Votes        []Vote
}

type ExperimentModel struct {
	DB               *sql.DB
	OptionModel      *OptionModel
	ParticipantModel *ParticipantModel
	VoteModel        *VoteModel
}

// Returns the winning Option and the number of votes it received
func (e Experiment) winner() (Option, int) {
	max := 0
	var winner Option

	// for i := 0; i < len(e.Options); i++ {
	// 	if e.Options[i].Votes > max {
	// 		max = e.Options[i].Votes
	// 		winner = e.Options[i]
	// 	}
	// }

	return winner, max
}

// Returns true if the Experiment is waiting on 1 or more Participants to vote
func (e Experiment) isOpen() bool {
	return true
}

func (e ExperimentModel) All() ([]Experiment, error) {
	rows, err := e.DB.Query("SELECT id, prompt, organiserId FROM Experiment")
	if err != nil {
		log.Println(err)
		return []Experiment{}, err
	}
	defer rows.Close()

	var experiments []Experiment

	for rows.Next() {
		var experiment Experiment

		err := rows.Scan(&experiment.Id, &experiment.Prompt, &experiment.OrganiserId)
		if err != nil {
			log.Println(err)
			return []Experiment{}, err
		}

		options, err := e.OptionModel.AttachedToExperiment(experiment.Id)
		if err != nil {
			log.Println(err)
			return []Experiment{}, err
		}
		experiment.Options = options

		participants, err := e.ParticipantModel.AttachedToExperiment(experiment.Id)
		if err != nil {
			log.Println(err)
			return []Experiment{}, err
		}
		experiment.Participants = participants

		organiser, err := e.ParticipantModel.One(experiment.OrganiserId)
		if err != nil {
			log.Println(err)
			return []Experiment{}, err
		}
		experiment.Organiser = organiser

		experiments = append(experiments, experiment)
	}

	if err = rows.Err(); err != nil {
		return []Experiment{}, err
	}

	return experiments, nil
}

func (e ExperimentModel) One(id string) (Experiment, error) {
	rows, err := e.DB.Query(
		`SELECT Experiment.id, Experiment.prompt, Experiment.organiserId, Option.id, Option.value
		FROM Experiment
		LEFT JOIN Option
		ON Experiment.id = Option.experiment_id
		WHERE Experiment.id = ?`,
		id,
	)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return Experiment{}, err
	}

	var experiment Experiment

	for rows.Next() {
		var option Option

		err := rows.Scan(&experiment.Id, &experiment.Prompt, &experiment.OrganiserId, &option.Id, &option.Value)
		if err != nil {
			log.Println(err)
			return Experiment{}, err
		}

		experiment.Options = append(experiment.Options, option)
	}
	if err = rows.Err(); err != nil {
		return Experiment{}, err
	}

	return experiment, nil

	// err := row.Scan(&experiment.Id, &experiment.Prompt, &experiment.OrganiserId)
	// if err != nil {
	// 	log.Println(err)
	// 	return Experiment{}, err
	// }

	// options, err := e.OptionModel.AttachedToExperiment(experiment.Id)
	// if err != nil {
	// 	log.Println(err)
	// 	return Experiment{}, err
	// }
	// experiment.Options = options

	// participants, err := e.ParticipantModel.AttachedToExperiment(experiment.Id)
	// if err != nil {
	// 	log.Println(err)
	// 	return Experiment{}, err
	// }
	// experiment.Participants = participants

	// organiser, err := e.ParticipantModel.One(experiment.OrganiserId)
	// if err != nil {
	// 	log.Println(err)
	// 	return Experiment{}, err
	// }
	// experiment.Organiser = organiser

	// votes, err := e.VoteModel.AttachedToExperiment(experiment.Id)
	// if err != nil {
	// 	log.Println(err)
	// 	return Experiment{}, err
	// }
	// experiment.Votes = votes

}
