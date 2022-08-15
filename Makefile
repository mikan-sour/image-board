setup-env:
	cp ./.env.sample ./.env

up:
	docker-compose up --build -d 

down:
	docker-compose down