BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (username);

COMMENT ON COLUMN users.id IS 'User ID';
COMMENT ON COLUMN users.username IS 'Login username';
COMMENT ON COLUMN users.password IS 'Password hash';
COMMENT ON COLUMN users.created_at IS 'Creation time';
COMMENT ON COLUMN users.updated_at IS 'Last update time';

CREATE TABLE IF NOT EXISTS projects (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    private_key TEXT NOT NULL,
    public_key TEXT NOT NULL,
    key_bits BIGINT NOT NULL DEFAULT 2048,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

COMMENT ON COLUMN projects.id IS 'Project ID';
COMMENT ON COLUMN projects.name IS 'Project name';
COMMENT ON COLUMN projects.description IS 'Project description';
COMMENT ON COLUMN projects.private_key IS 'Encrypted private key';
COMMENT ON COLUMN projects.public_key IS 'Public key';
COMMENT ON COLUMN projects.key_bits IS 'RSA key size';
COMMENT ON COLUMN projects.created_at IS 'Creation time';
COMMENT ON COLUMN projects.updated_at IS 'Last update time';

CREATE TABLE IF NOT EXISTS features (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_features_code ON features (code);

COMMENT ON COLUMN features.id IS 'Feature ID';
COMMENT ON COLUMN features.code IS 'Feature code';
COMMENT ON COLUMN features.name IS 'Feature name';
COMMENT ON COLUMN features.description IS 'Feature description';
COMMENT ON COLUMN features.enabled IS 'Whether feature is enabled';
COMMENT ON COLUMN features.created_at IS 'Creation time';
COMMENT ON COLUMN features.updated_at IS 'Last update time';

CREATE TABLE IF NOT EXISTS plans (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    max_users BIGINT NOT NULL DEFAULT 0,
    max_devices BIGINT NOT NULL DEFAULT 1,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_plans_code ON plans (code);

COMMENT ON COLUMN plans.id IS 'Plan ID';
COMMENT ON COLUMN plans.code IS 'Plan code';
COMMENT ON COLUMN plans.name IS 'Plan name';
COMMENT ON COLUMN plans.description IS 'Plan description';
COMMENT ON COLUMN plans.max_users IS 'Max user count, 0 means unlimited';
COMMENT ON COLUMN plans.max_devices IS 'Max device count';
COMMENT ON COLUMN plans.enabled IS 'Whether plan is enabled';
COMMENT ON COLUMN plans.created_at IS 'Creation time';
COMMENT ON COLUMN plans.updated_at IS 'Last update time';

CREATE TABLE IF NOT EXISTS plan_features (
    plan_id BIGINT NOT NULL,
    feature_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ,
    PRIMARY KEY (plan_id, feature_id),
    CONSTRAINT fk_plan_features_plan FOREIGN KEY (plan_id) REFERENCES plans (id) ON DELETE CASCADE,
    CONSTRAINT fk_plan_features_feature FOREIGN KEY (feature_id) REFERENCES features (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_plan_features_feature_id ON plan_features (feature_id);

COMMENT ON COLUMN plan_features.plan_id IS 'Plan ID';
COMMENT ON COLUMN plan_features.feature_id IS 'Feature ID';
COMMENT ON COLUMN plan_features.created_at IS 'Creation time';

CREATE TABLE IF NOT EXISTS licenses (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR(100) NOT NULL,
    project_id BIGINT NOT NULL,
    subject_name VARCHAR(255) NOT NULL,
    subject_email VARCHAR(255),
    subject_org VARCHAR(255),
    plan_id BIGINT,
    plan_name VARCHAR(255),
    plan VARCHAR(100),
    not_before TIMESTAMPTZ NOT NULL,
    not_after TIMESTAMPTZ NOT NULL,
    hardware JSON,
    features JSON NOT NULL,
    limits JSON NOT NULL DEFAULT '{"max_users":0,"max_devices":0}',
    metadata JSON NOT NULL,
    lic_content TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    CONSTRAINT fk_licenses_project FOREIGN KEY (project_id) REFERENCES projects (id),
    CONSTRAINT fk_licenses_plan FOREIGN KEY (plan_id) REFERENCES plans (id) ON DELETE SET NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_licenses_uuid ON licenses (uuid);
CREATE INDEX IF NOT EXISTS idx_licenses_project_id ON licenses (project_id);
CREATE INDEX IF NOT EXISTS idx_licenses_plan_id ON licenses (plan_id);

COMMENT ON COLUMN licenses.id IS 'License ID';
COMMENT ON COLUMN licenses.uuid IS 'License UUID';
COMMENT ON COLUMN licenses.project_id IS 'Project ID';
COMMENT ON COLUMN licenses.subject_name IS 'Licensed subject name';
COMMENT ON COLUMN licenses.subject_email IS 'Licensed subject email';
COMMENT ON COLUMN licenses.subject_org IS 'Licensed subject organization';
COMMENT ON COLUMN licenses.plan_id IS 'Plan ID snapshot source';
COMMENT ON COLUMN licenses.plan_name IS 'License plan name snapshot';
COMMENT ON COLUMN licenses.plan IS 'License plan';
COMMENT ON COLUMN licenses.not_before IS 'Valid from time';
COMMENT ON COLUMN licenses.not_after IS 'Valid until time';
COMMENT ON COLUMN licenses.hardware IS 'Hardware binding JSON';
COMMENT ON COLUMN licenses.features IS 'Licensed features JSON';
COMMENT ON COLUMN licenses.limits IS 'License limits JSON';
COMMENT ON COLUMN licenses.metadata IS 'License metadata JSON';
COMMENT ON COLUMN licenses.lic_content IS 'Signed license content';
COMMENT ON COLUMN licenses.created_at IS 'Creation time';

COMMIT;
