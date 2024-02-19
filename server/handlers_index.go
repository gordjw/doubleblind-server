package server

import (
	"fmt"
	"log"
	"net/http"
)

func (env *Env) getIndex(w http.ResponseWriter, r *http.Request) {
	experiments, err := env.experiments.All()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	env.responder(w, r, experiments)
}

func (env *Env) getProtectedFoo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getProtectedFoo")
	env.responder(w, r, nil)
}
