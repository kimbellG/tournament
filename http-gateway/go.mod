module github.com/kimbellG/tournament/http

go 1.16

replace github.com/kimbellG/tournament/core => ../bl

require (
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/kimbellG/kerror v0.0.0-20210819100523-8eb79808c2bd
	github.com/kimbellG/tournament/core v0.0.0-00010101000000-000000000000
	google.golang.org/genproto v0.0.0-20210818220304-27ea9cc85d9f // indirect
	google.golang.org/grpc v1.40.0
)
