# rivalslfg.com (backend)

![visitors](https://img.shields.io/endpoint?url=https://vu-mi.com/api/v1/views?id=jcserv/rivalslfg.com/backend)

this folder contains the backend for [rivalslfg](https://rivalslfg.vercel.app/), which is deployed to [fly.io](https://fly.io).

*Generated from [go-api-template](https://github.com/jcserv/go-api-template)*

## installation

### prerequisites/stack
- [go](https://go.dev/doc/install)
- [docker](https://docs.docker.com/get-started/get-docker/)
- [sqlc](https://docs.sqlc.dev/en/stable/overview/install.html)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [gomock](https://github.com/uber-go/mock)
- OPTIONAL:
  - [postgres](https://www.postgresql.org/)
  - [flyctl](https://fly.io/docs/flyctl/installing/)
  - [postico](https://eggerapps.at/postico2/)
  - [godepgraph](https://github.com/kisielk/godepgraph)

### running locally

1. clone the repo
2. `cp .env.example .env`, fill in the required environment variables
3. run `make dev-db`
4. run `make migrate`
5. run `make dev`