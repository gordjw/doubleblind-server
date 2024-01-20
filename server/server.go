package server

import (
	"database/sql"
	"doubleblind/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Env struct {
	experiments		models.ExperimentModel
	options			models.OptionModel
	participants	models.ParticipantModel
}

func Run(host string, port int) {	

	// Logger setup
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)


	// Database setup
	db, err := sql.Open("sqlite3", "./doubleblind.db")

	if err != nil {
		log.Fatal("Error establishing DB connection", err)
	}
	defer os.Remove("./doubleblind.db")
	defer db.Close()

	options := models.OptionModel{
		DB: db, 
	}

	participants := models.ParticipantModel{
		DB: db,
	}

	experiments := models.ExperimentModel{
		DB: db,
		OptionModel: &options,
		ParticipantModel: &participants,
	}

	env := &Env {
		experiments:	experiments,
		options:		options,
		participants:	participants,
	}

	err = env.experiments.Setup()
	if err != nil {
		log.Fatal("Error setting up DB: ", err)
	}
	err = env.options.Setup()
	if err != nil {
		log.Fatal("Error setting up DB: ", err)
	}
	err = env.participants.Setup()
	if err != nil {
		log.Fatal("Error setting up DB: ", err)
	}



	// Multiplexer setup

	/**
	get		/							(show homepage)
	post	/experiment					(create a new experiment)
	get		/experiment/:id				(get the details of an experiment. including votes)
	post	/experiment/:id/option		(create a new option for this experiment)
	post	/experiment/:id/vote/:id	(vote for an option)
	*/

	mux := http.NewServeMux()
	
	mux.HandleFunc("/experiment/{experiment_id}/vote/{option_id}", env.postVote)
	mux.HandleFunc("/experiment/{experiment_id}/option", env.postOption)
	mux.HandleFunc("/experiment/{experiment_id}", env.getExperiment)
	mux.HandleFunc("/experiment", env.postExperiment)
	mux.HandleFunc("/", env.handleRoot)

	// Start the server and listen on port 80
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), mux))
}

func (env *Env) handleRoot		(w http.ResponseWriter, r *http.Request) { 
	experiments, err := env.experiments.All()
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
	}

	fmt.Println("Viewed homepage") 

	output, err := json.Marshal(experiments)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
	}

	w.Write(output)
}


func (env *Env) getExperiment	(w http.ResponseWriter, r *http.Request) { 
	experiment, err := env.experiments.One(1)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
	}

	fmt.Println("Viewed experiment") 

	output, err := json.Marshal(experiment)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
	}

	w.Write(output)

}
func (env *Env) postExperiment	(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, fmt.Sprintf("Created new experiment")); fmt.Println("Created new experiment") }
func (env *Env) postVote		(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, fmt.Sprintf("Voted for option")); fmt.Println("Voted for option") }
func (env *Env) postOption		(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, fmt.Sprintf("Created new option")); fmt.Println("Created new option") }




