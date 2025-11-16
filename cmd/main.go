package main

import (
	"context"
	"github.com/GameXost/Avito_Test_Case/internal/repository"
	"github.com/GameXost/Avito_Test_Case/internal/server"
	"github.com/GameXost/Avito_Test_Case/internal/server/handlers"
	"github.com/GameXost/Avito_Test_Case/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	connStr := os.Getenv("DB_CONN")
	if connStr == "" {
		log.Fatal("DB_CONN is not set")
	}
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("cant create pgxpool: %v", err)
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal("cant ping db: %v", err)
	}

	log.Println("successfull onnection to postgres")

	teamRepo := repository.NewTeamRepo(pool)
	userRepo := repository.NewUserRepo(pool)
	prRepo := repository.NewPullRequestRepo(pool)

	teamService := service.NewTeamService(teamRepo)
	userService := service.NewUserService(userRepo)
	prService := service.NewPRService(prRepo, teamRepo)

	teamHandler := handlers.NewTeamHandler(teamService)
	userHandler := handlers.NewUserHandler(userService)
	prHandler := handlers.NewPRHandler(prService)

	rt := server.NewRouter(userHandler, teamHandler, prHandler)
	router := rt.Init()

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Server started at http://localhost:8080")
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutdown signal received")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped correctly")
}
