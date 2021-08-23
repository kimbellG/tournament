module github.com/kimbellG/tournament/core

go 1.16

replace github.com/kimbellG/kerror => ../../tournament-error

require (
	github.com/google/uuid v1.3.0
	github.com/jackc/pgx/v4 v4.13.0
	github.com/joho/godotenv v1.3.0
	github.com/kimbellG/kerror v0.0.0-20210820142247-2f3f8ab8756f
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/sys v0.0.0-20210823070655-63515b42dcdf // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210821163610-241b8fcbd6c8 // indirect
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
)
