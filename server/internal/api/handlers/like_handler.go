package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/pkg/utils"

	"github.com/ywl0806/yuno_kiroku/internal/apperr"
)

type LikeHandler struct {
	likeService *services.LikeService
}

func NewLikeHandler(likeService *services.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

// @Tags Like
// @Description 미디어 아이템 좋아요
// @Param id path int true "MediaItem ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/{id}/like [post]
func (h *LikeHandler) LikeMediaItem(c echo.Context) error {
	mediaItemID, err := utils.ConvertToInt32(c.Param("id"))
	if err != nil {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "id"})
	}
	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()

	if err := h.likeService.LikeMediaItem(ctx, mediaItemID, authUser.ID); err != nil {
		return err
	}
	return c.JSON(200, map[string]bool{"is_liked": true})
}

// @Tags Like
// @Description 미디어 아이템 좋아요 취소
// @Param id path int true "MediaItem ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/{id}/like [delete]
func (h *LikeHandler) UnlikeMediaItem(c echo.Context) error {
	mediaItemID, err := utils.ConvertToInt32(c.Param("id"))
	if err != nil {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "id"})
	}
	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()

	if err := h.likeService.UnlikeMediaItem(ctx, mediaItemID, authUser.ID); err != nil {
		return err
	}
	return c.JSON(200, map[string]bool{"is_liked": false})
}

// @Tags Like
// @Description 미디어 아이템 좋아요 상태 조회
// @Param id path int true "MediaItem ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/{id}/like [get]
func (h *LikeHandler) IsMediaItemLiked(c echo.Context) error {
	mediaItemID, err := utils.ConvertToInt32(c.Param("id"))
	if err != nil {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "id"})
	}
	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()

	isLiked, err := h.likeService.IsMediaItemLiked(ctx, mediaItemID, authUser.ID)
	if err != nil {
		return err
	}
	return c.JSON(200, map[string]bool{"is_liked": isLiked})
}
