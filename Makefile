up:
	docker-compose up loadtest

down:
	docker-compose down --remove-orphans

test:
	docker-compose run unittest