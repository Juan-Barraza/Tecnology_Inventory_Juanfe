-- ============================================================
-- INVENTORY MANAGEMENT SYSTEM — Schema
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE logical_status_enum  AS ENUM ('active','inactive','written_off');
CREATE TYPE physical_status_enum AS ENUM ('optimal','good','fair','deteriorated','out_of_service');
CREATE TYPE assignment_status_enum AS ENUM ('active','released','written_off');
CREATE TYPE period_status_enum AS ENUM ('open', 'closed');

-- ── Catalogs ──────────────────────────────────────────────────
CREATE TABLE cities (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    department VARCHAR(100) NOT NULL,
    CONSTRAINT uq_cities_name UNIQUE (name)
);

CREATE TABLE areas (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    CONSTRAINT uq_areas_name UNIQUE (name)
);

-- Accounting group: the parent (code is unique)
CREATE TABLE accounting_groups (
    id         SERIAL PRIMARY KEY,
    code       BIGINT       NOT NULL,
    name       VARCHAR(150) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_accounting_groups_code UNIQUE (code)
);

-- Asset account: the specific sub-account (account_code is unique)
CREATE TABLE asset_accounts (
    id                  SERIAL PRIMARY KEY,
    accounting_group_id INT          NOT NULL REFERENCES accounting_groups(id),
    account_code        BIGINT       NOT NULL,
    open_ledger         VARCHAR(100),
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_asset_accounts_code UNIQUE (account_code)
);

CREATE TABLE asset_categories (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    CONSTRAINT uq_asset_categories_name UNIQUE (name)
);

-- ── Users ─────────────────────────────────────────────────────
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(150) NOT NULL,
    email         VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    last_login    TIMESTAMP WITH TIME ZONE,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_users_email UNIQUE (email)
);

-- ── Assets ────────────────────────────────────────────────────
CREATE TABLE assets (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code             VARCHAR(50)  NOT NULL,
    description      TEXT         NOT NULL,
    category_id      INT          NOT NULL REFERENCES asset_categories(id),
    asset_account_id INT          NOT NULL REFERENCES asset_accounts(id),
    city_id          INT          NOT NULL REFERENCES cities(id),
    area_id          INT          REFERENCES areas(id),
    historical_cost  NUMERIC(14,2),
    activation_date  DATE         NOT NULL,
    logical_status   logical_status_enum  NOT NULL DEFAULT 'active',
    physical_status  physical_status_enum NOT NULL DEFAULT 'optimal',
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_assets_code UNIQUE (code)
);

-- ── Assignments ───────────────────────────────────────────────
CREATE TABLE assignments (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id             UUID NOT NULL REFERENCES assets(id),
    responsible_name     VARCHAR(150),
    responsible_position VARCHAR(150),
    assigned_at          DATE NOT NULL,
    deactivated_at       DATE,
    deactivation_reason  TEXT,
    status               assignment_status_enum NOT NULL DEFAULT 'active',
    created_by           UUID NOT NULL REFERENCES users(id),
    created_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_deactivation CHECK (
        (status = 'active'  AND deactivated_at IS NULL) OR
        (status <> 'active' AND deactivated_at IS NOT NULL)
    )
);

-- ── Status history ────────────────────────────────────────────
CREATE TABLE status_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id        UUID NOT NULL REFERENCES assets(id),
    previous_status logical_status_enum,
    new_status      logical_status_enum NOT NULL,
    notes           TEXT,
    recorded_by     UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE inventory_periods (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    period_year  INT NOT NULL,
    period_month INT NOT NULL CHECK (period_month BETWEEN 1 AND 12),
    period_day   INT NOT NULL CHECK (period_day BETWEEN 1 AND 31),
    status       period_status_enum NOT NULL DEFAULT 'open',
    created_by   UUID NOT NULL REFERENCES users(id),
    closed_at    TIMESTAMP WITH TIME ZONE,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_period UNIQUE (period_year, period_month, period_day)
);

CREATE TABLE inventory_records (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    period_id   UUID    NOT NULL REFERENCES inventory_periods(id),
    asset_id    UUID    NOT NULL REFERENCES assets(id),
    confirmed   BOOLEAN NOT NULL DEFAULT false,
    deactivated BOOLEAN NOT NULL DEFAULT false,
    notes       TEXT,
    has_label   BOOLEAN NOT NULL DEFAULT false,
    recorded_by UUID NOT NULL REFERENCES users(id),
    recorded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_record_per_period UNIQUE (period_id, asset_id)
);

-- ── Indexes ───────────────────────────────────────────────────
CREATE INDEX idx_assets_code             ON assets(code);
CREATE INDEX idx_assets_logical_status   ON assets(logical_status);
CREATE INDEX idx_assets_physical_status  ON assets(physical_status);
CREATE INDEX idx_assets_activation_date  ON assets(activation_date);
CREATE INDEX idx_assets_city_id          ON assets(city_id);
CREATE INDEX idx_assets_area_id          ON assets(area_id);
CREATE INDEX idx_assets_category_id      ON assets(category_id);
CREATE INDEX idx_assets_asset_account_id ON assets(asset_account_id);
CREATE INDEX idx_assignments_asset_id    ON assignments(asset_id);
CREATE INDEX idx_assignments_status      ON assignments(status);
CREATE INDEX idx_assignments_dates       ON assignments(assigned_at, deactivated_at);
CREATE INDEX idx_status_history_asset    ON status_history(asset_id);
CREATE INDEX idx_status_history_date     ON status_history(created_at);
CREATE INDEX idx_inventory_records_period  ON inventory_records(period_id);
CREATE INDEX idx_inventory_records_asset   ON inventory_records(asset_id);

-- ── updated_at trigger ────────────────────────────────────────
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN NEW.updated_at = NOW(); RETURN NEW; END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at          BEFORE UPDATE ON users          FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_assets_updated_at         BEFORE UPDATE ON assets         FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_accounting_groups_updated BEFORE UPDATE ON accounting_groups FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_asset_accounts_updated    BEFORE UPDATE ON asset_accounts  FOR EACH ROW EXECUTE FUNCTION set_updated_at();