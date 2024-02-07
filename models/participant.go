package models

import (
	"database/sql"
	"log"
)

type Participant struct {
	Id    string
	Name  string
	Email string
}

type ParticipantModel struct {
	DB *sql.DB
}

func (p ParticipantModel) One(id string) (Participant, error) {
	row := p.DB.QueryRow(`SELECT id, name, email FROM Participant WHERE id = ?`, id)

	var participant Participant

	err := row.Scan(&participant.Id, &participant.Name, &participant.Email)
	if err != nil {
		log.Printf("Error in ParticipantModel.One(%s): %v\n", id, err)
		return Participant{}, err
	}

	return participant, nil
}

func (p ParticipantModel) AttachedToExperiment(id string) ([]Participant, error) {
	var participants []Participant

	rows, err := p.DB.Query(`SELECT id, name, email FROM Participant WHERE experiment_id = ?`, id)
	if err != nil {
		log.Printf("Error in ParticipantModel.AttachedToExperiment(%s): %v\n", id, err)
		return nil, err
	}

	for rows.Next() {
		var participant Participant

		err := rows.Scan(&participant.Id, &participant.Name, &participant.Email)
		if err != nil {
			log.Println("Error in ParticipantModel.AttachedToExperiment(%s): %v", id, err)
			return nil, err
		}

		participants = append(participants, participant)
	}

	return participants, nil
}
