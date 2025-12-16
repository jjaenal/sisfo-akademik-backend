package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/handler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSchoolHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockSchoolUseCase(ctrl)
	h := handler.NewSchoolHandler(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id": "tenant-1",
			"name":      "SMA 1",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/schools", bytes.NewBuffer(body))

		h.Create(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestSchoolHandler_GetByTenantID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockSchoolUseCase(ctrl)
	h := handler.NewSchoolHandler(mockUseCase)

	tenantID := "tenant-1"

	t.Run("Success", func(t *testing.T) {
		school := &entity.School{TenantID: tenantID}
		mockUseCase.EXPECT().GetByTenantID(gomock.Any(), tenantID).Return(school, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/schools/tenant/"+tenantID, nil)
		c.Params = gin.Params{{Key: "tenant_id", Value: tenantID}}

		h.GetByTenantID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestSchoolHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockSchoolUseCase(ctrl)
	h := handler.NewSchoolHandler(mockUseCase)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "SMA 1 Updated",
		}
		body, _ := json.Marshal(reqBody)

		existing := &entity.School{ID: id, Name: "SMA 1"}
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
		mockUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/schools/"+id.String(), bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.Update(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestSchoolHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockSchoolUseCase(ctrl)
	h := handler.NewSchoolHandler(mockUseCase)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockUseCase.EXPECT().Delete(gomock.Any(), id).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodDelete, "/schools/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.Delete(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestSchoolHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockSchoolUseCase(ctrl)
	h := handler.NewSchoolHandler(mockUseCase)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		school := &entity.School{ID: id}
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(school, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/schools/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
