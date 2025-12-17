-- Admission Service initial schema
-- Tables: admission_periods, applications, application_documents

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS admission_periods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id TEXT NOT NULL DEFAULT 'default',
    name TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_announced BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_admission_periods_tenant_name
    ON admission_periods (tenant_id, name)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id TEXT NOT NULL DEFAULT 'default', -- Added default for now as code might not set it
    admission_period_id UUID NOT NULL REFERENCES admission_periods(id),
    registration_number TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT,
    phone_number TEXT,
    status TEXT NOT NULL DEFAULT 'submitted',
    previous_school TEXT,
    average_score NUMERIC(5,2),
    submission_date TIMESTAMP WITH TIME ZONE,
    test_score NUMERIC(5,2),
    interview_score NUMERIC(5,2),
    final_score NUMERIC(5,2),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_applications_number_unique
    ON applications (tenant_id, registration_number)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_applications_period_status
    ON applications (tenant_id, admission_period_id, status)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS application_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    application_id UUID NOT NULL REFERENCES applications(id),
    document_type TEXT NOT NULL,
    file_url TEXT NOT NULL,
    file_name TEXT,
    file_size BIGINT,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_application_documents_app_type
    ON application_documents (application_id, document_type)
    WHERE deleted_at IS NULL;
