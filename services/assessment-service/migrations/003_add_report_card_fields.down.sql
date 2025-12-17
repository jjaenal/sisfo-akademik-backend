ALTER TABLE report_card_details
DROP COLUMN subject_name,
DROP COLUMN credit;

ALTER TABLE report_cards
RENAME COLUMN status TO publish_status;

ALTER TABLE report_cards DROP CONSTRAINT report_cards_status_check;
ALTER TABLE report_cards ADD CONSTRAINT report_cards_publish_status_check CHECK (publish_status IN ('draft', 'published'));

ALTER TABLE report_cards
DROP COLUMN class_id,
DROP COLUMN gpa,
DROP COLUMN total_credits,
DROP COLUMN attendance,
DROP COLUMN comments,
DROP COLUMN generated_at;
