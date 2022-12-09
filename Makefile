
export PROJECT_NAME = blockchain-parser
export PROJECT_PATH = github/$(PROJECT_NAME)

#----------------
# building service
#----------------

build-in-docker:
	go build -o /bin/$(PROJECT_NAME) ./cmd/$(PROJECT_NAME)-server/main.go

build-service-image:
	docker build -t $(PROJECT_NAME) \
		--build-arg PROJECT_NAME=$(PROJECT_NAME) \
		--build-arg PROJECT_PATH=$(PROJECT_PATH) \
		--progress plain \
		-f ./build/Dockerfile .

#----------------
# running service
#----------------

run-in-docker:
	docker-compose -f ./build/dev/docker-compose.yml up -d

stop-in-docker:
	docker-compose -f ./build/dev/docker-compose.yml stop