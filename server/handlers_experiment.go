package server

import (
	"doubleblind/models"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
)

func (env *Env) getAllExperiments(w http.ResponseWriter, r *http.Request) {
	experiments, err := env.experiments.All()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	env.responder(w, r, experiments)
}

func (env *Env) getExperiment(w http.ResponseWriter, r *http.Request) {
	experimentId := chi.URLParam(r, "experimentId")
	experiment, err := env.experiments.One(experimentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	env.responder(w, r, experiment)
}

func (env *Env) postExperiment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	organiserId := 1

	var experiment models.Experiment

	decoder := schema.NewDecoder()
	err = decoder.Decode(&experiment, r.PostForm)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	env.experiments.Add(experiment.Prompt, organiserId, experiment.Options)

	env.responder(w, r, experiment)
}

func (env *Env) patchExperiment(w http.ResponseWriter, r *http.Request) {
	experimentId := chi.URLParam(r, "experimentId")
	prompt := "I'm a new question"
	experiment := models.Experiment{
		Id:     experimentId,
		Prompt: prompt,
	}
	err := env.experiments.Update(experimentId, prompt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	env.responder(w, r, experiment)
}

func (env *Env) deleteExperiment(w http.ResponseWriter, r *http.Request) {
	experimentId := chi.URLParam(r, "experimentId")

	err := env.experiments.Delete(experimentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	env.responder(w, r, nil)
}
