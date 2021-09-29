test: compose
	cd itest && go test -tags integration
	make clean

compose:
	docker-compose up -d
	sleep 3

clean:
	docker stop core_postgres_1
	docker rm core_postgres_1
