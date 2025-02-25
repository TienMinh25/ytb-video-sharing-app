package dto

type ShareVideoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
	UpVote      int64  `json:"upvote" binding:"omitempty"`
	DownVote    int64  `json:"downvote" binding:"omitempty"`
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
}

type VideoResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UpVote      int64  `json:"upvote"`
	DownVote    int64  `json:"downvote"`
	Thumbnail   string `json:"thumbnail"`
	VideoUrl    string `json:"video_url"`
	SharedBy    string `json:"shared_by"`
}
