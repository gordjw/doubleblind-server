package server

import (
	"context"

	"github.com/go-chi/chi"
)

type Router struct {
	Router *chi.Mux
}

func (env *Env) NewRouter() Router {
	/**
	 *  Setting up the API Router and routes
	 */
	apiRouter := chi.NewRouter()
	apiRouter.Use(middlewareAuth)
	// apiRouter.Use(middlewareJSONResponse)
	apiRouter.Use(middlewareHTMLResponse)
	// apiRouter.Use(TrackRoute)

	apiRouter.Route("/experiments", func(apiRouter chi.Router) {
		apiRouter.Get("/", env.getAllExperiments) // Get all experiments
		apiRouter.Post("/", env.postExperiment)   // Create a new experiment

		apiRouter.Route("/{experimentId}", func(apiRouter chi.Router) {
			apiRouter.Get("/", env.getExperiment)            // Get one experiment
			apiRouter.Patch("/", env.patchExperiment)        // Update
			apiRouter.Delete("/", env.deleteExperiment)      // Delete an experiment
			apiRouter.Post("/vote/{optionId}", env.postVote) // Vote for an experiment

			apiRouter.Route("/options", func(apiRouter chi.Router) {
				apiRouter.Post("/", env.postOption)               // Create a new option
				apiRouter.Patch("/{optionId}", env.patchOption)   // Update an option
				apiRouter.Delete("/{optionId}", env.deleteOption) // Delete an option
			})
		})
	})

	apiRouter.Route("/sse", func(apiRouter chi.Router) {
		apiRouter.Get("/", env.getSSE) // Register a client for SSE updates
	})

	fooRouter := chi.NewRouter()
	fooRouter.Use(env.middlewareOauth)
	fooRouter.Use(middlewareHTMLResponse)
	fooRouter.Get("/", env.getProtectedFoo)

	/**
	 *  Setting up the Main Router and routes
	 */
	r := chi.NewRouter()
	r.Use(middlewareCORS)
	r.Use(middlewareLogger)
	r.Mount("/api", apiRouter)
	r.Mount("/foo", fooRouter)

	r.Get("/", env.getIndex)
	r.Get("/auth/google/login", env.oauthGoogleLogin)
	r.Get("/auth/google/callback", env.oauthGoogleCallback)

	return Router{
		Router: r,
	}
}

func setAuthToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, ContextKeyAuthToken, token)
}

func getAuthToken(ctx context.Context) (string, bool) {
	tokenStr, ok := ctx.Value(ContextKeyAuthToken).(string)
	return tokenStr, ok
}
