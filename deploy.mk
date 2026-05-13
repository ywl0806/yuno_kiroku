include .env.production
export

ECR_BASE := $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com

build-api-lambda:
	cd server && \
	GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/lambda/api/main.go && \
	zip function.zip bootstrap && \
	rm bootstrap && \
	aws-vault exec $(AWS_VAULT_PROFILE) -- aws s3 cp function.zip s3://$(LAMBDA_ZIP_BUCKET)/function.zip && \
	rm function.zip

deploy-api-lambda:
	aws-vault exec $(AWS_VAULT_PROFILE) -- aws lambda update-function-code \
		--function-name $(LAMBDA_API_NAME) \
		--s3-bucket $(LAMBDA_ZIP_BUCKET) --s3-key function.zip

build-resize-lambda:
	aws-vault exec $(AWS_VAULT_PROFILE) -- sh -c '\
		aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR_BASE) && \
		docker buildx build --platform linux/arm64 --provenance=false --push \
			-t $(ECR_BASE)/$(ECR_RESIZE_REPO):production \
			-f docker/production/lambda-resize/Dockerfile ./server'

deploy-resize-lambda:
	aws-vault exec $(AWS_VAULT_PROFILE) -- aws lambda update-function-code \
		--function-name $(LAMBDA_RESIZE_NAME) \
		--image-uri $(ECR_BASE)/$(ECR_RESIZE_REPO):production

build-ai-task:
	aws-vault exec $(AWS_VAULT_PROFILE) -- sh -c '\
		aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR_BASE) && \
		docker buildx build --platform linux/amd64 --provenance=false --push \
			-t $(ECR_BASE)/$(ECR_AI_REPO):production \
			-f docker/production/ai/Dockerfile .'

deploy-ai-task:
	aws-vault exec $(AWS_VAULT_PROFILE) -- sh -c '\
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR_BASE) && \
	docker buildx build --platform linux/amd64 --provenance=false --push \
		-t $(ECR_BASE)/$(ECR_AI_REPO):production \
		-f docker/production/ai/Dockerfile .'
	aws-vault exec $(AWS_VAULT_PROFILE) -- aws ecs update-service \
		--cluster $(ECS_CLUSTER) --service $(ECS_AI_SERVICE) --force-new-deployment

deploy-frontend:
	cd front && yarn build:production && \
	aws-vault exec $(AWS_VAULT_PROFILE) -- aws s3 sync dist s3://$(FRONTEND_BUCKET) --delete
	aws-vault exec $(AWS_VAULT_PROFILE) -- aws cloudfront create-invalidation \
		--distribution-id $(CLOUDFRONT_ID) --paths "/*"

build-video-task:
	aws-vault exec $(AWS_VAULT_PROFILE) -- sh -c '\
		aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR_BASE) && \
		docker buildx build --platform linux/amd64 --provenance=false --push \
			-t $(ECR_BASE)/$(ECR_VIDEO_REPO):production \
			-f docker/production/video-processing-worker/Dockerfile ./server'

deploy-video-task:
	aws-vault exec $(AWS_VAULT_PROFILE) -- aws ecs update-service \
		--cluster $(ECS_CLUSTER) --service $(ECS_VIDEO_SERVICE) --force-new-deployment

deploy-production:
	make -f deploy.mk build-ai-task
	make -f deploy.mk deploy-ai-task
	make -f deploy.mk build-video-task
	make -f deploy.mk deploy-video-task
	make -f deploy.mk build-resize-lambda
	make -f deploy.mk deploy-resize-lambda
	make -f deploy.mk build-api-lambda
	make -f deploy.mk deploy-api-lambda
	make -f deploy.mk deploy-frontend