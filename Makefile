
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
# run service
#----------------

run-in-docker:
	docker-compose -f ./build/dev/docker-compose.yml up -d

stop-in-docker:
	docker-compose -f ./build/dev/docker-compose.yml stop

#----------------
# run tests
#----------------

run-test-in-docker:
	docker run -it -v $(PWD):$(PWD) -w $(PWD) golang:1.19 go test ./...

#----------------
# generate mocks
#----------------

mocks:
	go generate ./...

#----------------
# generate docs
#----------------

docs:
	docker run --rm -it -p 8001:8001 --user $(shell id -u):$(shell id -g) -e GOPATH=$(shell go env GOPATH):/go -v $(PWD):$(PWD) -w $(PWD) quay.io/goswagger/swagger:v0.30.3 serve \
		-p 8001 \
		--host=0.0.0.0 \
		--no-open \
		./api/swagger.yml
