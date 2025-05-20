swag:
	~/go/bin/swag i --output docs --generalInfo cmd/server/main.go
kill:
	kill -9 $(shell lsof -t -i:1323)

run:
	docker compose up -d
	docker compose logs -f app
