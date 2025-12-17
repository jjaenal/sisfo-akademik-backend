ALTER TABLE teacher_attendance ADD COLUMN IF NOT EXISTS semester_id UUID;
ALTER TABLE student_attendance ADD COLUMN IF NOT EXISTS semester_id UUID;
