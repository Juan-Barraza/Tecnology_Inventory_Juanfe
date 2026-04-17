-- ============================================================
-- Migration: Add new categories + owner column on assets
-- ============================================================

-- 1) New asset categories
INSERT INTO asset_categories (name) VALUES
    ('Electrodomestico'),
    ('Escritorio'),
    ('Archivador'),
    ('Sillas'),
    ('Descansa Pies')
ON CONFLICT (name) DO NOTHING;

-- 2) Add "owner" (VARCHAR) column to assets table
ALTER TABLE assets
    ADD COLUMN IF NOT EXISTS owner VARCHAR(150);

-- 3) Make asset fields nullable (no longer required)
ALTER TABLE assets ALTER COLUMN code DROP NOT NULL;
ALTER TABLE assets ALTER COLUMN description DROP NOT NULL;
ALTER TABLE assets ALTER COLUMN category_id DROP NOT NULL;
ALTER TABLE assets ALTER COLUMN asset_account_id DROP NOT NULL;
ALTER TABLE assets ALTER COLUMN city_id DROP NOT NULL;
ALTER TABLE assets ALTER COLUMN activation_date DROP NOT NULL;

-- 4) New Area
INSERT INTO areas (name) VALUES
    ('Logistica')
ON CONFLICT (name) DO NOTHING;