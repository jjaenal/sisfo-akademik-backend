package postgres_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/repository/postgres"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestScheduleRepository_CheckConflicts(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := postgres.NewScheduleRepository(mock)
	
	tenantID := "tenant-1"
	classID := uuid.New()
	teacherID := uuid.New()
	subjectID := uuid.New()
	room := "Room 101"
	
	t.Run("Conflict Detected", func(t *testing.T) {
		s := &entity.Schedule{
			ID:        uuid.New(),
			TenantID:  tenantID,
			ClassID:   classID,
			SubjectID: subjectID,
			TeacherID: teacherID,
			DayOfWeek: 1,
			StartTime: "08:00",
			EndTime:   "10:00",
			Room:      room,
		}

		rows := pgxmock.NewRows([]string{
			"id", "tenant_id", "class_id", "subject_id", "teacher_id", "day_of_week",
			"start_time", "end_time", "room", "created_at", "updated_at", "created_by", "updated_by",
		}).AddRow(
			uuid.New(), tenantID, classID, subjectID, teacherID, 1,
			"09:00", "11:00", room, time.Now(), time.Now(), nil, nil,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT 
			id, tenant_id, class_id, subject_id, teacher_id, day_of_week,
			start_time, end_time, room, created_at, updated_at, created_by, updated_by
		FROM schedules`)).
			WithArgs(
				s.TenantID,
				s.DayOfWeek,
				s.ClassID,
				s.TeacherID,
				s.Room,
				s.EndTime,
				s.StartTime,
				s.ID,
			).
			WillReturnRows(rows)

		conflicts, err := repo.CheckConflicts(context.Background(), s)
		assert.NoError(t, err)
		assert.Len(t, conflicts, 1)
	})

	t.Run("No Conflict", func(t *testing.T) {
		s := &entity.Schedule{
			ID:        uuid.New(),
			TenantID:  tenantID,
			ClassID:   classID,
			SubjectID: subjectID,
			TeacherID: teacherID,
			DayOfWeek: 1,
			StartTime: "08:00",
			EndTime:   "10:00",
			Room:      room,
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT 
			id, tenant_id, class_id, subject_id, teacher_id, day_of_week,
			start_time, end_time, room, created_at, updated_at, created_by, updated_by
		FROM schedules`)).
			WithArgs(
				s.TenantID,
				s.DayOfWeek,
				s.ClassID,
				s.TeacherID,
				s.Room,
				s.EndTime,
				s.StartTime,
				s.ID,
			).
			WillReturnRows(pgxmock.NewRows([]string{}))

		conflicts, err := repo.CheckConflicts(context.Background(), s)
		assert.NoError(t, err)
		assert.Empty(t, conflicts)
	})
	
	t.Run("DB Error", func(t *testing.T) {
		s := &entity.Schedule{
			ID:        uuid.New(),
			TenantID:  tenantID,
			ClassID:   classID,
			SubjectID: subjectID,
			TeacherID: teacherID,
			DayOfWeek: 1,
			StartTime: "08:00",
			EndTime:   "10:00",
			Room:      room,
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT 
			id, tenant_id, class_id, subject_id, teacher_id, day_of_week,
			start_time, end_time, room, created_at, updated_at, created_by, updated_by
		FROM schedules`)).
			WillReturnError(pgx.ErrTxClosed)

		conflicts, err := repo.CheckConflicts(context.Background(), s)
		assert.Error(t, err)
		assert.Nil(t, conflicts)
	})
}

func TestScheduleRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := postgres.NewScheduleRepository(mock)

	t.Run("Success", func(t *testing.T) {
		s := &entity.Schedule{
			TenantID:  "tenant-1",
			ClassID:   uuid.New(),
			SubjectID: uuid.New(),
			TeacherID: uuid.New(),
			DayOfWeek: 1,
			StartTime: "08:00",
			EndTime:   "10:00",
			Room:      "Room 101",
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO schedules`)).
			WithArgs(
				pgxmock.AnyArg(), // ID
				s.TenantID,
				s.ClassID,
				s.SubjectID,
				s.TeacherID,
				s.DayOfWeek,
				s.StartTime,
				s.EndTime,
				s.Room,
				pgxmock.AnyArg(), // CreatedAt
				pgxmock.AnyArg(), // UpdatedAt
				(*uuid.UUID)(nil), // CreatedBy
				(*uuid.UUID)(nil), // UpdatedBy
			).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := repo.Create(context.Background(), s)
		assert.NoError(t, err)
	})
}

func TestScheduleRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := postgres.NewScheduleRepository(mock)

	t.Run("Found", func(t *testing.T) {
		id := uuid.New()
		rows := pgxmock.NewRows([]string{
			"id", "tenant_id", "class_id", "subject_id", "teacher_id", "day_of_week",
			"start_time", "end_time", "room", "created_at", "updated_at", "created_by", "updated_by", "deleted_at",
		}).AddRow(
			id, "tenant-1", uuid.New(), uuid.New(), uuid.New(), 1,
			"08:00", "10:00", "Room 101", time.Now(), time.Now(), (*uuid.UUID)(nil), (*uuid.UUID)(nil), (*time.Time)(nil),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT 
			id, tenant_id, class_id, subject_id, teacher_id, day_of_week,
			start_time, end_time, room, created_at, updated_at, created_by, updated_by, deleted_at
		FROM schedules`)).
			WithArgs(id).
			WillReturnRows(rows)

		s, err := repo.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.NotNil(t, s)
		assert.Equal(t, id, s.ID)
	})

	t.Run("Not Found", func(t *testing.T) {
		id := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WithArgs(id).
			WillReturnError(pgx.ErrNoRows)

		s, err := repo.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Nil(t, s)
	})
}

func TestScheduleRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := postgres.NewScheduleRepository(mock)

	t.Run("Success", func(t *testing.T) {
		s := &entity.Schedule{
			ID:        uuid.New(),
			ClassID:   uuid.New(),
			SubjectID: uuid.New(),
			TeacherID: uuid.New(),
			DayOfWeek: 2,
			StartTime: "10:00",
			EndTime:   "12:00",
			Room:      "Room 102",
		}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE schedules SET`)).
			WithArgs(
				s.ClassID,
				s.SubjectID,
				s.TeacherID,
				s.DayOfWeek,
				s.StartTime,
				s.EndTime,
				s.Room,
				pgxmock.AnyArg(), // UpdatedAt
				(*uuid.UUID)(nil), // UpdatedBy
				s.ID,
			).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.Update(context.Background(), s)
		assert.NoError(t, err)
	})
}

func TestScheduleRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := postgres.NewScheduleRepository(mock)

	t.Run("Success", func(t *testing.T) {
		id := uuid.New()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE schedules SET deleted_at = $1 WHERE id = $2`)).
			WithArgs(pgxmock.AnyArg(), id).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.Delete(context.Background(), id)
		assert.NoError(t, err)
	})
}

func TestScheduleRepository_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := postgres.NewScheduleRepository(mock)

	t.Run("Success", func(t *testing.T) {
		tenantID := "tenant-1"
		limit := 10
		offset := 0

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM schedules`)).
			WithArgs(tenantID).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))

		rows := pgxmock.NewRows([]string{
			"id", "tenant_id", "class_id", "subject_id", "teacher_id", "day_of_week",
			"start_time", "end_time", "room", "created_at", "updated_at", "created_by", "updated_by", "deleted_at",
		}).AddRow(
			uuid.New(), tenantID, uuid.New(), uuid.New(), uuid.New(), 1,
			"08:00", "10:00", "Room 101", time.Now(), time.Now(), (*uuid.UUID)(nil), (*uuid.UUID)(nil), (*time.Time)(nil),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT 
			id, tenant_id, class_id, subject_id, teacher_id, day_of_week,
			start_time, end_time, room, created_at, updated_at, created_by, updated_by, deleted_at
		FROM schedules`)).
			WithArgs(tenantID, limit, offset).
			WillReturnRows(rows)

		schedules, total, err := repo.List(context.Background(), tenantID, limit, offset)
		assert.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Len(t, schedules, 1)
	})
}

func TestScheduleRepository_ListByClass(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := postgres.NewScheduleRepository(mock)

	t.Run("Success", func(t *testing.T) {
		classID := uuid.New()

		rows := pgxmock.NewRows([]string{
			"id", "tenant_id", "class_id", "subject_id", "teacher_id", "day_of_week",
			"start_time", "end_time", "room", "created_at", "updated_at", "created_by", "updated_by", "deleted_at",
		}).AddRow(
			uuid.New(), "tenant-1", classID, uuid.New(), uuid.New(), 1,
			"08:00", "10:00", "Room 101", time.Now(), time.Now(), (*uuid.UUID)(nil), (*uuid.UUID)(nil), (*time.Time)(nil),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT 
			id, tenant_id, class_id, subject_id, teacher_id, day_of_week,
			start_time, end_time, room, created_at, updated_at, created_by, updated_by, deleted_at
		FROM schedules`)).
			WithArgs(classID).
			WillReturnRows(rows)

		schedules, err := repo.ListByClass(context.Background(), classID)
		assert.NoError(t, err)
		assert.Len(t, schedules, 1)
	})
}

func TestScheduleRepository_BulkCreate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := postgres.NewScheduleRepository(mock)

	t.Run("Success", func(t *testing.T) {
		schedules := []*entity.Schedule{
			{
				TenantID:  "tenant-1",
				ClassID:   uuid.New(),
				SubjectID: uuid.New(),
				TeacherID: uuid.New(),
				DayOfWeek: 1,
				StartTime: "08:00",
				EndTime:   "10:00",
				Room:      "Room 101",
			},
			{
				TenantID:  "tenant-1",
				ClassID:   uuid.New(),
				SubjectID: uuid.New(),
				TeacherID: uuid.New(),
				DayOfWeek: 2,
				StartTime: "10:00",
				EndTime:   "12:00",
				Room:      "Room 102",
			},
		}

		mock.ExpectBegin()
		
		query := regexp.QuoteMeta(`INSERT INTO schedules`)
		
		// Expect 2 inserts
		mock.ExpectExec(query).WithArgs(
			pgxmock.AnyArg(), // ID
			schedules[0].TenantID,
			schedules[0].ClassID,
			schedules[0].SubjectID,
			schedules[0].TeacherID,
			schedules[0].DayOfWeek,
			schedules[0].StartTime,
			schedules[0].EndTime,
			schedules[0].Room,
			pgxmock.AnyArg(), // CreatedAt
			pgxmock.AnyArg(), // UpdatedAt
			(*uuid.UUID)(nil), // CreatedBy
			(*uuid.UUID)(nil), // UpdatedBy
		).WillReturnResult(pgxmock.NewResult("INSERT", 1))

		mock.ExpectExec(query).WithArgs(
			pgxmock.AnyArg(), // ID
			schedules[1].TenantID,
			schedules[1].ClassID,
			schedules[1].SubjectID,
			schedules[1].TeacherID,
			schedules[1].DayOfWeek,
			schedules[1].StartTime,
			schedules[1].EndTime,
			schedules[1].Room,
			pgxmock.AnyArg(), // CreatedAt
			pgxmock.AnyArg(), // UpdatedAt
			(*uuid.UUID)(nil), // CreatedBy
			(*uuid.UUID)(nil), // UpdatedBy
		).WillReturnResult(pgxmock.NewResult("INSERT", 1))

		mock.ExpectCommit()

		err := repo.BulkCreate(context.Background(), schedules)
		assert.NoError(t, err)
	})

	t.Run("Transaction Rollback on Error", func(t *testing.T) {
		schedules := []*entity.Schedule{
			{
				TenantID:  "tenant-1",
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO schedules`)).WillReturnError(pgx.ErrTxClosed)
		mock.ExpectRollback()

		err := repo.BulkCreate(context.Background(), schedules)
		assert.Error(t, err)
	})
}
