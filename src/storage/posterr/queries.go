package posterr

const (
	selectAllPosts = `SELECT post_id, username, COALESCE(content, ''), COALESCE(reposted_id, ''), created_at
                 FROM posts
                 ORDER BY created_at DESC
                 LIMIT 10
                 OFFSET $1`

	selectFollowingPosts = `SELECT post_id, username, COALESCE(content, ''), COALESCE(reposted_id, ''), created_at
                 FROM posts
                 WHERE username IN (
                     SELECT username
                     FROM followers
                     WHERE followed_by = $1)
                 ORDER BY created_at DESC
                 LIMIT 10
                 OFFSET $2`

	selectProfilePosts = `SELECT post_id, username, COALESCE(content, ''), COALESCE(reposted_id, ''), created_at
                 FROM posts
                 WHERE username = $1
                 ORDER BY created_at DESC
                 LIMIT 5
                 OFFSET $2`

	countDailyPosts = `SELECT COUNT(*) as daily_posts
                 FROM posts
                 WHERE username = $1
                 AND date_trunc('day', created_at) = date_trunc('day', NOW())`

	searchPosts = `SELECT post_id, username, COALESCE(content, ''), COALESCE(reposted_id, ''), created_at
                 FROM posts
                 WHERE content IS NOT NULL AND content LIKE '%' || $1 || '%'
                 ORDER BY created_at DESC
                 LIMIT $2
                 OFFSET $3`
)
