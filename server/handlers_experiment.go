package server

import (
	"context"
	"doubleblind/models"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
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
	experimentId := chi.URLParam(r, "experimentId")
	experiment, err := env.experiments.One(experimentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	templatePaths := []string{
		"templates/components/option_list.html",
	}

	tmpl := template.Must(template.ParseFiles(templatePaths...))
	err = tmpl.Execute(w, experiment)
	if err != nil {
		fmt.Println(err)
	}
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

	fmt.Fprint(w, fmt.Sprintf("Created new experiment\n"))
	fmt.Fprintf(w, "Experiment: %+v\n", experiment)
	fmt.Println("Created new experiment")
}
