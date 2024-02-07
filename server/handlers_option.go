package server

import (
	"fmt"
	"net/http"
)

func (env *Env) postOption(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, fmt.Sprintf("Created new option"))
	fmt.Println("Created new option")
}
