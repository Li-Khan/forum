.PHONY: build

build:
	go build -o bin/main ./cmd

.PHONY: docker

run:
	@go run ./cmd/

docker:
	@echo "Building Docker Image:"
	docker image build -f Dockerfile -t forum-image .
	@echo	
	@echo "List of images:"
	docker images
	@echo	
	@echo "Initiating Container:"
	docker container run -t -p 27969:27960 --detach --name forum-container forum-image
	@echo	
	@echo "Running command:"
	docker exec -it forum-container ls -la
	@echo	
	@echo "Running server:"
	docker exec -it forum-container ./main
	@echo

.DEFAULT_GOAL := build