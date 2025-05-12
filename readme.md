# Golang Server

- Auth Server
  - Supports MFA
    - 1FA: username+password pair
    - 2FA: username+password pair, then otp_email
    - Email and Mode configuration

- Boilerplate for future projects.
- Nats Server (WIP)

See [todos](docs/todos)

## local development 
`go run ./cmd/client_demos/nats` - Demonstration of a durable message broker. \
`go run ./cmd/client_demos/aws` - Using AWS S3 SDK.

- http servers: `go run cmd/demos/servers/<pkg>` - 
  - Auth Server `cmd/demos/servers/multifact`
    - Swagger Docs @`cmd/demos/servers/multifact/swagger`
    - Test Coverage @ `cmd/demos/servers/multifact/e2e_test`



# cmd
demos
- `cond`
- `sync`
- `nats`
- `pool`
- `context`