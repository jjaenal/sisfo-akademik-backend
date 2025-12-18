ALTER TABLE assessments RENAME COLUMN grade_category_id TO category_id;
ALTER TABLE assessments DROP COLUMN IF NOT EXISTS teacher_id;
