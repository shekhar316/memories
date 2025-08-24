-- Create database (already created by POSTGRES_DB, but kept for clarity)
CREATE EXTENSION IF NOT EXISTS vector;

CREATE DATABASE memories;

-- Grant permissions to user
GRANT ALL PRIVILEGES ON DATABASE memories TO postgres;
