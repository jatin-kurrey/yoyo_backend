-- 000003_add_admin_sidebar_toggles.up.sql
ALTER TABLE site_settings 
ADD COLUMN IF NOT EXISTS admin_sidebar_toggles JSONB NOT NULL DEFAULT '{}'::jsonb;

-- Create missing attractions table
CREATE TABLE IF NOT EXISTS attractions (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image_url TEXT,
    icon_name VARCHAR(100),
    tag VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_attractions_deleted_at ON attractions(deleted_at);
