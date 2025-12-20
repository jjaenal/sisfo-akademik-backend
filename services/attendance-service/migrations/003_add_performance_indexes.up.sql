CREATE INDEX IF NOT EXISTS idx_student_attendance_class_date ON student_attendance (class_id, attendance_date);
CREATE INDEX IF NOT EXISTS idx_student_attendance_student_semester ON student_attendance (student_id, semester_id);
CREATE INDEX IF NOT EXISTS idx_teacher_attendance_teacher_semester_date ON teacher_attendance (teacher_id, semester_id, attendance_date);
