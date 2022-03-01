.PHONY: build

build:
	go build -o bin/main ./cmd

.PHONY: docker

run:
	@go run ./cmd/

docker:
	docker volume create web
	docker build -t forum .
	docker run --rm --name web -p 27969:27960 -v web:/app/ forum

	@echo "Running server:"
	@echo "\n\t***************************"
	@echo "\t* http://localhost:27969/ *"
	@echo "\t***************************\n"

docker-delete:
	docker rmi forum

.DEFAULT_GOAL := build