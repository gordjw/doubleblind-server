package server

import (
	"database/sql"
	"doubleblind/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Env struct {
	experiments  models.ExperimentModel
	options      models.OptionModel
	participants models.ParticipantModel
	votes        models.VoteModel
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

	votes := models.VoteModel{
		DB: db,
	}

	experiments := models.ExperimentModel{
		DB:               db,
		OptionModel:      &options,
		ParticipantModel: &participants,
		VoteModel:        &votes,
	}

	env := &Env{
		experiments:  experiments,
		options:      options,
		participants: participants,
		votes:        votes,
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

	/**
	 *  Setting up the API Router and routes
	 */
	apiRouter := chi.NewRouter()
	apiRouter.Use(middlewareAuth)
	apiRouter.Use(middlewareJSONResponse)

	apiRouter.Get("/experiment", env.getExperiments)
	apiRouter.Get("/experiment/{experiment_id}", env.getExperiment)
	apiRouter.Post("/experiment", env.postExperiment)

	apiRouter.Post("/experiment/{experiment_id}/vote/{option_id}", env.postVote)

	/**
	 *  Setting up the Main Router and routes
	 */
	r := chi.NewRouter()
	r.Use(middlewareCORS)
	r.Use(middlewareLogger)

	r.Handle("/", http.FileServer(http.Dir("./public")))

	r.Mount("/api", apiRouter)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), r))
}

func jsonResponse(output any, w http.ResponseWriter) {
	bytes, err := json.Marshal(output)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
