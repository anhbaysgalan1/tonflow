local:
	docker-compose --env-file .env.local -f docker-compose.local.yml up -d --build
	go run cmd/main.go