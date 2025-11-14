-- migrations/000001_init_schema.down.sql

DROP TRIGGER IF EXISTS update_tech_categories_updated_at ON tech_categories;
DROP TRIGGER IF EXISTS update_countries_updated_at ON countries;
DROP TRIGGER IF EXISTS update_technologies_updated_at ON technologies;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_tech_name_search;
DROP INDEX IF EXISTS idx_country_code;
DROP INDEX IF EXISTS idx_tech_status;
DROP INDEX IF EXISTS idx_tech_year_deployed;
DROP INDEX IF EXISTS idx_tech_year_developed;
DROP INDEX IF EXISTS idx_tech_category;
DROP INDEX IF EXISTS idx_tech_country;

DROP TABLE IF EXISTS technologies;
DROP TABLE IF EXISTS tech_categories;
DROP TABLE IF EXISTS countries;

DROP EXTENSION IF EXISTS "uuid-ossp";
