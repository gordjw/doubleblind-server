package server

import (
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

	// templatePaths := []string{
	// 	"templates/index.html",
	// 	"templates/components/experiment.html",
	// 	"templates/components/experiment_list.html",
	// 	"templates/partials/header.html",
	// 	"templates/partials/navigation.html",
	// 	"templates/partials/footer.html",
	// }

	// tmpl := template.Must(template.ParseFiles(templatePaths...))
	// err = tmpl.Execute(w, experiments)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
