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

- [ ] Better username validation (\*)
  - 14 chars max
  - Must be at least 3 characters
  - Invalid characters: ! @ # $ % ^ & \* ( ) = +
  - [ { ] } \ |
  - ; :
  - / ?
  - Allowed: . - \_ ' < >
- [x] Group should be private by default
- [ ] Chat
- [ ] Matchmaking
- [ ] Set default filters based on user profile (region=, gamemode=, canjoin=) (\*)
- [ ] Filter or sort characters based on selected roles
- [ ] vice versa, derive roles from characters?
- [x] Remove player from group

API Integration:

- [x] Integrate Find Group Form
- [x] Integrate Create Group Form
- [x] Integrate Profile Form
- [x] Integrate Browse Page
- [ ] Browse Page Pagination & Filtering (\*)
- [x] Integrate Group Page
- [x] Integrate Join Group
- [x] Integrate remove player from group
- [ ] Chat
- [ ] Matchmaking

Bugs:

- [x] Fix flicker on group page?
  - [x] repro'd when joining a public group
- [x] Group table does not show until groups are loaded, should load in but with empty data.
- [ ] After creating a group, the owner should not have to authenticate (sessions integration will probs solve this)
  - repro: create a group, click away, then come back to the group page. this works, but then if you refresh the dialog will open again
