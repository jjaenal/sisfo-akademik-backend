ALTER TABLE applications ADD COLUMN IF NOT EXISTS previous_school TEXT;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS average_score NUMERIC(5,2);
