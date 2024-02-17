package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (env *Env) postVote(w http.ResponseWriter, r *http.Request) {
	experimentId, err := strconv.Atoi(chi.URLParam(r, "experimentId"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	optionId, err := strconv.Atoi(chi.URLParam(r, "optionId"))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	participantId, err := strconv.Atoi(r.Context().Value(ContextAuthKey).(string))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = env.votes.VoteFor(experimentId, optionId, participantId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
