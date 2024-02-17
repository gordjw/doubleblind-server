package server

import (
	"context"
	"doubleblind/models"
	"encoding/json"
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

	ctx := context.WithValue(r.Context(), ContextJsonResponseKey, experiments)
	r = r.WithContext(ctx)
	sendJson(w, r)
}

func (env *Env) getExperiment(w http.ResponseWriter, r *http.Request) {
	fmt.Println(chi.URLParam(r, "experimentId"))
	experimentId := chi.URLParam(r, "experimentId")
	experiment, err := env.experiments.One(experimentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	ctx := context.WithValue(r.Context(), ContextJsonResponseKey, experiment)
	r = r.WithContext(ctx)
	sendJson(w, r)
}

func (env *Env) postExperiment(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	fmt.Println(r.Body)

	organiserId := 1

	var e models.Experiment

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&e)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	env.experiments.Add(e.Prompt, organiserId, e.Options)

	fmt.Fprint(w, fmt.Sprintf("Created new experiment\n"))
	fmt.Fprintf(w, "Experiment: %+v\n", e)
	fmt.Println("Created new experiment")
}
