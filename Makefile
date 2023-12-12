run:
	docker-compose --env-file ./.env up --build

db:
	docker-compose up postgres 

test:
	go get ./...
	go test -p 1 -v -cover -race -timeout 30s -coverprofile=coverage.out -a ./... 
	go tool cover -func=coverage.out

stop:
	docker-compose down