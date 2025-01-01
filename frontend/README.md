# rivalslfg.com (frontend)

![visitors](https://img.shields.io/endpoint?url=https://vu-mi.com/api/v1/views?id=jcserv/rivalslfg.com/frontend)

this folder contains the frontend for [rivalslfg](https://rivalslfg.vercel.app/), which is deployed to [vercel](https://vercel.com).

_Generated from [react-vite-shadcn-template](https://github.com/jcserv/react-vite-shadcn-template)_

## setup

### prerequisites

- [node](https://nodejs.org/en)
- [pnpm](https://pnpm.io/installation)

### installation

1. clone the repo
2. `cp .env.example .env`, fill in the required environment variables
3. run `pnpm install`
4. run `pnpm run dev`

## references

- [favicon.io/](https://favicon.io/)

### TODO

FE-Only:

- [ ] Better username validation
  - 14 chars max
  - Must be at least 3 characters
  - Invalid characters: ! @ # $ % ^ & \* ( ) = +
  - [ { ] } \ |
  - ; :
  - / ?
  - Allowed: . - \_ ' < >
- Group should be private by default
- [ ] Chat
- [ ] Matchmaking
- [ ] Set default filters based on user profile (region=, gamemode=, canjoin=)

API Integration:

- [x] Integrate Find Group Form
- [ ] Integrate Create Group Form
- [x] Integrate Profile Form
- [x] Integrate Browse Page
- [ ] Browse Page Pagination & Filtering
- [x] Integrate Group Page
- [ ] Integrate Join Group
- [ ] Chat
- [ ] Matchmaking

Bugs:

- [ ] Fix flicker on group page
- [ ] Group table does not show until groups are loaded, should load in but with empty data.
