package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/BalabanovA898/project/account/src/config"
	"github.com/BalabanovA898/project/account/src/delivery"
	"github.com/BalabanovA898/project/account/src/repository/postgres"
	"github.com/BalabanovA898/project/account/src/usecase"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewRefreshTokenRepository(db)
	teamRepo := postgres.NewTeamRepository(db)

	authUseCase := usecase.NewAuthUseCase(userRepo, tokenRepo, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	userUseCase := usecase.NewUserUseCase(userRepo)
	teamUseCase := usecase.NewTeamUseCase(teamRepo)

	server := delivery.NewServer(authUseCase, userUseCase, teamUseCase, cfg.JWTSecret)

	log.Printf("account service listening on %s", cfg.ServerAddr)
	if err := http.ListenAndServe(cfg.ServerAddr, server); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
