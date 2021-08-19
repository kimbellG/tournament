test: compose
	if cd itest && go test -tags integration -v; then cd ../ && make clean; else cd ../ && make clean; fi

compose:
	docker-compose up -d
	sleep 3

clean:
	docker stop bl_postgres_1
	docker rm bl_postgres_1
