package models

import(
	"fmt"
	"database/sql"
)

type Participant struct {
	Id int
	Name string
	Email string
}

type ParticipantModel struct {
	DB *sql.DB
}


func (p ParticipantModel) One(id int) (Participant, error) {
	row := p.DB.QueryRow(`SELECT id, name, email FROM Participant WHERE id = ?`, id)

	var participant Participant

	err := row.Scan(&participant.Id, &participant.Name, &participant.Email)
	if err != nil {
		fmt.Println("Error in ParticipantModel.One(%d): %v", id, err)
		return participant, err
	}

	return participant, nil
}


func (p ParticipantModel) AttachedToExperiment(id int) ([]Participant, error) {
	var participants []Participant

	rows, err := p.DB.Query(`SELECT id, name, email FROM Participant WHERE experiment_id = ?`, id)
	if err != nil {
		fmt.Println("Error in ParticipantModel.AttachedToExperiment(%d): %v", id, err)
		return participants, err
	}
	
	for rows.Next() {
		var participant Participant

		err := rows.Scan( &participant.Id, &participant.Name, &participant.Email )
		if err != nil {
			fmt.Println("Error in ParticipantModel.AttachedToExperiment(%d): %v", id, err)
			return nil, err
		}

		participants = append(participants, participant)
	}

	return participants, nil
}



func (p ParticipantModel) Setup() error {
	_, err := p.DB.Exec(`
		DROP TABLE IF EXISTS Participant;
		CREATE TABLE Participant (
			id				INTEGER PRIMARY KEY AUTOINCREMENT,
			experiment_id	INTEGER NOT NULL,
			name			VARCHAR(128) NOT NULL,
			email			VARCHAR(128) NOT NULL,
			FOREIGN KEY (experiment_id)
				REFERENCES Experiment(id)
				ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	_, err = p.DB.Exec(`
		INSERT INTO Participant (experiment_id, name, email) VALUES ('1', 'Jimmy Smits', 'jimmy@smits.com'), ('1', 'Jane Doe', 'jane@doe.com'), ('1', 'Mephisto the Cat', 'evilgrin@alice.com')
	`)
	if err != nil {
		return err
	}

	return nil
}