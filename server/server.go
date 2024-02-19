package server

import (
	"bytes"
	"database/sql"
	"doubleblind/models"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Env struct {
	experiments   models.ExperimentModel
	options       models.OptionModel
	participants  models.ParticipantModel
	votes         models.VoteModel
	dataCh        chan models.Experiment
	clients       []http.ResponseWriter
	templatePaths map[string]map[string][]string
}

func Run(host string, port int) {
	// Logger setup
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Database setup
	db, err := sql.Open("sqlite3", "./doubleblind.db")

	if err != nil {
		log.Fatal("Error establishing DB connection", err)
	}
	defer os.Remove("./doubleblind.db")
	defer db.Close()

	// Model setup
	options := models.OptionModel{
		DB: db,
	}

	participants := models.ParticipantModel{
		DB: db,
	}

	votes := models.VoteModel{
		DB: db,
	}

	experiments := models.ExperimentModel{
		DB:               db,
		OptionModel:      &options,
		ParticipantModel: &participants,
	}

	dataCh := make(chan models.Experiment)
	clients := []http.ResponseWriter{}
	templatePaths := map[string]map[string][]string{
		"/": {
			"GET": {
				"templates/index.html",
				"templates/components/experiment.html",
				"templates/components/experiment_list.html",
				"templates/partials/header.html",
				"templates/partials/navigation.html",
				"templates/partials/footer.html",
			},
		},
		"/api/*/experiments/*": {
			"GET": {
				"templates/components/option_container.html",
				"templates/components/option_list.html",
			},
		},
		"/api/*/experiments/*/{experimentId}/": {
			"GET": {
				"templates/components/option_container.html",
				"templates/components/option_list.html",
			},
			"PATCH": {
				"templates/components/experiment.html",
			},
			"DELETE": {},
		},
	}

	env := &Env{
		experiments:   experiments,
		options:       options,
		participants:  participants,
		votes:         votes,
		dataCh:        dataCh,
		clients:       clients,
		templatePaths: templatePaths,
	}

	env.populateDB()
	r := env.NewRouter()
	go env.clientBroker()
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), r.Router))
}

func (env *Env) clientBroker() {
	templatePaths := []string{
		"templates/components/option_list.html",
	}

	tmpl := template.Must(template.ParseFiles(templatePaths...))

	for data := range env.dataCh {
		var b bytes.Buffer

		err := tmpl.Execute(&b, data)
		if err != nil {
			fmt.Println(err)
		}
		message := b.String()

		eventName := fmt.Sprintf("vote-%s", data.Id)

		fmt.Printf("event: %s\ndata: %s\n\n", eventName, message)
		for _, w := range env.clients {
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventName, message)
			w.(http.Flusher).Flush()
		}

		fmt.Printf("Sent message to %d clients\n", len(env.clients))
	}
}

func (env *Env) populateDB() {
	err := env.experiments.Setup()
	if err != nil {
		log.Fatal("Error setting up DB: ", err)
	}
	err = env.options.Setup()
	if err != nil {
		log.Fatal("Error setting up DB: ", err)
	}
	err = env.participants.Setup()
	if err != nil {
		log.Fatal("Error setting up DB: ", err)
	}
	err = env.votes.Setup()
	if err != nil {
		log.Fatal("Error setting up DB: ", err)
	}
}

func (env *Env) responder(w http.ResponseWriter, r *http.Request, data any) {
	if r.Header.Get("Accept") == "application/json" {
		jsonResponse(w, data)
	} else {
		rctx := chi.RouteContext(r.Context())
		routePattern := strings.Join(rctx.RoutePatterns, "")
		fmt.Println("route:", routePattern)
		fmt.Println("method:", r.Method)

		templatePaths := env.templatePaths[routePattern][r.Method]
		fmt.Println(templatePaths)

		htmlResponse(w, data, templatePaths)
	}
}

func jsonResponse(w http.ResponseWriter, output any) {
	bytes, err := json.Marshal(output)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func htmlResponse(w http.ResponseWriter, data any, templatePaths []string) {
	if len(templatePaths) == 0 {
		return
	}

	tmpl := template.Must(template.ParseFiles(templatePaths...))
	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}
