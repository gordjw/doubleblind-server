package models

func (e ExperimentModel) Setup() error {
	_, err := e.DB.Exec(`
		DROP TABLE IF EXISTS Experiment;
		CREATE TABLE Experiment (
			id				INTEGER PRIMARY KEY AUTOINCREMENT,
			prompt			VARCHAR(128) NOT NULL,
			organiserId		INTEGER NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	_, err = e.DB.Exec(`
		INSERT INTO Experiment (prompt, organiserId) VALUES('Where do you want to go to dinner?', '1')
	`)
	if err != nil {
		return err
	}

	return nil
}

func (o OptionModel) Setup() error {
	_, err := o.DB.Exec(`
		DROP TABLE IF EXISTS Option;
		CREATE TABLE Option (
			id				INTEGER PRIMARY KEY AUTOINCREMENT,
			experiment_id	INTEGER NOT NULL,
			value			VARCHAR(128) NOT NULL,
			FOREIGN KEY (experiment_id)
				REFERENCES Experiment(id)
				ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	_, err = o.DB.Exec(`
		INSERT INTO Option (experiment_id, value) VALUES ('1', 'Two Blind Mice'), ('1', 'CBD Dumplings'), ('1', 'Asian Cafe')
	`)
	if err != nil {
		return err
	}

	return nil
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
