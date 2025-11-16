package server

import (
	"github.com/GameXost/Avito_Test_Case/internal/server/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	userHandler *handlers.UserHandler
	teamHandler *handlers.TeamHandler
	prHandler   *handlers.PRHandler
}

func NewRouter(
	userHandler *handlers.UserHandler,
	teamHandler *handlers.TeamHandler,
	prHandler *handlers.PRHandler,
) *Router {
	return &Router{
		userHandler: userHandler,
		teamHandler: teamHandler,
		prHandler:   prHandler,
	}
}

func (rt *Router) Init() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", rt.userHandler.SetIsActive)
		r.Get("/getReview", rt.userHandler.GetReview)
	})

	r.Route("/team", func(r chi.Router) {
		r.Post("/add", rt.teamHandler.CreateTeam)
		r.Get("/get", rt.teamHandler.GetTeam)
	})

	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", rt.prHandler.CreatePR)
		r.Post("/merge", rt.prHandler.Merge)
		r.Post("/reassign", rt.prHandler.Reassign)
	})
	return r
}
