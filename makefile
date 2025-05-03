rebuild:
	@docker compose down
	@docker rmi ewanlav/diploma:latest
	@docker build -t ewanlav/diploma:latest .
	@docker compose up -d
	sleep 5
	@docker start diploma-fin-app-1
	@docker ps
	@docker logs diploma-fin-app-1
	sleep 60
	@docker logs diploma-fin-migrate-1

build:
	@docker rmi ewanlav/diploma:latest
	@docker build -t ewanlav/diploma:latest .
	@docker compose up -d
	sleep 5
	@docker start diploma-fin-app-1
	@docker ps
	@docker logs diploma-fin-app-1
	sleep 60
	@docker logs diploma-fin-migrate-1

enterp:
	@docker exec -it diploma-fin-postgres-1 bash -c "PGPASSWORD=qwerty psql -U ewan -p 5005 -d debts"

stop:
	@docker compose down