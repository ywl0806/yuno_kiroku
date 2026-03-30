package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/services"
)

type KidHandler struct {
	kidService *services.KidService
}

func NewKidHandler(kidService *services.KidService) *KidHandler {
	return &KidHandler{kidService: kidService}
}

// @Tags Kid
// @Description Create a new kid
// @Router /kid [post]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param body body models.CreateKidRequest true "Create kid request"
// @Success 201 {object} models.KidResponse
func (h *KidHandler) CreateKid(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	req := new(models.CreateKidRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	var birthDate *time.Time
	if req.BirthDate != nil && *req.BirthDate != "" {
		t, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid birth_date format, use YYYY-MM-DD")
		}
		birthDate = &t
	}

	kid, err := h.kidService.CreateKid(c.Request().Context(), authUser.FamilyId, req.Name, birthDate, req.IdentityID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, models.NewKidResponse(&kid))
}

// @Tags Kid
// @Description Update a kid
// @Router /kid/{kidId} [put]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param kidId path int true "Kid ID"
// @Param body body models.UpdateKidRequest true "Update kid request"
// @Success 200 {object} models.KidResponse
func (h *KidHandler) UpdateKid(c echo.Context) error {
	kidId, err := strconv.Atoi(c.Param("kidId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid kid id")
	}

	req := new(models.UpdateKidRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	var birthDate *time.Time
	if req.BirthDate != nil && *req.BirthDate != "" {
		t, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid birth_date format, use YYYY-MM-DD")
		}
		birthDate = &t
	}

	kid, err := h.kidService.UpdateKid(c.Request().Context(), int32(kidId), req.Name, birthDate, req.IdentityID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.NewKidResponse(&kid))
}

// @Tags Kid
// @Description 월별 아이 얼굴 사진 조회
// @Router /kid/face-imgs [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param body body models.GetKidsWithFaceImgRequest true "Get kids with face img request"
// @Success 200 {array} models.KidWithFaceImgResponse
func (h *KidHandler) GetKidsWithFaceImg(c echo.Context) error {

	req := new(models.GetKidsWithFaceImgRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	authUser := middlewares.GetAuthUser(c)
	kids, err := h.kidService.GetKidsWithFaceImg(c.Request().Context(), authUser.FamilyId, req.Year, req.Month)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.NewKidWithFaceImgResponses(kids))
}

// @Tags Kid
// @Description Delete a kid
// @Router /kid/{kidId} [delete]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param kidId path int true "Kid ID"
// @Success 204
func (h *KidHandler) DeleteKid(c echo.Context) error {
	kidId, err := strconv.Atoi(c.Param("kidId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid kid id")
	}

	if err := h.kidService.DeleteKid(c.Request().Context(), int32(kidId)); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
