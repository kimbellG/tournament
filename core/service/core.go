package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/kimbellG/tournament/core/controller"
	"github.com/kimbellG/tournament/core/debugutil"
	"github.com/kimbellG/tournament/core/handler"
	pb "github.com/kimbellG/tournament/core/handler/grpc"
	"github.com/kimbellG/tournament/core/interceptor"
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()

		oscall := <-c
		log.Printf("system call: %v", oscall)
	}()

	startServer(ctx)
}

func startServer(ctx context.Context) {
	listener, err := net.Listen("tcp", os.Getenv("SERVICE_ADDRESS"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer debugutil.Close(listener)

	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialization database: %v", err)
	}
	defer debugutil.Close(db)

	srv := newServer(listener, startHandler(db))

	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("Failed to serve of grpc server: %v", err)
		}
	}()

	log.Println("Core service started on", os.Getenv("SERVICE_ADDRESS"))
	<-ctx.Done()
	gracefulShutdown(srv)
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

func startHandler(db *sql.DB) *handler.ServiceHandler {
	store := tx.NewStore(db)
	userRepo := &repository.UserRepository{}
	tournamentRepo := &repository.TournamentRepository{}

	userController := controller.NewUserController(userRepo, store)
	tournamentController := controller.NewTournamentController(tournamentRepo, userRepo, store)

	return handler.NewServiceHandler(userController, tournamentController)
}

func newServer(listener net.Listener, handler *handler.ServiceHandler) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Error(interceptor.Log)),
	}
	srv := grpc.NewServer(opts...)
	pb.RegisterTournamentServiceServer(srv, handler)

	return srv
}

func gracefulShutdown(srv *grpc.Server) {
	log.Println("Core service is starting graceful shutdown")

	srv.GracefulStop()
	log.Println("Core service stopped")

}
