run:
	docker-compose --env-file ./.env up --build

db:
	docker-compose -f docker-compose-test.yml up 

test:
	go get ./...
	go test -p 1 -v -cover -race -timeout 30s -coverprofile=coverage.out -a ./... 
	go tool cover -func=coverage.out

clean: 
	docker-compose -f docker-compose-test.yml down --remove-orphans

stop:
	docker-compose down