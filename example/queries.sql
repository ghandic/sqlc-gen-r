-- name: getReplyIds :many
SELECT
    id
FROM
    post
WHERE
    parent_id = ?;

