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
	db, err := sql.Open("sqlite3", "./doubleblind.db")

	if err != nil {
		log.Fatal("Error establishing DB connection", err)
	}
	defer os.Remove("./doubleblind.db")
	defer db.Close()

	env := &Env {
		experiments:	models.ExperimentModel{DB: db},
		options:		models.OptionModel{DB: db},
		participants:	models.ParticipantModel{DB: db},
	}

	err = env.experiments.Setup()
	if err != nil {
		log.Fatal("Error setting up DB: ", err)
	}
	err = env.options.Setup()
	if err != nil {
		log.Fatal("Error setting up DB: ", err)
	}


	mux := http.NewServeMux()

	/**
	get		/							(show homepage)
	post	/experiment					(create a new experiment)
	get		/experiment/:id				(get the details of an experiment. including votes)
	post	/experiment/:id/option		(create a new option for this experiment)
	post	/experiment/:id/vote/:id	(vote for an option)
	*/
	
	mux.HandleFunc("/experiment/{experiment_id}/vote/{option_id}", postVote)
	mux.HandleFunc("/experiment/{experiment_id}/option", postOption)
	mux.HandleFunc("/experiment/{experiment_id}", getExperiment)
	mux.HandleFunc("/experiment", postExperiment)
	mux.HandleFunc("/", env.handleRoot)

	// Serve static files
	fs := http.FileServer(http.Dir("static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

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
func getExperiment	(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, fmt.Sprintf("Viewed experiment")); fmt.Println("Viewed experiment") }
func postExperiment	(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, fmt.Sprintf("Created new experiment")); fmt.Println("Created new experiment") }
func postVote		(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, fmt.Sprintf("Voted for option")); fmt.Println("Voted for option") }
func postOption		(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, fmt.Sprintf("Created new option")); fmt.Println("Created new option") }




