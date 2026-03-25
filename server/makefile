swag:
	docker compose exec app swag i --output docs --generalInfo cmd/server/main.go

reload:
	docker compose restart app ai

kill:
	kill -9 $(shell lsof -t -i:1323)

run:
	docker compose up -d
	docker compose logs -f app ai

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