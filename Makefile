VERSION := latest

build:
	docker build -t tonflow:$(VERSION) --target tonflow .

local:
#	docker-compose -f docker-compose-local.yml --env-file .env.local down
	make build
	docker-compose -f docker-compose-local.yml --env-file .env.local up -d
	make log

prod:
#	docker-compose -f docker-compose.yml --env-file .env.prod down
	make build
	docker-compose -f docker-compose.yml --env-file .env.prod up -d

log:
	docker logs --follow tonflow