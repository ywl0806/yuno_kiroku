package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/services"
)

type InviteHandler struct {
	inviteService *services.InviteService
	frontURL      string
}

func NewInviteHandler(inviteService *services.InviteService) *InviteHandler {
	frontURL := viper.GetString("FRONT_URL")
	if frontURL == "" {
		frontURL = "http://localhost:5173"
	}
	return &InviteHandler{inviteService: inviteService, frontURL: frontURL}
}

// CreateInviteRequest 초대 토큰 생성 요청
type CreateInviteRequest struct {
	GroupID int32 `json:"group_id" validate:"required"`
}

// @Description 초대 토큰 생성
//
// @Summary Create Invite Token
// @Tags Invite
// @Accept json
// @Produce json
// @Param createInviteRequest body CreateInviteRequest true "Create Invite Request"
// @Success 200 {object} models.CreateInviteResponse
// @Router /invite [post]
func (h *InviteHandler) CreateInvite(c echo.Context) error {
	var req CreateInviteRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}
	authUser := middlewares.GetAuthUser(c)

	invite, err := h.inviteService.CreateInviteToken(c.Request().Context(), authUser.FamilyId, req.GroupID, authUser.ID)
	if err != nil {
		return err
	}
	inviteURL := h.frontURL + "/login?invite_token=" + invite.Token
	return c.JSON(http.StatusOK, models.CreateInviteResponse{
		Token:     invite.Token,
		InviteURL: inviteURL,
		ExpiresAt: invite.ExpiresAt,
	})
}

// @Description 초대 토큰 검증
//
// @Summary Validate Invite Token
// @Tags Invite
// @Accept json
// @Produce json
// @Param token query string true "Invite Token"
// @Success 200 {object} models.ValidateInviteResponse
// @Router /invite/validate [get]
func (h *InviteHandler) ValidateInvite(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "field.token"})
	}
	familyName, groupName, familyID, groupID, err := h.inviteService.ValidateInviteToken(c.Request().Context(), token)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.ValidateInviteResponse{
		Valid:       true,
		FamilyName:  familyName,
		GroupName:   groupName,
		FamilyID:    familyID,
		GroupID:     groupID,
	})
}
