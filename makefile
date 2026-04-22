swag:
	docker compose exec app swag i --output docs --generalInfo cmd/server/main.go

reload:
	docker compose restart app resize-worker ai-batch

kill:
	kill -9 $(shell lsof -t -i:1323)

up:
	docker compose up -d
	docker compose logs -f app resize-worker ai-batch

run-front:
	cd front && yarn dev

stop:
	docker compose down

sqlc:
	docker compose exec app sqlc generate -f db/sqlc.yaml

migrate:
	docker compose exec app go run cmd/db/migrate/main.go

seed:
	docker compose exec app go run cmd/db/seed/main.go

destroy:
	docker compose exec app go run cmd/db/destroy/main.go

refresh:
	make destroy
	make migrate
	make seed

ai-batch:
	docker compose up -d ai-batch

init:
	make up
	make migrate
	make seed

build-api-lambda:
	cd server && \
	GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/lambda/api/main.go && \
	zip function.zip bootstrap && \
	rm bootstrap