.PHONY: build

build:
	go build -o bin/main ./cmd

.PHONY: docker

run:
	@go run ./cmd/

docker:
	docker volume create li-khan-forum
	docker build -t forum .

	@echo "Running server:"
	@echo "\n***************************"
	@echo "*                         *"
	@echo "* http://localhost:27969/ *"
	@echo "*                         *"
	@echo "***************************\n"

	docker run --rm --name web -p 27969:27960 -v li-khan-forum:/app/ forum

docker-delete:
	docker rmi forum

.DEFAULT_GOAL := build