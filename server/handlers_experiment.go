package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func (env *Env) getExperiments(w http.ResponseWriter, r *http.Request) {
	experiments, err := env.experiments.All()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonResponse(experiments, w)
}

func (env *Env) getExperiment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "experiment_id")
	experiment, err := env.experiments.One(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonResponse(experiment, w)
}

func (env *Env) postExperiment(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	fmt.Println(r.Body)

	fmt.Fprint(w, fmt.Sprintf("Created new experiment"))
	fmt.Println("Created new experiment")
}
