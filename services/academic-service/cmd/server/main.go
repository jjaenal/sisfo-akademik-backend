package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/database"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/logger"
	redisutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/redis"
	"go.uber.org/zap"

	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/repository/postgres"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log, err := logger.New(cfg.Env)
	if err != nil {
		panic(err)
	}
	dbPool, err := database.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		panic(err)
	}
	redis := redisutil.New(cfg.RedisAddr)

	// Init dependencies
	schoolRepo := postgres.NewSchoolRepository(dbPool)
	schoolUseCase := usecase.NewSchoolUseCase(schoolRepo, 5*time.Second) // 5s timeout
	schoolHandler := handler.NewSchoolHandler(schoolUseCase)

	academicYearRepo := postgres.NewAcademicYearRepository(dbPool)
	academicYearUseCase := usecase.NewAcademicYearUseCase(academicYearRepo, 5*time.Second)
	academicYearHandler := handler.NewAcademicYearHandler(academicYearUseCase)

	semesterRepo := postgres.NewSemesterRepository(dbPool)
	semesterUseCase := usecase.NewSemesterUseCase(semesterRepo, 5*time.Second)
	semesterHandler := handler.NewSemesterHandler(semesterUseCase)

	studentRepo := postgres.NewStudentRepository(dbPool)
	studentUseCase := usecase.NewStudentUseCase(studentRepo, 5*time.Second)
	studentHandler := handler.NewStudentHandler(studentUseCase)

	teacherRepo := postgres.NewTeacherRepository(dbPool)
	teacherUseCase := usecase.NewTeacherUseCase(teacherRepo, 5*time.Second)
	teacherHandler := handler.NewTeacherHandler(teacherUseCase)

	classRepo := postgres.NewClassRepository(dbPool)
	classUseCase := usecase.NewClassUseCase(classRepo, 5*time.Second)
	classHandler := handler.NewClassHandler(classUseCase)

	subjectRepo := postgres.NewSubjectRepository(dbPool)
	subjectUseCase := usecase.NewSubjectUseCase(subjectRepo, 5*time.Second)
	subjectHandler := handler.NewSubjectHandler(subjectUseCase)

	scheduleRepo := postgres.NewScheduleRepository(dbPool)
	scheduleUseCase := usecase.NewScheduleUseCase(scheduleRepo, 5*time.Second)
	scheduleHandler := handler.NewScheduleHandler(scheduleUseCase)

	curriculumRepo := postgres.NewCurriculumRepository(dbPool)
	curriculumUseCase := usecase.NewCurriculumUseCase(curriculumRepo, 5*time.Second)
	curriculumHandler := handler.NewCurriculumHandler(curriculumUseCase)

	enrollmentRepo := postgres.NewEnrollmentRepository(dbPool)
	enrollmentUseCase := usecase.NewEnrollmentUseCase(enrollmentRepo, classRepo, 5*time.Second)
	enrollmentHandler := handler.NewEnrollmentHandler(enrollmentUseCase)

	classSubjectRepo := postgres.NewClassSubjectRepository(dbPool)
	classSubjectUseCase := usecase.NewClassSubjectUseCase(classSubjectRepo, 5*time.Second)
	classSubjectHandler := handler.NewClassSubjectHandler(classSubjectUseCase)

	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Use(gin.Logger(), gin.Recovery())

	// Health Check
	r.GET("/api/v1/health", func(c *gin.Context) {
		ctx := c.Request.Context()
		if err := dbPool.Ping(ctx); err != nil {
			httputil.Error(c.Writer, http.StatusServiceUnavailable, "5003", "Database Unavailable", err.Error())
			return
		}
		if err := redis.Raw().Ping(ctx).Err(); err != nil {
			httputil.Error(c.Writer, http.StatusServiceUnavailable, "5004", "Redis Unavailable", err.Error())
			return
		}
		httputil.Success(c.Writer, map[string]string{"status": "ok", "service": "academic-service"})
	})

	// Routes
	v1 := r.Group("/api/v1")
	{
		schools := v1.Group("/schools")
		{
			schools.POST("", schoolHandler.Create)
			schools.GET("/:id", schoolHandler.GetByID)
			schools.GET("/tenant/:tenant_id", schoolHandler.GetByTenantID)
			schools.PUT("/:id", schoolHandler.Update)
			schools.DELETE("/:id", schoolHandler.Delete)
		}

		academicYears := v1.Group("/academic-years")
		{
			academicYears.POST("", academicYearHandler.Create)
			academicYears.GET("/:id", academicYearHandler.GetByID)
			academicYears.GET("/tenant/:tenant_id", academicYearHandler.List)
			academicYears.PUT("/:id", academicYearHandler.Update)
			academicYears.DELETE("/:id", academicYearHandler.Delete)
		}

		semesters := v1.Group("/semesters")
		{
			semesters.POST("", semesterHandler.Create)
			semesters.GET("/:id", semesterHandler.GetByID)
			semesters.GET("", semesterHandler.List)
			semesters.PUT("/:id", semesterHandler.Update)
			semesters.PATCH("/:id/activate", semesterHandler.Activate)
			semesters.DELETE("/:id", semesterHandler.Delete)
		}

		students := v1.Group("/students")
		{
			students.POST("", studentHandler.Create)
			students.GET("/:id", studentHandler.GetByID)
			students.GET("", studentHandler.List) // Query: ?tenant_id=...&limit=...&offset=...
			students.PUT("/:id", studentHandler.Update)
			students.DELETE("/:id", studentHandler.Delete)
		}

		teachers := v1.Group("/teachers")
		{
			teachers.POST("", teacherHandler.Create)
			teachers.GET("/:id", teacherHandler.GetByID)
			teachers.GET("", teacherHandler.List) // Query: ?tenant_id=...&limit=...&offset=...
			teachers.PUT("/:id", teacherHandler.Update)
			teachers.DELETE("/:id", teacherHandler.Delete)
		}

		classes := v1.Group("/classes")
		{
			classes.POST("", classHandler.Create)
			classes.GET("/:id", classHandler.GetByID)
			classes.GET("", classHandler.List) // Query: ?tenant_id=...&limit=...&offset=...
			classes.PUT("/:id", classHandler.Update)
			classes.DELETE("/:id", classHandler.Delete)

			// Class Subjects
			classes.POST("/:id/subjects", classSubjectHandler.AddSubjectToClass)
			classes.GET("/:id/subjects", classSubjectHandler.ListByClass)
			classes.DELETE("/:id/subjects/:subject_id", classSubjectHandler.RemoveSubject)
			classes.POST("/:id/subjects/:subject_id/teacher", classSubjectHandler.AssignTeacher)
		}

		subjects := v1.Group("/subjects")
		{
			subjects.POST("", subjectHandler.Create)
			subjects.GET("/:id", subjectHandler.GetByID)
			subjects.GET("", subjectHandler.List) // Query: ?tenant_id=...&limit=...&offset=...
			subjects.PUT("/:id", subjectHandler.Update)
			subjects.DELETE("/:id", subjectHandler.Delete)
		}

		schedules := v1.Group("/schedules")
		{
			schedules.POST("", scheduleHandler.Create)
			schedules.GET("/:id", scheduleHandler.GetByID)
			schedules.GET("", scheduleHandler.List) // Query: ?tenant_id=...&limit=...&offset=...
			schedules.GET("/class/:class_id", scheduleHandler.ListByClass)
			schedules.PUT("/:id", scheduleHandler.Update)
			schedules.DELETE("/:id", scheduleHandler.Delete)
		}

		curricula := v1.Group("/curricula")
		{
			curricula.POST("", curriculumHandler.Create)
			curricula.GET("/:id", curriculumHandler.GetByID)
			curricula.GET("", curriculumHandler.List)
			curricula.PUT("/:id", curriculumHandler.Update)
			curricula.DELETE("/:id", curriculumHandler.Delete)

			// Curriculum Subjects
			curricula.POST("/:id/subjects", curriculumHandler.AddSubject)
			curricula.GET("/:id/subjects", curriculumHandler.ListSubjects)
			curricula.DELETE("/:id/subjects/:subject_id", curriculumHandler.RemoveSubject)

			// Grading Rules
			curricula.POST("/:id/grading-rules", curriculumHandler.AddGradingRule)
			curricula.GET("/:id/grading-rules", curriculumHandler.ListGradingRules)
		}

		enrollments := v1.Group("/enrollments")
		{
			enrollments.POST("", enrollmentHandler.Enroll)
			enrollments.DELETE("/:id", enrollmentHandler.Unenroll)
			enrollments.GET("/:id", enrollmentHandler.GetByID)
			enrollments.PUT("/:id/status", enrollmentHandler.UpdateStatus)
		}

		// Additional routes
		classes.GET("/:id/students", enrollmentHandler.ListByClass)
		students.GET("/:id/classes", enrollmentHandler.ListByStudent)
	}

	addr := fmt.Sprintf(":%d", cfg.HTTPPort)
	log.Info("Starting academic-service", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		log.Fatal("server failed", zap.Error(err))
	}
}
