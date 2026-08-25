package request

import (
	"hub/internal/websocket"
	"hub/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestController struct {
	requestService RequestService
	hub            *websocket.Hub
}

func NewRequestController(requestService RequestService, hub *websocket.Hub) *RequestController {
	return &RequestController{requestService: requestService, hub: hub}
}

// GetAll 			godoc
// @Security 		BearerAuth
// @Summary 		Get all friend requests
// @Description 	Get all friend requests
// @Tags 			Request
// @Accept 			json
// @Produce 		json
// @Success 		200 {object} []RequestResponse
// @Router 			/requests [get]
func (r *RequestController) GetAll(ctx *gin.Context) {
	res, err := r.requestService.GetAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "Lấy danh sách lời mời kết bạn thành công", res)
}

// SendRequest 		godoc
// @Security 		BearerAuth
// @Summary 		Send a friend request
// @Description 	Send a friend request
// @Tags 			Request
// @Accept 			json
// @Produce 		json
// @Param 			request body SendRequestReq true "Request"
// @Success 		200 {object} SendRequestRes
// @Router 			/requests [post]
func (r *RequestController) SendRequest(ctx *gin.Context) {
	req := &SendRequestReq{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, _ := ctx.Get("userId")
	req.SenderId = userId.(uuid.UUID)

	res, err := r.requestService.SendRequest(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	msg := &websocket.NotificationMessage{
		TargetUserID: res.ReceiverId.String(),
		Data:         []byte(res.SenderId.String() + " đã gửi lời mời kết bạn đến bạn"),
	}
	r.hub.SendToUser <- msg

	response.Success(ctx, http.StatusOK, "Gửi lời mời kết bạn thành công", res)
}

// UpdateRequest 			godoc
// @Security 				BearerAuth
// @Summary 				Update a friend request
// @Description 			Update a friend request
// @Tags 					Request
// @Accept 					json
// @Produce 				json
// @Param 					id path string true "Request ID"
// @Param 					request body UpdateRequestReq true "Request"
// @Success 				200 {object} SendRequestRes
// @Router 					/requests/{id} [patch]
func (r *RequestController) UpdateRequest(ctx *gin.Context) {
	requestIdParam := ctx.Param("id")
	requestId, err := uuid.Parse(requestIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid Request ID format")
		return
	}

	req := &UpdateRequestReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, _ := ctx.Get("userId")
	req.RequestId = requestId
	req.ReceiverId = userId.(uuid.UUID)

	res, err := r.requestService.UpdateRequest(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	msg := &websocket.NotificationMessage{
		TargetUserID: res.SenderId.String(),
		Data:         []byte(res.ReceiverId.String() + " đã chấp nhận lời mời kết bạn"),
	}
	r.hub.SendToUser <- msg
	response.Success(ctx, http.StatusOK, "Cập nhật lời mời kết bạn thành công", res)
}

// DeleteRequest 	godoc
// @Security 		BearerAuth
// @Summary 		Delete the request
// @Description 	Delete the request
// @Tags 			Request
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Request ID"
// @Success 		200 {object} nil
// @Router 			/requests/{id} [delete]
func (r *RequestController) DeleteRequest(ctx *gin.Context) {
	requestIdParam := ctx.Param("id")
	requestId, err := uuid.Parse(requestIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid Request ID format")
		return
	}

	userId, _ := ctx.Get("userId")

	err = r.requestService.DeleteRequest(ctx.Request.Context(), requestId, userId.(uuid.UUID))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "Xóa lời mời kết bạn thành công", nil)
}
