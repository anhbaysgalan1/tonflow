VERSION := latest
build:
	docker build -t tonflow:$(VERSION) --target tonflow .
up:
	docker-compose up -d
down:
	docker-compose down
log:
	docker logs --follow tonflow
rebuild:
	make build
	docker-compose --env-file local.env up