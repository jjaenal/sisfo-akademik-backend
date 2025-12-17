-- Attendance Service initial schema
-- Tables: student_attendance, teacher_attendance

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS student_attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id TEXT NOT NULL,
    student_id UUID NOT NULL,
    class_id UUID,
    attendance_date DATE NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('present','absent','late','excused','sick')),
    check_in_time TIMESTAMP WITH TIME ZONE,
    check_out_time TIMESTAMP WITH TIME ZONE,
    check_in_latitude DOUBLE PRECISION,
    check_in_longitude DOUBLE PRECISION,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_student_attendance_student_date
    ON student_attendance (student_id, attendance_date)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_student_attendance_tenant_date
    ON student_attendance (tenant_id, attendance_date)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS teacher_attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id TEXT NOT NULL,
    teacher_id UUID NOT NULL,
    attendance_date DATE NOT NULL,
    check_in_time TIMESTAMP WITH TIME ZONE,
    check_out_time TIMESTAMP WITH TIME ZONE,
    status TEXT NOT NULL CHECK (status IN ('present','absent','late','excused','sick')),
    location_latitude DOUBLE PRECISION,
    location_longitude DOUBLE PRECISION,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_teacher_attendance_teacher_date
    ON teacher_attendance (teacher_id, attendance_date)
    WHERE deleted_at IS NULL;

