build:
	CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ./bin/mammon ./cmd/mammon
run: 
	go run ./cmd/mammon
test:
	go test ./... -cover
docker:
	sudo -E docker-compose up --build