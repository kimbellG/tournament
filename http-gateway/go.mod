module github.com/kimbellG/tournament/http

go 1.16

replace github.com/kimbellG/tournament/core => ../core

replace github.com/kimbellG/kerror => ../../tournament-error

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.3.0
	github.com/kimbellG/kerror v0.0.0-20210820142247-2f3f8ab8756f
	github.com/kimbellG/tournament/core v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.40.0
)
