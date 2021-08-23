module github.com/kimbellG/tournament/http

go 1.16

replace github.com/kimbellG/tournament/core => ../core

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/kimbellG/kerror v0.0.0-20210820142247-2f3f8ab8756f
	github.com/kimbellG/tournament v0.0.0-20210809141859-98b34e5f6f05 // indirect
	github.com/kimbellG/tournament/core v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	google.golang.org/grpc v1.40.0
)
