package service

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/kimbellG/tournament/core/controller"
	"github.com/kimbellG/tournament/core/handler"
	pb "github.com/kimbellG/tournament/core/handler/grpc"
	"github.com/kimbellG/tournament/core/infrastructure"
	"github.com/kimbellG/tournament/core/repository"
	"github.com/kimbellG/tournament/core/tx"
	"github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvl = "debug"
	}

	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}

	logrus.SetLevel(ll)
}

func StartServer() {
	listener, err := net.Listen("tcp", os.Getenv("SERVICE_ADDRESS"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialization database: %v", err)
	}

	store := tx.NewStore(db)
	userRepo := &repository.UserRepository{}
	tournamentRepo := &repository.TournamentRepository{}

	userController := controller.NewUserController(userRepo, store)
	tournamentController := controller.NewTournamentController(tournamentRepo, userRepo, store)

	handler := handler.NewServiceHandler(userController, tournamentController)

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(infrastructure.UnaryInterceptor),
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTournamentServiceServer(grpcServer, handler)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve of grpc server: %v", err)
	}
}

func InitDB() (*sql.DB, error) {
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
