-- name: FindAllGroups :many
SELECT * FROM Groups;

-- name: FindAllPlayers :many
SELECT * FROM Players;

-- name: GetGroupByID :one
SELECT * 
FROM Groups
JOIN GroupMembers ON Groups.id = GroupMembers.group_id
JOIN Players ON Players.id = GroupMembers.player_id
WHERE Groups.id = $1;
