DROP INDEX IF EXISTS idx_application_documents_app_type;
DROP TABLE IF EXISTS application_documents;

DROP INDEX IF EXISTS idx_applications_period_status;
DROP INDEX IF EXISTS idx_applications_number_unique;
ALTER TABLE applications DROP CONSTRAINT IF EXISTS applications_period_id_fkey;
DROP TABLE IF EXISTS applications;

DROP INDEX IF EXISTS idx_admission_periods_tenant_name;
DROP TABLE IF EXISTS admission_periods;

