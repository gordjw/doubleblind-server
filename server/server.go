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
	votes			models.VoteModel
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

	// Model setup
	options := models.OptionModel{
		DB: db, 
	}

	participants := models.ParticipantModel{
		DB: db,
	}

	votes := models.VoteModel {
		DB: db,
	}

	experiments := models.ExperimentModel{
		DB: db,
		OptionModel: &options,
		ParticipantModel: &participants,
		VoteModel: &votes,
	}

	env := &Env {
		experiments:	experiments,
		options:		options,
		participants:	participants,
		votes:			votes,
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
	err = env.votes.Setup()
	if err != nil {
		log.Fatal("Error setting up DB: ", err)
	}

	// Multiplexer setup
	mux := http.NewServeMux()
	
	mux.HandleFunc("/experiment/{experiment_id}/vote/{option_id}", env.postVote)
	mux.HandleFunc("/experiment/{experiment_id}/option", env.postOption)
	mux.HandleFunc("/experiment/{experiment_id}", env.getExperiment)
	mux.HandleFunc("/experiment", env.postExperiment)
	mux.HandleFunc("/", env.handleRoot)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), mux))
}

func (env *Env) handleRoot (w http.ResponseWriter, r *http.Request) { 
	experiments, err := env.experiments.All()
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}

	output, err := json.Marshal(experiments)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}

	w.Write(output)
}


func (env *Env) getExperiment (w http.ResponseWriter, r *http.Request) { 
	experiment, err := env.experiments.One(1)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}

	output, err := json.Marshal(experiment)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}

	w.Write(output)
}


func (env *Env) postExperiment (w http.ResponseWriter, r *http.Request) { 
	fmt.Println(r.Method)
	fmt.Println(r.Body)

	fmt.Fprint(w, fmt.Sprintf("Created new experiment")); 
	fmt.Println("Created new experiment") 
}

func (env *Env) postVote (w http.ResponseWriter, r *http.Request) { 
	fmt.Fprint(w, fmt.Sprintf("Voted for option")); 
	fmt.Println("Voted for option") 
}

func (env *Env) postOption (w http.ResponseWriter, r *http.Request) { 
	fmt.Fprint(w, fmt.Sprintf("Created new option")); 
	fmt.Println("Created new option") 
}