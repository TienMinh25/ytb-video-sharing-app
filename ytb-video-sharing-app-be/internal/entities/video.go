package entities

type Video struct {
	ID          int64  `db:"id"`
	Description string `db:"description"`
	UpVote      int64  `db:"upvote"`
	DownVote    int64  `db:"downvote"`
	Thumbnail   string `db:"thumbnail"`
	VideoUrl    string `db:"video_url"`
	AccountID   int64  `db:"account_id"`
}
