CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS admin_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(120) NOT NULL,
    email VARCHAR(180) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(30) NOT NULL DEFAULT 'staff' CHECK (role IN ('super_admin', 'admin', 'moderator', 'staff')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(160) NOT NULL,
    slug VARCHAR(180) NOT NULL UNIQUE,
    description TEXT,
    price BIGINT NOT NULL CHECK (price >= 0),
    original_price BIGINT CHECK (original_price >= 0),
    category VARCHAR(80),
    features JSONB NOT NULL DEFAULT '[]'::jsonb,
    validity VARCHAR(120),
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    sold_count INTEGER NOT NULL DEFAULT 0 CHECK (sold_count >= 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id VARCHAR(40) NOT NULL UNIQUE,
    customer_name VARCHAR(140) NOT NULL,
    customer_email VARCHAR(180) NOT NULL,
    customer_phone VARCHAR(30) NOT NULL,
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    amount BIGINT NOT NULL CHECK (amount >= 0),
    payment_status VARCHAR(30) NOT NULL DEFAULT 'pending' CHECK (payment_status IN ('pending', 'paid', 'failed', 'refunded')),
    razorpay_order_id VARCHAR(120),
    razorpay_payment_id VARCHAR(120),
    razorpay_signature VARCHAR(255),
    visit_date DATE NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'cancelled', 'refunded')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS contact_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(140) NOT NULL,
    email VARCHAR(180) NOT NULL,
    phone VARCHAR(30),
    subject VARCHAR(180),
    message TEXT NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'read', 'replied', 'archived')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS site_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_name VARCHAR(140) NOT NULL,
    logo_url TEXT,
    contact_email VARCHAR(180),
    phone_numbers JSONB NOT NULL DEFAULT '[]'::jsonb,
    whatsapp_number VARCHAR(20),
    address TEXT,
    google_maps_url TEXT,
    opening_hours TEXT,
    social_links JSONB NOT NULL DEFAULT '{}'::jsonb,
    meta_title VARCHAR(180),
    meta_description TEXT,
    razorpay_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    maintenance_mode BOOLEAN NOT NULL DEFAULT FALSE,
    feature_toggles JSONB NOT NULL DEFAULT '{}'::jsonb,
    homepage_sections JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_user_id UUID REFERENCES admin_users(id) ON UPDATE CASCADE ON DELETE SET NULL,
    action VARCHAR(120) NOT NULL,
    module VARCHAR(120) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    ip_address VARCHAR(80),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS hero_slides (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255),
    subtitle VARCHAR(255),
    description TEXT,
    image_url TEXT NOT NULL,
    mobile_image_url TEXT,
    cta_label VARCHAR(100),
    cta_url VARCHAR(255),
    secondary_cta_label VARCHAR(100),
    secondary_cta_url VARCHAR(255),
    badge_text VARCHAR(100),
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS content_pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slug VARCHAR(180) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    meta_title VARCHAR(255),
    meta_description TEXT,
    is_published BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS media_assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url TEXT NOT NULL,
    storage_key VARCHAR(255) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    original_filename VARCHAR(255),
    mime_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT NOT NULL,
    storage_provider VARCHAR(50) NOT NULL,
    uploaded_by_id UUID NOT NULL REFERENCES admin_users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    alt_text VARCHAR(255),
    folder VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS seo_pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_slug VARCHAR(180) NOT NULL UNIQUE,
    meta_title VARCHAR(255),
    meta_description TEXT,
    canonical_url TEXT,
    og_title VARCHAR(255),
    og_description TEXT,
    og_image TEXT,
    robots_index BOOLEAN NOT NULL DEFAULT TRUE,
    robots_follow BOOLEAN NOT NULL DEFAULT TRUE,
    schema_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gallery_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255),
    description TEXT,
    image_url TEXT NOT NULL,
    category VARCHAR(100),
    alt_text VARCHAR(255),
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS restaurant_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image_url TEXT,
    category VARCHAR(100),
    price BIGINT,
    is_featured BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS suite_rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    image_url TEXT NOT NULL,
    gallery JSONB NOT NULL DEFAULT '[]'::jsonb,
    price_per_night BIGINT NOT NULL,
    max_guests INTEGER NOT NULL DEFAULT 2,
    amenities JSONB NOT NULL DEFAULT '[]'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS hall_packages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image_url TEXT NOT NULL,
    capacity INTEGER,
    starting_price BIGINT,
    suitable_for JSONB NOT NULL DEFAULT '[]'::jsonb,
    features JSONB NOT NULL DEFAULT '[]'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS hall_enquiries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    email VARCHAR(255),
    event_type VARCHAR(100),
    expected_guests INTEGER,
    preferred_date DATE,
    message TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'new',
    source VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS offers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    code VARCHAR(50) NOT NULL UNIQUE,
    discount_type VARCHAR(50) NOT NULL,
    discount_value BIGINT NOT NULL,
    starts_at TIMESTAMPTZ,
    ends_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Indices
CREATE INDEX IF NOT EXISTS idx_admin_users_deleted_at ON admin_users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_tickets_active_sort ON tickets(is_active, sort_order);
CREATE INDEX IF NOT EXISTS idx_tickets_deleted_at ON tickets(deleted_at);
CREATE INDEX IF NOT EXISTS idx_bookings_search ON bookings(booking_id, customer_email, customer_phone);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status, payment_status);
CREATE INDEX IF NOT EXISTS idx_bookings_visit_date ON bookings(visit_date);
CREATE INDEX IF NOT EXISTS idx_bookings_deleted_at ON bookings(deleted_at);
CREATE INDEX IF NOT EXISTS idx_contact_messages_status ON contact_messages(status);
CREATE INDEX IF NOT EXISTS idx_contact_messages_deleted_at ON contact_messages(deleted_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_module_action ON audit_logs(module, action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_admin ON audit_logs(admin_user_id);
CREATE INDEX IF NOT EXISTS idx_hero_slides_sort ON hero_slides(is_active, sort_order);
CREATE INDEX IF NOT EXISTS idx_media_assets_storage ON media_assets(storage_key, storage_provider);
CREATE INDEX IF NOT EXISTS idx_gallery_category ON gallery_items(category, is_active);
CREATE INDEX IF NOT EXISTS idx_restaurant_category ON restaurant_items(category, is_active);
CREATE INDEX IF NOT EXISTS idx_hall_enquiries_phone ON hall_enquiries(phone);
