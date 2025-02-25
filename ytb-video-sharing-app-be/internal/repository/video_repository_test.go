package repository

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"ytb-video-sharing-app-be/internal/entities"
)

type videoConfig struct {
	testConfig
	repo VideoRepository
}

func SetupVideoConfig(t *testing.T) *videoConfig {
	testConf := SetupTest(t)

	return &videoConfig{
		testConfig: *testConf,
		repo:       NewVideoRepository(testConf.db),
	}
}

func TestCreateVideo(t *testing.T) {
	t.Run("Should create a video successfully", func(t *testing.T) {
		cfg := SetupVideoConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		video := &entities.Video{
			ID:          1,
			Title:       "test",
			Description: "Test Video",
			UpVote:      10,
			DownVote:    2,
			Thumbnail:   "https://thumbnail.url",
			VideoUrl:    "https://video.url",
			AccountID:   1,
		}

		mockSqlResult := &MockSQLResult{
			LastInsertID: video.ID,
			RowAffected:  1,
		}
		cfg.db.EXPECT().ExecWithResult(ctx, gomock.Any(), video.Title, video.Description, video.UpVote, video.DownVote, video.Thumbnail, video.VideoUrl, video.AccountID).Return(mockSqlResult, nil)
		cfg.db.EXPECT().QueryRow(ctx, gomock.Any(), mockSqlResult.LastInsertID).Return(cfg.row)

		cfg.row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(args ...interface{}) error {
				*args[0].(*int64) = video.ID
				*args[1].(*string) = video.Title
				*args[2].(*string) = video.Description
				*args[3].(*int64) = video.UpVote
				*args[4].(*int64) = video.DownVote
				*args[5].(*string) = video.Thumbnail
				*args[6].(*string) = video.VideoUrl
				*args[7].(*int64) = video.AccountID

				return nil
			})

		rs, err := cfg.repo.CreateVideo(ctx, video)

		assert.NoError(t, err)
		assert.Equal(t, video, rs)
	})

	t.Run("Should return error when DB execution fails", func(t *testing.T) {
		cfg := SetupVideoConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		video := &entities.Video{
			Title:       "test",
			Description: "Test Video",
			UpVote:      10,
			DownVote:    2,
			Thumbnail:   "https://thumbnail.url",
			VideoUrl:    "https://video.url",
			AccountID:   1,
		}

		expectedErr := errors.New("db execution failed")
		cfg.db.EXPECT().ExecWithResult(ctx, gomock.Any(), video.Title, video.Description, video.UpVote, video.DownVote, video.Thumbnail, video.VideoUrl, video.AccountID).Return(nil, expectedErr)

		_, err := cfg.repo.CreateVideo(ctx, video)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Should return error when failing to get last insert ID", func(t *testing.T) {
		cfg := SetupVideoConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		video := &entities.Video{
			Title:       "test",
			Description: "Test Video",
			UpVote:      10,
			DownVote:    2,
			Thumbnail:   "https://thumbnail.url",
			VideoUrl:    "https://video.url",
			AccountID:   1,
		}

		expectedErr := errors.New("db execution failed")
		mockSqlResult := &MockSQLResult{}
		cfg.db.EXPECT().ExecWithResult(ctx, gomock.Any(), video.Title, video.Description, video.UpVote, video.DownVote, video.Thumbnail, video.VideoUrl, video.AccountID).Return(mockSqlResult, nil)

		_, err := cfg.repo.CreateVideo(ctx, video)
		assert.Error(t, err)
		assert.Equal(t, expectedErr.Error(), err.Error())
	})
}

func TestGetVideo(t *testing.T) {
	t.Run("Should return video when found", func(t *testing.T) {
		cfg := SetupVideoConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		expectedVideo := &entities.Video{
			ID:          1,
			Title:       "test",
			Description: "Test Video",
			UpVote:      10,
			DownVote:    2,
			Thumbnail:   "https://thumbnail.url",
			VideoUrl:    "https://video.url",
			AccountID:   1,
		}

		cfg.db.EXPECT().QueryRow(ctx, gomock.Any(), expectedVideo.ID).Return(cfg.row)

		cfg.row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(args ...interface{}) error {
			*args[0].(*int64) = expectedVideo.ID
			*args[1].(*string) = expectedVideo.Title
			*args[2].(*string) = expectedVideo.Description
			*args[3].(*int64) = expectedVideo.UpVote
			*args[4].(*int64) = expectedVideo.DownVote
			*args[5].(*string) = expectedVideo.Thumbnail
			*args[6].(*string) = expectedVideo.VideoUrl
			*args[7].(*int64) = expectedVideo.AccountID
			return nil
		})

		video, err := cfg.repo.GetVideo(ctx, expectedVideo.ID)
		assert.NoError(t, err)
		assert.Equal(t, expectedVideo, video)
	})

	t.Run("Should return error if video not found", func(t *testing.T) {
		cfg := SetupVideoConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		cfg.db.EXPECT().QueryRow(ctx, gomock.Any(), gomock.Any()).Return(cfg.row)

		cfg.row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("no rows"))

		video, err := cfg.repo.GetVideo(ctx, 1)
		assert.Nil(t, video)
		assert.Error(t, err)
	})

	t.Run("Should return error if scan fails", func(t *testing.T) {
		cfg := SetupVideoConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		cfg.db.EXPECT().QueryRow(ctx, gomock.Any(), gomock.Any()).Return(cfg.row)

		err := errors.New("scan error")
		cfg.row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(err)

		result, errRes := cfg.repo.GetVideo(ctx, 1)

		assert.Nil(t, result)
		assert.Error(t, errRes)
		assert.Equal(t, errRes, err)
	})
}

func TestGetListVideos(t *testing.T) {
	t.Run("Should return list of videos", func(t *testing.T) {
		cfg := SetupVideoConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		expectedVideos := []*entities.Video{
			{ID: 1, Title: "test 1", Description: "Video 1", UpVote: 5, DownVote: 1, Thumbnail: "thumb1.jpg", VideoUrl: "url1", AccountID: 1},
			{ID: 2, Title: "test 2", Description: "Video 2", UpVote: 3, DownVote: 0, Thumbnail: "thumb2.jpg", VideoUrl: "url2", AccountID: 2},
		}

		cfg.db.EXPECT().QueryRow(ctx, gomock.Any()).Return(cfg.row)
		cfg.row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
			*args[0].(*int) = len(expectedVideos)
			return nil
		})

		cfg.db.EXPECT().Query(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(cfg.rows, nil)

		cfg.rows.EXPECT().Next().Return(true).Times(len(expectedVideos))
		cfg.rows.EXPECT().Next().Return(false).Times(1)
		cfg.rows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(args ...interface{}) error {
			video := expectedVideos[0]
			expectedVideos = expectedVideos[1:]
			*args[0].(*int64) = video.ID
			*args[1].(*string) = video.Title
			*args[2].(*string) = video.Description
			*args[3].(*int64) = video.UpVote
			*args[4].(*int64) = video.DownVote
			*args[5].(*string) = video.Thumbnail
			*args[6].(*string) = video.VideoUrl
			*args[7].(*int64) = video.AccountID
			return nil
		}).Times(len(expectedVideos))

		cfg.rows.EXPECT().Close().Times(1)

		videos, total, err := cfg.repo.GetListVideos(ctx, 1, 2)
		assert.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, videos, 2)
	})

	t.Run("Should return list <nil>, total pages 0 and error", func(t *testing.T) {
		cfg := SetupVideoConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()

		err := errors.New("no rows")
		cfg.db.EXPECT().QueryRow(ctx, gomock.Any()).Return(cfg.row)
		cfg.row.EXPECT().Scan(gomock.Any()).Return(err)

		videos, total, errRes := cfg.repo.GetListVideos(ctx, 1, 2)

		assert.Error(t, err)
		assert.Equal(t, err, errRes)
		assert.Nil(t, videos)
		assert.Equal(t, 0, total)
	})

	t.Run("Should return error because db execute failed", func(t *testing.T) {
		cfg := SetupVideoConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		expectedVideos := []*entities.Video{
			{ID: 1, Title: "test 1", Description: "Video 1", UpVote: 5, DownVote: 1, Thumbnail: "thumb1.jpg", VideoUrl: "url1", AccountID: 1},
			{ID: 2, Title: "test 2", Description: "Video 2", UpVote: 3, DownVote: 0, Thumbnail: "thumb2.jpg", VideoUrl: "url2", AccountID: 2},
		}

		cfg.db.EXPECT().QueryRow(ctx, gomock.Any()).Return(cfg.row)
		cfg.row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
			*args[0].(*int) = len(expectedVideos)
			return nil
		})

		err := errors.New("db execute failed")
		cfg.db.EXPECT().Query(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, err)

		videos, total, errRes := cfg.repo.GetListVideos(ctx, 1, 2)

		assert.Error(t, err)
		assert.Equal(t, err, errRes)
		assert.Nil(t, videos)
		assert.Equal(t, 0, total)
	})

	t.Run("Should return error because cannot scan one record", func(t *testing.T) {
		cfg := SetupVideoConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		expectedVideos := []*entities.Video{
			{ID: 1, Title: "test 1", Description: "Video 1", UpVote: 5, DownVote: 1, Thumbnail: "thumb1.jpg", VideoUrl: "url1", AccountID: 1},
			{ID: 2, Title: "test 2", Description: "Video 2", UpVote: 3, DownVote: 0, Thumbnail: "thumb2.jpg", VideoUrl: "url2", AccountID: 2},
		}

		cfg.db.EXPECT().QueryRow(ctx, gomock.Any()).Return(cfg.row)
		cfg.row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
			*args[0].(*int) = len(expectedVideos)
			return nil
		})

		cfg.db.EXPECT().Query(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(cfg.rows, nil)

		err := errors.New("scan error")
		cfg.rows.EXPECT().Next().Return(true).Times(len(expectedVideos))
		cfg.rows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(args ...interface{}) error {
			video := expectedVideos[0]
			expectedVideos = expectedVideos[1:]
			*args[0].(*int64) = video.ID
			*args[1].(*string) = video.Title
			*args[2].(*string) = video.Description
			*args[3].(*int64) = video.UpVote
			*args[4].(*int64) = video.DownVote
			*args[5].(*string) = video.Thumbnail
			*args[6].(*string) = video.VideoUrl
			*args[7].(*int64) = video.AccountID
			return nil
		}).Times(len(expectedVideos) - 1)
		cfg.rows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(err).Times(1)

		cfg.rows.EXPECT().Close().Times(1)

		videos, total, errRes := cfg.repo.GetListVideos(ctx, 1, 2)

		assert.Error(t, err)
		assert.Equal(t, err, errRes)
		assert.Nil(t, videos)
		assert.Equal(t, 0, total)
	})
}
