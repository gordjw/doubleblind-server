package models

import (
	"database/sql"
	"log"
)

type Option struct {
	Id    string
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

		err := rows.Scan(&option.Id, &option.Value)
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

func (o OptionModel) AttachedToExperiment(id string) ([]Option, error) {
	var options []Option

	rows, err := o.DB.Query(`SELECT id, value FROM Option WHERE experiment_id = ?`, id)
	if err != nil {
		log.Printf("Error in OptionModel.AttachedToExperiment(%s): %v\n", id, err)
		return nil, err
	}

	for rows.Next() {
		var option Option

		err := rows.Scan(&option.Id, &option.Value)
		if err != nil {
			log.Printf("Error in OptionModel.AttachedToExperiment(%s): %v\n", id, err)
			return nil, err
		}

		options = append(options, option)
	}

	return options, nil
}
