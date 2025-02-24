package dto

type ShareVideoRequest struct {
	Description string `json:"description"`
	UpVote      int64  `json:"upvote"`
	DownVote    int64  `json:"downvote"`
	Thumbnail   string `json:"thumbnail"`
	VideoUrl    string `json:"video_url"`
}

type ShareVideoResponse struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	UpVote      int64  `json:"upvote"`
	DownVote    int64  `json:"downvote"`
	Thumbnail   string `json:"thumbnail"`
	VideoUrl    string `json:"video_url"`
	AccountID   int64  `json:"account_id"`
}

type VideoResponse struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	UpVote      int64  `json:"upvote"`
	DownVote    int64  `json:"downvote"`
	Thumbnail   string `json:"thumbnail"`
	VideoUrl    string `json:"video_url"`
	AccountID   int64  `json:"account_id"`
}
