# [Confluent Doc](https://github.com/confluentinc/confluent-kafka-go#using-go-modules)
# Remember:
# If you are building for Alpine Linux (musl), -tags musl must be specified.
# go build -tags musl ./...

up:
	docker-compose up -d

down:
	docker-compose down
