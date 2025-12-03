-- migrations/000001_init_schema.up.sql

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- countries table
CREATE TABLE IF NOT EXISTS countries(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    code VARCHAR(3) NOT NULL UNIQUE,
    flag_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Technology categories
CREATE TABLE IF NOT EXISTS tech_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    icon_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Military technologies
CREATE TABLE IF NOT EXISTS technologies (
    id SERIAL PRIMARY KEY,
    country_id INTEGER NOT NULL REFERENCES countries(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES tech_categories(id) ON DELETE RESTRICT,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    designer VARCHAR(200),
    year_developed INTEGER,
    year_deployed INTEGER,
    manufacturer VARCHAR(250),
    unit_cost BIGINT,
    mass INTEGER,
    length INTEGER,
    width INTEGER,
    height INTEGER,
    status VARCHAR(50) CHECK (status IN ('historical', 'current', 'future', 'concept', 'prototype')),
    image_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_tech_country ON technologies(country_id);
CREATE INDEX IF NOT EXISTS idx_tech_category ON technologies(category_id);
CREATE INDEX IF NOT EXISTS idx_tech_year_developed ON technologies(year_developed);
CREATE INDEX IF NOT EXISTS idx_tech_year_deployed ON technologies(year_deployed);
CREATE INDEX IF NOT EXISTS idx_tech_status ON technologies(status);
CREATE INDEX IF NOT EXISTS idx_country_code ON countries(code);
CREATE INDEX IF NOT EXISTS idx_tech_name_search ON technologies USING gin(to_tsvector('english', name));

-- Updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';


-- Apply triggers
DROP TRIGGER IF EXISTS update_technologies_updated_at ON technologies;
CREATE TRIGGER update_technologies_updated_at
    BEFORE UPDATE ON technologies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_countries_updated_at ON countries;
CREATE TRIGGER update_countries_updated_at
    BEFORE UPDATE ON countries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_tech_categories_updated_at ON tech_categories;
CREATE TRIGGER update_tech_categories_updated_at
    BEFORE UPDATE ON tech_categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
