# rivalslfg

## Features

1. Create Group
2. List Groups
   - Filter by:
     -  Region
     -  Gamemode
     -  Requirements Met (provided player info)
     -  Visibility
     -  Platform
     -  Gamemode
     -  Size
     -  Group Settings
       - Platforms
       - Voice Chat
       - Mic
3. Get Group by ID
4. Update Group Info
5. Join Group (if private, authenticate provided passcode)
6. Remove Player from Group

7. Create Player
8. Update Player Info
9.  Authenticate as Player

10. Chat

## Endpoints

- GET /groups/:id
- GET /groups


## Requirements
- Docker
- [sqlc](https://docs.sqlc.dev/en/stable/overview/install.html)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)