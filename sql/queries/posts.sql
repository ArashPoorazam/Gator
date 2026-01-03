-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, published_at, title, url, description, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetPostsForUser :many
SELECT * FROM posts
WHERE feed_id = $1
ORDER BY published_at ASC
LIMIT $2;

-- name: ClearPosts :exec
TRUNCATE TABLE posts RESTART IDENTITY CASCADE;