DROP INDEX IF EXISTS idx_report_card_details_card_subject;
DROP TABLE IF EXISTS report_card_details;

DROP INDEX IF EXISTS idx_report_cards_unique;
DROP TABLE IF EXISTS report_cards;

DROP INDEX IF EXISTS idx_grades_assessment_student;
DROP TABLE IF EXISTS grades;

DROP INDEX IF EXISTS idx_assessments_tenant_class;
ALTER TABLE assessments DROP CONSTRAINT IF EXISTS assessments_category_id_fkey;
DROP TABLE IF EXISTS assessments;

DROP INDEX IF EXISTS idx_grade_categories_tenant_name;
DROP TABLE IF EXISTS grade_categories;

