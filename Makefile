VERSION := latest

build:
	docker build -t tonflow:$(VERSION) --target tonflow .

local:
#	docker-compose -f docker-compose-local.yml --env-file local.env down
	make build
	docker-compose -f docker-compose-local.yml --env-file local.env up -d
	make log

prod:
#	docker-compose -f docker-compose.yml --env-file prod.env down
	make build
	docker-compose -f docker-compose.yml --env-file prod.env up -d

log:
	docker logs --follow tonflow