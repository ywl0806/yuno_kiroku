package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/pkg/utils"
)

type TagHandler struct {
	tagService services.TagService
}

func NewTagHandler(tagService services.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

type CreateTagRequest struct {
	Name string `json:"name" validate:"required"`
}

type AddTagToMediaItemRequest struct {
	TagID int32 `json:"tag_id" validate:"required"`
}

// @Tags Tag
// @Description 가족 태그 목록 조회 (프리셋 + 커스텀)
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /tag [get]
func (h *TagHandler) GetTags(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()

	tags, err := h.tagService.GetTagsByFamilyID(ctx, authUser.FamilyId)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewTagResponses(tags))
}

// @Tags Tag
// @Description 커스텀 태그 생성
// @Param body body CreateTagRequest true "태그 정보"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /tag [post]
func (h *TagHandler) CreateTag(c echo.Context) error {
	req := new(CreateTagRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()

	tag, err := h.tagService.CreateTag(ctx, authUser.FamilyId, authUser.ID, req.Name)
	if err != nil {
		return err
	}
	return c.JSON(201, models.NewTagResponse(tag))
}

// @Tags Tag
// @Description 커스텀 태그 삭제 (프리셋 태그는 삭제 불가)
// @Param id path int true "Tag ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /tag/{id} [delete]
func (h *TagHandler) DeleteTag(c echo.Context) error {
	tagID, err := utils.ConvertToInt32(c.Param("id"))
	if err != nil {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "id"})
	}
	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()

	if err := h.tagService.DeleteTag(ctx, tagID, authUser.FamilyId); err != nil {
		return err
	}
	return c.NoContent(204)
}

// @Tags Tag
// @Description 미디어 아이템에 태그 추가
// @Param id path int true "MediaItem ID"
// @Param body body AddTagToMediaItemRequest true "태그 ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/{id}/tag [post]
func (h *TagHandler) AddTagToMediaItem(c echo.Context) error {
	mediaItemID, err := utils.ConvertToInt32(c.Param("id"))
	if err != nil {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "id"})
	}
	req := new(AddTagToMediaItemRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()

	if err := h.tagService.AddTagToMediaItem(ctx, mediaItemID, req.TagID, authUser.ID); err != nil {
		return err
	}
	return c.JSON(200, map[string]bool{"ok": true})
}

// @Tags Tag
// @Description 미디어 아이템에서 태그 제거
// @Param id path int true "MediaItem ID"
// @Param tagId path int true "Tag ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/{id}/tag/{tagId} [delete]
func (h *TagHandler) RemoveTagFromMediaItem(c echo.Context) error {
	mediaItemID, err := utils.ConvertToInt32(c.Param("id"))
	if err != nil {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "id"})
	}
	tagID, err := utils.ConvertToInt32(c.Param("tagId"))
	if err != nil {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "tagId"})
	}
	ctx := c.Request().Context()

	if err := h.tagService.RemoveTagFromMediaItem(ctx, mediaItemID, tagID); err != nil {
		return err
	}
	return c.NoContent(204)
}

// @Tags Tag
// @Description 미디어 아이템의 태그 목록 조회
// @Param id path int true "MediaItem ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/{id}/tag [get]
func (h *TagHandler) GetTagsForMediaItem(c echo.Context) error {
	mediaItemID, err := utils.ConvertToInt32(c.Param("id"))
	if err != nil {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "id"})
	}
	ctx := c.Request().Context()

	tags, err := h.tagService.GetTagsForMediaItem(ctx, mediaItemID)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewTagResponses(tags))
}
