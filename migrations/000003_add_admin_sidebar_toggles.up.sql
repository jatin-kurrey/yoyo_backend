-- 000003_add_admin_sidebar_toggles.up.sql
ALTER TABLE site_settings 
ADD COLUMN IF NOT EXISTS admin_sidebar_toggles JSONB NOT NULL DEFAULT '{}'::jsonb;
