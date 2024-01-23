package models

import(
	"log"
	"database/sql"
)

type Option struct {
	Id int
	Value string
}

type OptionModel struct {
	DB *sql.DB
}


func (o OptionModel) All() (*[]Option, error) {
	rows, err := o.DB.Query("SELECT * FROM Option")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []Option

	for rows.Next() {
		var option Option

		err := rows.Scan(&option.Id, &option.Value)
		if err != nil {
			return nil, err
		}

		options = append(options, option)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &options, nil
}


func (o OptionModel) AttachedToExperiment(id int) (*[]Option, error) {
	var options []Option

	rows, err := o.DB.Query(`SELECT id, value FROM Option WHERE experiment_id = ?`, id)
	if err != nil {
		log.Println("Error in OptionModel.AttachedToExperiment(%d): %v", id, err)
		return nil, err
	}
	
	for rows.Next() {
		var option Option

		err := rows.Scan( &option.Id, &option.Value )
		if err != nil {
			log.Println("Error in OptionModel.AttachedToExperiment(%d): %v", id, err)
			return nil, err
		}

		options = append(options, option)
	}

	return &options, nil
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