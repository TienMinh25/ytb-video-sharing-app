package service

import (
	"context"
	"math"
	"net/http"
	"ytb-video-sharing-app-be/internal/dto"
	"ytb-video-sharing-app-be/internal/entities"
	"ytb-video-sharing-app-be/internal/repository"
)

type VideoService interface {
	ShareVideoYTB(ctx context.Context, payload *entities.Video) (*dto.ShareVideoResponse, *dto.ErrorResponse)

	GetListVideos(ctx context.Context, limit int, page int) ([]*dto.VideoResponse, int, int, bool, bool, *dto.ErrorResponse)
}

type videoServie struct {
	videoRepository repository.VideoRepository
}

func NewVideoService(videoRepository repository.VideoRepository) VideoService {
	return &videoServie{
		videoRepository: videoRepository,
	}
}

func (v videoServie) ShareVideoYTB(ctx context.Context, payload *entities.Video) (*dto.ShareVideoResponse, *dto.ErrorResponse) {
	res, err := v.videoRepository.CreateVideo(ctx, payload)

	if err != nil {
		return nil, &dto.ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return &dto.ShareVideoResponse{
		ID:          res.ID,
		Title:       res.Title,
		Description: res.Description,
		UpVote:      res.UpVote,
		DownVote:    res.DownVote,
		Thumbnail:   res.Thumbnail,
		VideoUrl:    res.VideoUrl,
		AccountID:   res.AccountID,
	}, nil
}

func (v *videoServie) GetListVideos(ctx context.Context, limit int, page int) ([]*dto.VideoResponse, int, int, bool, bool, *dto.ErrorResponse) {
	videos, totalItems, err := v.videoRepository.GetListVideos(ctx, page, limit)
	if err != nil {
		return nil, 0, 0, false, false, &dto.ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	videoResponses := make([]*dto.VideoResponse, 0)
	for _, video := range videos {
		videoResponses = append(videoResponses, &dto.VideoResponse{
			ID:          video.ID,
			Title:       video.Title,
			Description: video.Description,
			UpVote:      video.UpVote,
			DownVote:    video.DownVote,
			Thumbnail:   video.Thumbnail,
			VideoUrl:    video.VideoUrl,
			AccountID:   video.AccountID,
		})
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	isNext := page < totalPages
	isPrevious := page > 1

	return videoResponses, totalItems, totalPages, isNext, isPrevious, nil
}
