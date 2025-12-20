CREATE INDEX IF NOT EXISTS idx_schedules_conflict_check ON schedules (tenant_id, day_of_week, deleted_at);
CREATE INDEX IF NOT EXISTS idx_schedules_class_id ON schedules (class_id);
CREATE INDEX IF NOT EXISTS idx_schedules_teacher_id ON schedules (teacher_id);
CREATE INDEX IF NOT EXISTS idx_schedules_room ON schedules (room);
