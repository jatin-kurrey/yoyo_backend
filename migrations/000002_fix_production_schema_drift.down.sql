-- 000002_fix_production_schema_drift.down.sql

-- Remove added columns
ALTER TABLE site_settings 
DROP COLUMN IF EXISTS whatsapp_number;

ALTER TABLE tickets 
DROP COLUMN IF EXISTS is_bestseller;

-- Restore HomepageSections constraints (optional/conservative)
-- Note: We keep the data as is, but we could re-enforce NOT NULL if all NULLs were cleared
-- For a safe 'down', we just remove the columns we added.
