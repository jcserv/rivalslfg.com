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
- OPTIONAL:
  - [postgres](https://www.postgresql.org/)
  - [flyctl](https://fly.io/docs/flyctl/installing/)

### running locally

1. clone the repo
2. `cp .env.example .env`, fill in the required environment variables
3. run `make dev-db`
4. run `make migrate`
5. run `make dev`

## todo

1. [X] List Groups
   - [ ] Filter by: (*)
     -  Region `/v1/groups?filter=region eq "na"`
     -  Gamemode `/v1/groups?filter=gamemode eq "competitive"`
     -  Requirements Met (provided player info) `/v1/groups?filter=requirementsMet eq true`

     -  Visibility
     -  Platform
     -  Gamemode
     -  Size
     -  Group Settings
       - Platforms
       - Voice Chat
       - Mic
2. [X] Upsert Group
3. [ ] Delete Group
4. [ ] Join Group (if private, authenticate provided passcode)
5. [ ] Remove Player from Group
6. [ ] Leave Group

7. Chat

8. Matchmaking
   - Find groups the user can join
     - Prioritize groups that are close to completion
   - Acquire lock on group
   - Join
   - Release lock
   - If no groups are found, create a new group with as many queued players as possible?

Bugs:
- No auth right now, so users can modify other users' info if they know their id

Before release:
- [ ] Squash migrations into one, remove seed data

* indicates these are good first issues

tweaks:
- add group count field to group table? (can be derived from players json but not sure if that's a performance issue)