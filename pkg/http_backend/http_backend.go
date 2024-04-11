package httpbackend

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ophum/tfstate-manager/pkg/models"
	"gorm.io/gorm"
)

type HTTPBackendServer struct {
	db *gorm.DB
}

func NewHTTPBackendServer(db *gorm.DB) *HTTPBackendServer {
	return &HTTPBackendServer{
		db: db,
	}
}

func (s *HTTPBackendServer) RegisterHandlers(router gin.IRouter) {
	r := router.Group("/states/:id")
	r.GET("", s.get)
	r.POST("", s.create)
}

type Params struct {
	ID uint64 `uri:"id"`
}

func (s *HTTPBackendServer) get(ctx *gin.Context) {
	var params Params
	if err := ctx.ShouldBindUri(&params); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var state models.State
	if err := s.db.Where("id = ?", params.ID).First(&state).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Data(http.StatusOK, "application/json", []byte(state.State))
}

func (s *HTTPBackendServer) create(ctx *gin.Context) {
	var params Params
	if err := ctx.ShouldBindUri(&params); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var state models.State
	if err := s.db.Where("id = ?", params.ID).First(&state).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := s.db.
		Model(&models.State{}).
		Where("id = ?", params.ID).
		Update("state", data).Error; err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
