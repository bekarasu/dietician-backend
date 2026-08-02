ALTER TABLE medical_uploads
    ADD COLUMN parsed_results JSONB,
    ADD COLUMN parsed_at TIMESTAMP WITH TIME ZONE;
