ALTER TABLE report_cards
ADD COLUMN class_id UUID,
ADD COLUMN gpa NUMERIC(4,2),
ADD COLUMN total_credits INTEGER,
ADD COLUMN attendance TEXT,
ADD COLUMN comments TEXT,
ADD COLUMN generated_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE report_cards
RENAME COLUMN publish_status TO status;

ALTER TABLE report_cards DROP CONSTRAINT report_cards_publish_status_check;
ALTER TABLE report_cards ADD CONSTRAINT report_cards_status_check CHECK (status IN ('draft', 'generated', 'published'));

ALTER TABLE report_card_details
ADD COLUMN subject_name TEXT,
ADD COLUMN credit INTEGER;
