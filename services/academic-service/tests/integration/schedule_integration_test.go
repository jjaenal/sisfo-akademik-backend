package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestScheduleIntegration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Repository
	mockRepo := mocks.NewMockScheduleRepository(ctrl)
	mockTemplateRepo := mocks.NewMockScheduleTemplateRepository(ctrl)

	// Real UseCase with Mock Repo
	u := usecase.NewScheduleUseCase(mockRepo, mockTemplateRepo, time.Second*2)

	// Real Handler with Real UseCase
	h := handler.NewScheduleHandler(u)

	// Setup Router
	r := gin.New()
	schedules := r.Group("/api/v1/schedules")
	{
		schedules.POST("", h.Create)
		schedules.POST("/bulk", h.BulkCreate)
		schedules.POST("/from-template", h.CreateFromTemplate)
		schedules.GET("/:id", h.GetByID)
		schedules.GET("", h.List)
		t.Run("Create From Template Success", func(t *testing.T) {
		templateID := uuid.New()
		classID := uuid.New()
		subjectID := uuid.New()
		teacherID := uuid.New()

		reqBody := map[string]interface{}{
			"template_id": templateID.String(),
			"class_id":    classID.String(),
			"assignments": []map[string]interface{}{
				{
					"subject_id": subjectID.String(),
					"teacher_id": teacherID.String(),
				},
			},
		}
		body, _ := json.Marshal(reqBody)

		// Mock Template Repo
		template := &entity.ScheduleTemplate{
			ID:       templateID,
			TenantID: "tenant-1",
		}
		items := []entity.ScheduleTemplateItem{
			{
				ID:         uuid.New(),
				TemplateID: templateID,
				SubjectID:  &subjectID,
				DayOfWeek:  1,
				StartTime:  "08:00",
				EndTime:    "09:30",
			},
		}

		mockTemplateRepo.EXPECT().GetByID(gomock.Any(), templateID).Return(template, nil)
		mockTemplateRepo.EXPECT().ListItems(gomock.Any(), templateID).Return(items, nil)

		// Mock Schedule Repo
		mockRepo.EXPECT().CheckConflicts(gomock.Any(), gomock.Any()).Return([]entity.Schedule{}, nil)
		mockRepo.EXPECT().BulkCreate(gomock.Any(), gomock.Any()).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/schedules/from-template", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp struct {
			Success bool `json:"success"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	})
}

	t.Run("Bulk Create Schedule Success", func(t *testing.T) {
		reqBody := []map[string]interface{}{
			{
				"tenant_id":   "tenant-1",
				"class_id":    uuid.New().String(),
				"subject_id":  uuid.New().String(),
				"teacher_id":  uuid.New().String(),
				"day_of_week": 1,
				"start_time":  "08:00",
				"end_time":    "09:30",
				"room":        "Room 101",
			},
			{
				"tenant_id":   "tenant-1",
				"class_id":    uuid.New().String(),
				"subject_id":  uuid.New().String(),
				"teacher_id":  uuid.New().String(),
				"day_of_week": 2,
				"start_time":  "10:00",
				"end_time":    "11:30",
				"room":        "Room 102",
			},
		}
		body, _ := json.Marshal(reqBody)

		// Expect CheckConflicts for each schedule
		mockRepo.EXPECT().CheckConflicts(gomock.Any(), gomock.Any()).Return([]entity.Schedule{}, nil).Times(2)

		// Expect BulkCreate
		mockRepo.EXPECT().BulkCreate(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, schedules []*entity.Schedule) error {
			assert.Equal(t, 2, len(schedules))
			assert.Equal(t, "Room 101", schedules[0].Room)
			assert.Equal(t, "Room 102", schedules[1].Room)
			return nil
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/schedules/bulk", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool                   `json:"success"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, float64(2), resp.Data["count"])
	})

	t.Run("Create Schedule Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":   "tenant-1",
			"class_id":    uuid.New().String(),
			"subject_id":  uuid.New().String(),
			"teacher_id":  uuid.New().String(),
			"day_of_week": 1,
			"start_time":  "08:00",
			"end_time":    "09:30",
			"room":        "Room 101",
		}
		body, _ := json.Marshal(reqBody)

		// Expect CheckConflicts (return empty -> no conflict)
		mockRepo.EXPECT().CheckConflicts(gomock.Any(), gomock.Any()).Return([]entity.Schedule{}, nil)

		// Expect Create
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, s *entity.Schedule) error {
			assert.Equal(t, "Room 101", s.Room)
			assert.Equal(t, "tenant-1", s.TenantID)
			s.ID = uuid.New()
			return nil
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/schedules", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool            `json:"success"`
			Data    entity.Schedule `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "Room 101", resp.Data.Room)
		assert.NotEmpty(t, resp.Data.ID)
	})

	t.Run("Create Schedule Conflict", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":   "tenant-1",
			"class_id":    uuid.New().String(),
			"subject_id":  uuid.New().String(),
			"teacher_id":  uuid.New().String(),
			"day_of_week": 1,
			"start_time":  "08:00",
			"end_time":    "09:30",
			"room":        "Room 101",
		}
		body, _ := json.Marshal(reqBody)

		// Expect CheckConflicts (return existing schedule -> conflict)
		existingSchedule := entity.Schedule{ID: uuid.New(), Room: "Room 101"}
		mockRepo.EXPECT().CheckConflicts(gomock.Any(), gomock.Any()).Return([]entity.Schedule{existingSchedule}, nil)

		// Create should NOT be called

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/schedules", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		// Handler currently returns 500 on usecase error
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		var resp struct {
			Success bool `json:"success"`
			Error   struct {
				Message string      `json:"message"`
				Details interface{} `json:"details"`
			} `json:"error"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.False(t, resp.Success)
		assert.Contains(t, resp.Error.Details, "conflict detected")
	})

	t.Run("Get Schedule Flow", func(t *testing.T) {
		id := uuid.New()
		expectedSchedule := &entity.Schedule{
			ID:        id,
			TenantID:  "tenant-1",
			Room:      "Room 101",
			StartTime: "08:00",
			EndTime:   "09:30",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(expectedSchedule, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/schedules/"+id.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool            `json:"success"`
			Data    entity.Schedule `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, id, resp.Data.ID)
		assert.Equal(t, "Room 101", resp.Data.Room)
	})
}
