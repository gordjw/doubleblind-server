package models

import(
	"fmt"
	"database/sql"
)

type Option struct {
	Id int
	Value string
	Votes int
}

type OptionModel struct {
	DB *sql.DB
}


func (o OptionModel) All() ([]Option, error) {
	rows, err := o.DB.Query("SELECT * FROM Option")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []Option

	for rows.Next() {
		var option Option

		err := rows.Scan(&option.Id, &option.Value, &option.Votes)
		if err != nil {
			return nil, err
		}

		options = append(options, option)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return options, nil
}


func (o OptionModel) AttachedToExperiment(id int) ([]Option, error) {
	var options []Option

	rows, err := o.DB.Query(`SELECT id, value, votes FROM Option WHERE experiment_id = ?`, id)
	if err != nil {
		fmt.Println("Error in OptionModel.AttachedToExperiment(%d): %v", id, err)
		return options, err
	}
	
	for rows.Next() {
		var option Option

		err := rows.Scan( &option.Id, &option.Value, &option.Votes )
		if err != nil {
			fmt.Println("Error in OptionModel.AttachedToExperiment(%d): %v", id, err)
			return nil, err
		}

		options = append(options, option)
	}

	return options, nil
}



func (o OptionModel) Setup() error {
	_, err := o.DB.Exec(`
		DROP TABLE IF EXISTS Option;
		CREATE TABLE Option (
			id				INTEGER PRIMARY KEY AUTOINCREMENT,
			experiment_id	INTEGER NOT NULL,
			value			VARCHAR(128) NOT NULL,
			votes			INTEGER NOT NULL,
			FOREIGN KEY (experiment_id)
				REFERENCES Experiment(id)
				ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	_, err = o.DB.Exec(`
		INSERT INTO Option (experiment_id, value, votes) VALUES ('1', 'Two Blind Mice', '0'), ('1', 'CBD Dumplings', '0'), ('1', 'Asian Cafe', '0')
	`)
	if err != nil {
		return err
	}

	return nil
}