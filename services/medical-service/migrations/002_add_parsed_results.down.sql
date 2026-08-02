ALTER TABLE medical_uploads
    DROP COLUMN IF EXISTS parsed_results,
    DROP COLUMN IF EXISTS parsed_at;
