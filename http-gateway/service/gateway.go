package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/kimbellG/tournament/http/controller"
	"github.com/kimbellG/tournament/http/controller/interceptor"
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()

		oscall := <-c
		log.Printf("system call: %v", oscall)
	}()

	startGateway(ctx)
}

func startGateway(ctx context.Context) {
	if err := initConfig(); err != nil {
		log.Fatalf("Failed to init config file: %v", err)
	}

	conn, err := grpc.Dial(os.Getenv("SERVICE_ADDRESS"), grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptor.Error))
	if err != nil {
		log.Fatalf("Failed to connect with core service: %v", err)
	}
	defer conn.Close()

	srv := &http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: startRouter(conn),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen http server: %v", err)
		}
	}()
	fmt.Println("Server is listening on", os.Getenv("PORT"))
	<-ctx.Done()
	fmt.Println("Server is starting gracefull shutdown")

	gracefullShutdown(srv, 5*time.Second)

}

func startRouter(conn *grpc.ClientConn) *mux.Router {
	router := mux.NewRouter()
	cont := controller.NewTournamentController(conn)
	authmid := handler.AuthenticationMiddleware{
		TokenPassword: os.Getenv("TK_PASSWORD"),
		NotAuthPaths: []string{
			"/" + handler.UserPath,
			"/" + handler.LogInPath,
		},
	}

	handler.RegisterUserEndpoints(router, cont)
	handler.RegisterTournamentEndpoints(router, cont)
	router.Use(authmid.Middleware)

	return router
}

func gracefullShutdown(srv *http.Server, timeout time.Duration) {
	ctxShutdown, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Gracefull shutdown is failed: %s", err)
	}

	log.Println("Server closed")
}
