package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/file-service/internal/domain"
	httputil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/middleware"
)

type FileHandler struct {
	uc domain.FileUseCase
}

func NewFileHandler(uc domain.FileUseCase) *FileHandler {
	return &FileHandler{uc: uc}
}

func (h *FileHandler) RegisterRoutes(r gin.IRouter) {
	api := r.Group("/api/v1/files")
	// Protected routes should be handled by middleware in main
	api.POST("/upload", h.upload)
	api.GET("/:id", h.get)
	api.GET("/:id/download", h.download)
	api.DELETE("/:id", h.delete)
	api.GET("", h.list)
}

func (h *FileHandler) download(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4003", "Bad Request", "invalid id format")
		return
	}

	fileMeta, fileContent, err := h.uc.Download(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if fileMeta == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "file not found")
		return
	}
	defer fileContent.Close()

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileMeta.OriginalName))
	c.Header("Content-Type", fileMeta.MimeType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileMeta.Size))

	io.Copy(c.Writer, fileContent)
}

func (h *FileHandler) upload(c *gin.Context) {
	val := c.Request.Context().Value(middleware.ClaimsKey)
	if val == nil {
		httputil.Error(c.Writer, http.StatusUnauthorized, "4001", "Unauthorized", "missing claims")
		return
	}
	claims, ok := val.(jwtutil.Claims)
	if !ok {
		httputil.Error(c.Writer, http.StatusUnauthorized, "4001", "Unauthorized", "invalid claims")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4002", "Bad Request", "file is required")
		return
	}
	defer file.Close()

	bucket := c.DefaultPostForm("bucket", "uploads")

	tenantID, err := uuid.Parse(claims.TenantID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusUnauthorized, "4001", "Unauthorized", "invalid tenant id")
		return
	}

	res, err := h.uc.Upload(c.Request.Context(), file, header, tenantID, claims.UserID, bucket)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, res)
}

func (h *FileHandler) get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4003", "Bad Request", "invalid id format")
		return
	}

	res, err := h.uc.Get(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if res == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "file not found")
		return
	}

	httputil.Success(c.Writer, res)
}

func (h *FileHandler) delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4003", "Bad Request", "invalid id format")
		return
	}

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, gin.H{"status": "deleted"})
}

func (h *FileHandler) list(c *gin.Context) {
	val := c.Request.Context().Value(middleware.ClaimsKey)
	if val == nil {
		httputil.Error(c.Writer, http.StatusUnauthorized, "4001", "Unauthorized", "missing claims")
		return
	}
	claims, ok := val.(jwtutil.Claims)
	if !ok {
		httputil.Error(c.Writer, http.StatusUnauthorized, "4001", "Unauthorized", "invalid claims")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	tenantID, err := uuid.Parse(claims.TenantID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusUnauthorized, "4001", "Unauthorized", "invalid tenant id")
		return
	}

	res, _, err := h.uc.List(c.Request.Context(), tenantID, page, limit)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, res)
}
