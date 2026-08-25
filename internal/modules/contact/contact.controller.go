package contact

import (
	"hub/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContactController struct {
	contactService ContactService
}

func NewContactController(contactService ContactService) *ContactController {
	return &ContactController{contactService: contactService}
}

// GetAll 			godoc
// @Security 		BearerAuth
// @Summary 		Get all contacts
// @Description 	Get all contacts with pagination
// @Tags 			Contact
// @Accept 			json
// @Produce 		json
// @Param 			userName query string false "User Name"
// @Param 			limit query int false "Limit"
// @Param 			offset query int false "Offset"
// @Param 			sortBy query string false "Sort By"
// @Param 			sortOrder query string false "Sort Order"
// @Success 		200 {object} PaginatedContactResponse
// @Router 			/contacts [get]
func (c *ContactController) GetAll(ctx *gin.Context) {
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, _ := ctx.Get("userId")
	params.UserId = userId.(uuid.UUID).String()

	contacts, total, err := c.contactService.GetAll(ctx.Request.Context(), params)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Lấy danh bạ thành công", PaginatedContactResponse{
		Data:  contacts,
		Total: total,
	})
}

// UpdateContact 	godoc
// @Security 		BearerAuth
// @Summary 		Update contact (Block/Unblock)
// @Description 	Update contact status
// @Tags 			Contact
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Contact ID"
// @Param 			request body UpdateContact true "Update info"
// @Success 		200 {object} nil
// @Router 			/contacts/{id} [patch]
func (c *ContactController) UpdateContact(ctx *gin.Context) {
	contactIdParam := ctx.Param("id")
	contactId, err := uuid.Parse(contactIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid Contact ID format")
		return
	}

	req := &UpdateContact{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, _ := ctx.Get("userId")
	req.ContactId = contactId
	req.UserId = userId.(uuid.UUID)

	contactRes, err := c.contactService.UpdateContact(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "Cập nhật danh bạ thành công", contactRes)
}

// DeleteContact 	godoc
// @Security 		BearerAuth
// @Summary 		Delete contact (Unfriend)
// @Description 	Remove someone from contact list
// @Tags 			Contact
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Contact ID"
// @Success 		200 {object} nil
// @Router 			/contacts/{id} [delete]
func (c *ContactController) DeleteContact(ctx *gin.Context) {
	contactIdParam := ctx.Param("id")
	friendId, err := uuid.Parse(contactIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid Contact ID format")
		return
	}

	userIdStr, _ := ctx.Get("userId")
	userId := userIdStr.(uuid.UUID)

	err = c.contactService.DeleteContact(ctx.Request.Context(), userId, friendId)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "Xóa liên hệ thành công", nil)
}
