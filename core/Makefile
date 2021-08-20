test: compose
	if cd itest && go test -tags integration -v; then make -f ../Makefile clean; else make -f ../Makefile clean; fi

compose:
	docker-compose up -d
	sleep 3

clean:
	docker stop core_postgres_1
	docker rm core_postgres_1
