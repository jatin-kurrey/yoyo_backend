-- 000003_add_admin_sidebar_toggles.down.sql
ALTER TABLE site_settings 
DROP COLUMN IF EXISTS admin_sidebar_toggles;

DROP TABLE IF EXISTS attractions;

