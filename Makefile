PHONY:
	docker-dev

docker-dev:
	docker-compose --env-file .env.dev up

docker-prod:
	docker-compose --env-file .env.prod up

docker-rebuild:
	docker-compose -d --no-deps --build
