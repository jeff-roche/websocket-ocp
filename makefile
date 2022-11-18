build-client:
	CGO_ENABLED=0 go build -o bin/websocketclient ./cli-client

build-server:
	CGO_ENABLED=0 go build -o bin/websocketserver ./server

build: build-client build-server

run-client: 
	./bin/websocketclient

run-server: 
	./bin/websocketserver

client: build-client run-client

server: build-server run-server

container:
	docker build -t websocketserver:latest .