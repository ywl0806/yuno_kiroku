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
	cp function.zip ../function.zip

upload-zip:
	aws-vault exec lee-tera -- aws s3 cp function.zip s3://yuno-lambda-zip-production/function.zip
	rm function.zip

update-api-lambda:
	aws-vault exec lee-tera -- aws lambda update-function-code --function-name yuno-api-production --s3-bucket yuno-lambda-zip-production --s3-key function.zip

deploy-resize-lambda:
	aws-vault exec lee-tera -- docker buildx build --platform linux/arm64 --provenance=false \
		-t 173549642885.dkr.ecr.ap-northeast-1.amazonaws.com/yuno-resize:production \
		-f docker/production/lambda-resize/Dockerfile ./server
	docker push 173549642885.dkr.ecr.ap-northeast-1.amazonaws.com/yuno-resize:production

deploy-face-recognition-lambda:
	aws-vault exec lee-tera -- docker buildx build --platform linux/arm64 --provenance=false \
		-t 173549642885.dkr.ecr.ap-northeast-1.amazonaws.com/yuno-face-recognition:production \
		-f docker/production/lambda-face-recognition/Dockerfile ./server
	docker push 173549642885.dkr.ecr.ap-northeast-1.amazonaws.com/yuno-face-recognition:production

update-resize-lambda:
	aws-vault exec lee-tera -- aws lambda update-function-code --function-name yuno-resize-production --image-uri 173549642885.dkr.ecr.ap-northeast-1.amazonaws.com/yuno-resize:production

update-face-recognition-lambda:
	aws-vault exec lee-tera -- aws lambda update-function-code --function-name yuno-face-recognition-production --image-uri 173549642885.dkr.ecr.ap-northeast-1.amazonaws.com/yuno-face-recognition:production
