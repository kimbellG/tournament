package service

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kimbellG/tournament/http/controller"
	"github.com/kimbellG/tournament/http/handler"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func initConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	return nil
}

func StartGateway() {
	if err := initConfig(); err != nil {
		log.Fatalf("Failed to init config file: %v", err)
	}

	conn, err := grpc.Dial(os.Getenv("SERVICE_ADDRESS"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect with core service: %v", err)
	}

	router := mux.NewRouter()
	cont := controller.NewTournamentController(conn)

	handler.RegisterUserEndpoints(router, cont)

	fmt.Println("Server is listening on %v", os.Getenv("PORT"))
	if err := http.ListenAndServe(os.Getenv("PORT"), router); err != nil {
		log.Fatalf("Failed to listen http server: %v", err)
	}
}
