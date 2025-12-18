CREATE TABLE IF NOT EXISTS grade_categories (
    id UUID PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    weight DECIMAL(5,2) NOT NULL, -- Percentage (e.g., 30.00 for 30%)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS assessments (
    id UUID PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    grade_category_id UUID NOT NULL REFERENCES grade_categories(id),
    teacher_id UUID NOT NULL,
    subject_id UUID NOT NULL,
    class_id UUID NOT NULL,
    semester_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    max_score DECIMAL(5,2) NOT NULL DEFAULT 100.00,
    date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS grades (
    id UUID PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    assessment_id UUID NOT NULL REFERENCES assessments(id),
    student_id UUID NOT NULL,
    score DECIMAL(5,2) NOT NULL,
    feedback TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(assessment_id, student_id)
);

CREATE TABLE IF NOT EXISTS report_cards (
    id UUID PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    student_id UUID NOT NULL,
    semester_id UUID NOT NULL,
    class_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT', -- DRAFT, GENERATED, PUBLISHED
    gpa DECIMAL(4,2),
    rank INTEGER,
    total_attendance INTEGER,
    attendance_summary JSONB, -- { "present": 10, "sick": 2, ... }
    teacher_comments TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS report_card_details (
    id UUID PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    report_card_id UUID NOT NULL REFERENCES report_cards(id),
    subject_id UUID NOT NULL,
    subject_name VARCHAR(100) NOT NULL,
    final_score DECIMAL(5,2) NOT NULL,
    grade_letter VARCHAR(2) NOT NULL, -- A, B, C, etc.
    teacher_id UUID NOT NULL,
    comments TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);
