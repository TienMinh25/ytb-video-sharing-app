package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"ytb-video-sharing-app-be/internal/dto"
	"ytb-video-sharing-app-be/internal/entities"
	"ytb-video-sharing-app-be/internal/service"
	"ytb-video-sharing-app-be/internal/websock"
	"ytb-video-sharing-app-be/utils"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	videoService  service.VideoService
	// messageBroker pkg.Queue
	wsManager     *websock.Manager
}

func NewVideoHandler(videoService service.VideoService, wsManager *websock.Manager) *VideoHandler {
	return &VideoHandler{
		videoService:  videoService,
		// messageBroker: messageBroker,
		wsManager:     wsManager,
	}
}

// ShareVideo godoc
//
//	@Summary		Share new video
//	@Tags			videos
//	@Description	Create new video and return itself.
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Param			conn_id	query		string					true	"WebSocket connection ID"
//	@Param			request	body		dto.ShareVideoRequest	true	"Share video payload"
//	@Success		201		{object}	dto.ShareVideoResponseDocs
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		500		{object}	dto.ResponseError
//	@Router			/videos [post]
func (v *VideoHandler) ShareVideo(ctx *gin.Context) {
	claimsStr, _ := ctx.Get("claims")
	claims := claimsStr.(*utils.UserClaims)

	req, _ := ctx.Get("data")

	// TODO: get connid of connection websocket in here
	connID := ctx.Query("conn_id")

	if connID == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Missing conn_id in query")
		return
	}

	data := req.(dto.ShareVideoRequest)

	// call service to share video
	res, err := v.videoService.ShareVideoYTB(ctx, &entities.Video{
		AccountID:   claims.AccountID,
		Description: data.Description,
		Title:       data.Title,
		UpVote:      data.UpVote,
		DownVote:    data.DownVote,
		Thumbnail:   data.Thumbnail,
		VideoUrl:    data.VideoUrl,
	})

	if err != nil {
		utils.ErrorResponse(ctx, err.Code, err)
		return
	}

	// push event to kafka
	// go func() {
	// 	payloadKafka := &third_party.VideoMessageEvent{
	// 		AccountId: claims.AccountID,
	// 		Title:     data.Title,
	// 		Thumbnail: data.Thumbnail,
	// 		SharedBy:  claims.Email,
	// 	}

	// 	payloadBytes, _ := third_party.SerializeVideoMessageEvent(payloadKafka)

	// 	v.messageBroker.Produce(os.Getenv("KAFKA_TOPIC"), payloadBytes)
	// }()

	// send through websocket
	// go func() {
	if connID != "" {
		newEvent := websock.EventNotificationMessage{
			Title:     data.Title,
			SharedBy:  claims.Email,
			Thumbnail: data.Thumbnail,
		}

		payload, err := json.Marshal(newEvent)

		if err != nil {
			log.Println("error when marshaling json: ", err)
			return
		}

		v.wsManager.SendBroadCast(websock.Event{
			Type:    "new_video",
			Payload: payload,
		}, connID)
	}
	// }()

	utils.SuccessResponse(ctx, http.StatusCreated, res)
}

// GetListVideos GetVideos godoc
//
//	@Summary		Get list videos
//	@Tags			videos
//	@Description	Get list videos
//	@Accept			json
//	@Produce		json
//
//	@Param			limit	query		int	true	"Limit number of records returned"
//	@Param			page	query		int	true	"page"
//	@Success		200		{object}	dto.ListVideosResponseDocs
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		500		{object}	dto.ResponseError
//	@Router			/videos [get]
func (v *VideoHandler) GetListVideos(ctx *gin.Context) {
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil || limit <= 0 {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid page parameter")
		return
	}

	// Call service to get videos
	res, totalItems, totalPages, isNext, isPrevious, errRes := v.videoService.GetListVideos(ctx, limit, page)
	if errRes != nil {
		utils.ErrorResponse(ctx, errRes.Code, errRes)
		return
	}

	utils.PaginatedResponse(ctx, res, page, limit, totalPages, totalItems, isNext, isPrevious)
}
