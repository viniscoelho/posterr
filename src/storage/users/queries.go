package users

const (
	selectUser = `SELECT *
                 FROM users
                 WHERE username = $1`

	countFollowers = `SELECT COUNT(*) as followers
                 FROM followers
                 WHERE username = $1`

	countFollowing = `SELECT COUNT(*) as following
                 FROM followers
                 WHERE followed_by = $1`

	listFollowers = `SELECT followed_by
                 FROM followers
                 WHERE username = $1`

	countUserPosts = `SELECT COUNT(*) as no_posts
                 FROM posts
                 WHERE username = $1`

	isFollowerOf = `SELECT COUNT(*) as is_follower
                 FROM followers
                 WHERE username = $1 AND followed_by = $2`
)
