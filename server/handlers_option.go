package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func (env *Env) postOption(w http.ResponseWriter, r *http.Request) {
	experimentId := chi.URLParam(r, "experimentId")
	optionId := chi.URLParam(r, "optionId")

	fmt.Fprintf(w, fmt.Sprintf("Not implemented yet\n"))
	fmt.Fprintf(w, "Experiment: %v %v\n", experimentId, optionId)
	fmt.Println("Not implemented yet")
}

func (env *Env) patchOption(w http.ResponseWriter, r *http.Request) {
	experimentId := chi.URLParam(r, "experimentId")
	optionId := chi.URLParam(r, "optionId")

	fmt.Fprintf(w, fmt.Sprintf("Not implemented yet\n"))
	fmt.Fprintf(w, "Experiment: %v %v\n", experimentId, optionId)
	fmt.Println("Not implemented yet")
}

func (env *Env) deleteOption(w http.ResponseWriter, r *http.Request) {
	experimentId := chi.URLParam(r, "experimentId")
	optionId := chi.URLParam(r, "optionId")

	fmt.Fprintf(w, fmt.Sprintf("Not implemented yet\n"))
	fmt.Fprintf(w, "Experiment: %v %v\n", experimentId, optionId)
	fmt.Println("Not implemented yet")
}
