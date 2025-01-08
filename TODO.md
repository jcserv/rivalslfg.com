# TODO

- [ ] Update platforms to PC/Console
- [ ] Add new rank "celestial" inbetween gm and eternity
- [ ] Add new characters Mr. Fantastic (Duelist) and Invisible Woman (Strategist)

## frontend

FE-Only:

- [x] Better username validation (\*)
  - 14 chars max
  - Must be at least 3 characters
  - Invalid characters: ! @ # $ % ^ & \* ( ) = +
  - [ { ] } \ |
  - ; :
  - / ?
  - Allowed: . - \_ ' < >
- [x] Group should be private by default
- [x] Chat
- [/] Matchmaking
- [ ] Set default filters based on user profile (region=, gamemode=, canjoin=) (\*)
- [x] Remove player from group
- [X] Replace hard-coded player IDs
- [ ] Fix responsiveness of multi-select & dark mode

API Integration:

- [x] Integrate Find Group Form
- [x] Integrate Create Group Form
- [x] Integrate Profile Form
- [x] Integrate Browse Page
- [x] Browse Page Pagination & Filtering (\*)
- [x] Integrate Group Page
- [x] Integrate Join Group
- [x] Integrate remove player from group
- [ ] Chat
- [ ] Matchmaking
- [ ] Pagination: Server-side filtering
  - [X] open
  - [ ] requirements met

Bugs:

- [x] Fix flicker on group page?
  - [x] repro'd when joining a public group
- [x] Group table does not show until groups are loaded, should load in but with empty data.
- [X] After creating a group, the owner should not have to authenticate (sessions integration will probs solve this)
  - repro: create a group, click away, then come back to the group page. this works, but then if you refresh the dialog will open again

## backend

1. [X] List Groups
   - [X] Filter by: (*)
     -  Region `/v1/groups?filter=region eq "na"`
     -  Gamemode `/v1/groups?filter=gamemode eq "competitive"`
     -  Requirements Met (provided player info) `/v1/groups?filter=requirementsMet eq true`

     -  [X] Visibility
     -  Platform
     -  Gamemode
     -  Size
     -  Group Settings
       - Platforms
       - Voice Chat
       - Mic
2. [X] Upsert Group
3. [ ] Delete Group
4. [X] Join Group (if private, authenticate provided passcode)
5. [X] Remove Player from Group
6. [X] Leave Group
7. [X] Get Group Passcode

8. Chat

9. Matchmaking
   - Find groups the user can join
     - Prioritize groups that are close to completion
   - Acquire lock on group
   - Join
   - Release lock
   - If no groups are found, create a new group with as many queued players as possible?

Bugs:
- [X] No auth right now, so users can modify other users' info if they know their id
- [X] prevent group creation if user already has group
- supported filters/sorters for ParseQueryParams
- [X] Unable to leave your own group if you are the only member
  - v1/players.go:124	can't scan into dest[0]: cannot scan NULL into *int32
- [X] Can get into a state where you can't leave a group

Before release:
- [ ] Squash migrations into one, remove seed data

* indicates these are good first issues

tweaks:
- [X] add group count field to group table? (can be derived from players json but not sure if that's a performance issue)
