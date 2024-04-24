DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_roles; 
DROP TYPE IF EXISTS user_status; 

CREATE TYPE user_roles AS ENUM ('admin', 'athlete');
CREATE TYPE user_status AS ENUM ('active', 'inactive');

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(100) NOT NULL UNIQUE,
  firstname VARCHAR(50),
  lastname VARCHAR(50),
  date_of_birth DATE,
  password_hash VARCHAR(100),
  status user_status,
  role user_roles,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Example data with created_at and updated_at fields
INSERT INTO users (email, firstname, lastname, date_of_birth, password_hash, status, role, created_at, updated_at)
VALUES
  ('admin@example.com', 'admin', 'admin', '1990-01-01', 'hash123', 'active', 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('jan@ullrich.de', 'Jan', 'Ullrich', '1973-12-02', 'hash456', 'inactive', 'athlete', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('mathieu@van-der-pole.be', 'Mathieu', 'van der Poel', '1995-01-19', 'hash456', 'active', 'athlete', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
