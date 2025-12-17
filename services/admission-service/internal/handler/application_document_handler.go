package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type ApplicationDocumentHandler struct {
	useCase usecase.ApplicationDocumentUseCase
}

func NewApplicationDocumentHandler(useCase usecase.ApplicationDocumentUseCase) *ApplicationDocumentHandler {
	return &ApplicationDocumentHandler{useCase: useCase}
}

func (h *ApplicationDocumentHandler) Upload(c *gin.Context) {
	idStr := c.Param("id")
	applicationID, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid application ID", err.Error())
		return
	}

	documentType := c.PostForm("document_type")
	if documentType == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4002", "document_type is required", nil)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4003", "file is required", err.Error())
		return
	}

	doc, err := h.useCase.Upload(c.Request.Context(), applicationID, documentType, file)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to upload document", err.Error())
		return
	}

	httputil.Success(c.Writer, doc)
}

func (h *ApplicationDocumentHandler) GetByApplicationID(c *gin.Context) {
	idStr := c.Param("id")
	applicationID, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid application ID", err.Error())
		return
	}

	docs, err := h.useCase.GetByApplicationID(c.Request.Context(), applicationID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to fetch documents", err.Error())
		return
	}

	httputil.Success(c.Writer, docs)
}

func (h *ApplicationDocumentHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid document ID", err.Error())
		return
	}

	if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to delete document", err.Error())
		return
	}

	httputil.Success(c.Writer, nil)
}
