# websocket-ocp
A basic websocket chat server and a CLI client

## Getting Started

Build the Client and Server:
`$ make build`

Build and run the server:
`$ make server`

Build and run the CLI client
`$ make client`

## Server
The server doesn't need any input on startup but you can provide it an optional `addr` input to change the port (default 8080)

Every client is given an id (uuid) which will be prefixed on the message.

A debug UI can be found at `http://0.0.0.0:8080/debug`

Websocket connections can be made to `ws://0.0.0.0:8080/ws`

### CLI Options
- `addr`: The address to make the server available on (default `":8080"`)
### Available commands:
- Build and run: `$ make server`
- Build: `$ make build-server`
- Run: `$ make run-server`

## Client
On startup the client will try to connect to the server. Once connection is successful you will immediately start receiving messages and can type messages into stdin

### CLI Options
- `addr`: the address to connect to (default: `"127.0.0.1:8080"`)
    - Exclude the `ws://` prefix
- `path`: the path on the server to connect to (default: `"/ws"`)

### Available commands:
- `$ make server`
- `$ make build-server`
- `$ make run-server`

## Go Workspace Footnote
I know the `go.work` file isn't recommended to be committed to the repo but It's small, saves time and most importantly, **I'm a rebel.**
