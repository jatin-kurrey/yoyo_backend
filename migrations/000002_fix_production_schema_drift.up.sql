-- 000002_fix_production_schema_drift.up.sql

-- Add WhatsApp Number to site_settings if missing
ALTER TABLE site_settings 
ADD COLUMN IF NOT EXISTS whatsapp_number VARCHAR(20) DEFAULT '';

-- Add IsBestseller to tickets if missing
ALTER TABLE tickets 
ADD COLUMN IF NOT EXISTS is_bestseller BOOLEAN NOT NULL DEFAULT false;

-- Fix HomepageSections constraints and defaults
-- 1. Drop NOT NULL if it exists (to allow safe migration)
ALTER TABLE site_settings 
ALTER COLUMN homepage_sections DROP NOT NULL;

-- 2. Set default for new rows
ALTER TABLE site_settings 
ALTER COLUMN homepage_sections SET DEFAULT '{}'::jsonb;

-- 3. Backfill NULL values with empty JSON object
UPDATE site_settings 
SET homepage_sections = '{}'::jsonb 
WHERE homepage_sections IS NULL;
