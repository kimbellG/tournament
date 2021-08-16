test: compose
	cd itest && go test -tags integration
	docker stop bl_postgres_1
	docker rm bl_postgres_1

compose:
	docker-compose up -d
	sleep 3

clean:
	docker stop bl_postgres_1
	docker rm bl_postgres_1
