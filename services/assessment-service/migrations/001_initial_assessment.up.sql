-- Assessment Service initial schema
-- Tables: grade_categories, assessments, grades, report_cards, report_card_details

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS grade_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    weight NUMERIC(5,2) NOT NULL CHECK (weight >= 0 AND weight <= 100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_grade_categories_tenant_name
    ON grade_categories (tenant_id, name)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS assessments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id TEXT NOT NULL,
    subject_id UUID,         -- external reference to academic-service
    class_id UUID,           -- external reference to academic-service
    name TEXT NOT NULL,
    description TEXT,
    category_id UUID NOT NULL REFERENCES grade_categories(id),
    max_score NUMERIC(7,2) NOT NULL DEFAULT 100,
    date DATE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_assessments_tenant_class
    ON assessments (tenant_id, class_id)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS grades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id TEXT NOT NULL,
    assessment_id UUID NOT NULL REFERENCES assessments(id),
    student_id UUID NOT NULL, -- external reference to academic-service
    score NUMERIC(7,2) NOT NULL CHECK (score >= 0),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_grades_assessment_student
    ON grades (assessment_id, student_id)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS report_cards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id TEXT NOT NULL,
    student_id UUID NOT NULL,
    class_id UUID NOT NULL,
    semester_id UUID NOT NULL, -- external reference to academic-service
    status TEXT NOT NULL DEFAULT 'draft',
    gpa NUMERIC(4,2) DEFAULT 0,
    total_credits INTEGER DEFAULT 0,
    attendance JSONB DEFAULT '{}',
    comments TEXT,
    generated_at TIMESTAMP WITH TIME ZONE,
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_report_cards_unique
    ON report_cards (tenant_id, student_id, semester_id)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS report_card_details (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_card_id UUID NOT NULL REFERENCES report_cards(id),
    subject_id UUID NOT NULL,
    subject_name TEXT,
    credit INTEGER DEFAULT 0,
    final_score NUMERIC(7,2) NOT NULL,
    grade_letter TEXT,
    comments TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_report_card_details_card_subject
    ON report_card_details (report_card_id, subject_id)
    WHERE deleted_at IS NULL;

