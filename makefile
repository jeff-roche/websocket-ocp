build-client:
	go build -o bin/websocketclient ./cli-client

build-server:
	go build -o bin/websocketserver ./server

build: build-client build-server

run-client: 
	./bin/websocketclient

run-server: 
	./bin/websocketserver

client: build-client run-client

server: build-server run-server
