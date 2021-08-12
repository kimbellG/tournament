package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/kimbellG/tournament-bl/controller"
	"github.com/kimbellG/tournament-bl/handler"
	pb "github.com/kimbellG/tournament-bl/handler/grpc"
	"github.com/kimbellG/tournament-bl/repository"
	"github.com/kimbellG/tournament-bl/tx"
	"google.golang.org/grpc"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load config file: %w", err)
	}
}

func main() {
	listener, err := net.Listen("tcp", os.Getenv("SERVICE_ADDRESS"))
	if err != nil {
		log.Fatalf("Failed to listen: %w", err)
	}

	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialization database: %w", err)
	}

	store := tx.NewStore(db)
	userRepo := &repository.UserRepository{}
	tournamentRepo := &repository.TournamentRepository{}

	userController := controller.NewUserController(userRepo, store)
	tournamentController := controller.NewTournamentController(tournamentRepo, userRepo, store)

	handler := handler.NewServiceHandler(userController, tournamentController)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTournamentServiceServer(grpcServer, handler)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve of grpc server: %w", err)
	}
}

func initDB() (*sql.DB, error) {
	connURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	db, err := sql.Open(os.Getenv("DB_DRIVER"), connURL)
	if err != nil {
		return nil, fmt.Errorf("db connection has failed: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db connection has failed: %w", err)
	}

	return db, nil
}
