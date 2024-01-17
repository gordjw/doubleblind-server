package server

import (
	"fmt"
	"log"
	"net/http"
)

func Run() {	
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
	mux.HandleFunc("/", handleRoot)

	// Serve static files
	fs := http.FileServer(http.Dir("static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start the server and listen on port 80
	log.Fatal(http.ListenAndServe(":8090", mux))
}

func handleRoot		(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, fmt.Sprintf("Viewed homepage")); fmt.Println("Viewed homepage") }
func getExperiment	(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, fmt.Sprintf("Viewed experiement")); fmt.Println("Viewed experiment") }
func postExperiment	(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, fmt.Sprintf("Created new experiment")); fmt.Println("Created new experiment") }
func postVote		(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, fmt.Sprintf("Voted for option")); fmt.Println("Voted for option") }
func postOption		(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, fmt.Sprintf("Created new option")); fmt.Println("Created new option") }
