# [Confluent Doc](https://github.com/confluentinc/confluent-kafka-go#using-go-modules)
# Remember:
# If you are building for Alpine Linux (musl), -tags musl must be specified.
# go build -tags musl ./...

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

run:
	go run cmd/franz/main.go

publish-random:
	@echo "Run:\n\tgo run cmd/publisher/main.go [topic_name] [number_of_messages]\n"
