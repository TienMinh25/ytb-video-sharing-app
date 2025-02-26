package repository

import (
	"context"
	"errors"
	"ytb-video-sharing-app-be/internal/entities"
	"ytb-video-sharing-app-be/pkg"
)

type VideoRepository interface {
	CreateVideo(ctx context.Context, payload *entities.Video) (*entities.Video, error)
	GetVideo(ctx context.Context, videoID int64) (*entities.Video, error)
	GetListVideos(ctx context.Context, page, limit int) ([]*entities.Video, int, error)
}

type videoRepository struct {
	db pkg.Database
}

func NewVideoRepository(db pkg.Database) VideoRepository {
	return &videoRepository{
		db: db,
	}
}

// CreateVideo implements VideoRepository.
func (v *videoRepository) CreateVideo(ctx context.Context, payload *entities.Video) (*entities.Video, error) {
	query := `INSERT INTO videos (title, description, upvote, downvote, thumbnail, video_url, account_id)
			VALUES (?, ?, ?, ?, ?, ?, ?)`

	rs, err := v.db.ExecWithResult(ctx, query, payload.Title, payload.Description, payload.UpVote, payload.DownVote, payload.Thumbnail, payload.VideoUrl, payload.AccountID)

	if err != nil {
		return nil, err
	}

	lastInsertId, err := rs.LastInsertId()

	if err != nil || lastInsertId == 0 {
		return nil, errors.New("db execution failed")
	}

	return v.GetVideo(ctx, lastInsertId)
}

// GetVideo implements VideoRepository.
func (v *videoRepository) GetVideo(ctx context.Context, videoID int64) (*entities.Video, error) {
	query := `SELECT * FROM videos WHERE id = ?`

	video := &entities.Video{}

	if err := v.db.QueryRow(ctx, query, videoID).
		Scan(&video.ID, &video.Title, &video.Description, &video.UpVote, &video.DownVote, &video.Thumbnail, &video.VideoUrl, &video.AccountID); err != nil {
		return nil, err
	}

	return video, nil
}

// GetListVideos implements VideoRepository.
func (v *videoRepository) GetListVideos(ctx context.Context, page, limit int) ([]*entities.Video, int, error) {
	var totalItems int
	countQuery := `SELECT COUNT(*) FROM videos`
	if err := v.db.QueryRow(ctx, countQuery).Scan(&totalItems); err != nil {
		return nil, 0, err
	}

	query := `SELECT v.id, v.title, v.description, v.upvote, v.downvote, v.thumbnail, v.video_url, v.account_id, a.fullname
	FROM videos v
	JOIN accounts a ON v.account_id = a.id
	ORDER BY v.id ASC LIMIT ? OFFSET ?`

	rows, err := v.db.Query(ctx, query, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var videos []*entities.Video
	for rows.Next() {
		video := &entities.Video{}
		if err := rows.Scan(&video.ID, &video.Title, &video.Description, &video.UpVote, &video.DownVote, &video.Thumbnail, &video.VideoUrl, &video.AccountID, &video.FullName); err != nil {
			return nil, 0, err
		}
		videos = append(videos, video)
	}

	return videos, totalItems, nil
}
