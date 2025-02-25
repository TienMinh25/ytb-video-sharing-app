package dto

type ShareVideoRequest struct {
	Title       string `json:"title,required" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
	UpVote      int64  `json:"upvote" binding:"required"`
	DownVote    int64  `json:"downvote" binding:"required"`
	Thumbnail   string `json:"thumbnail" binding:"required"`
	VideoUrl    string `json:"video_url" binding:"required,url"`
}

type ShareVideoResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UpVote      int64  `json:"upvote"`
	DownVote    int64  `json:"downvote"`
	Thumbnail   string `json:"thumbnail"`
	VideoUrl    string `json:"video_url"`
	AccountID   int64  `json:"account_id"`
}

type VideoResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UpVote      int64  `json:"upvote"`
	DownVote    int64  `json:"downvote"`
	Thumbnail   string `json:"thumbnail"`
	VideoUrl    string `json:"video_url"`
	AccountID   int64  `json:"account_id"`
}
