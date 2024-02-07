package server

import (
	"fmt"
	"net/http"
)

func (env *Env) postVote(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, fmt.Sprintf("Voted for option"))
	fmt.Println("Voted for option")
}
