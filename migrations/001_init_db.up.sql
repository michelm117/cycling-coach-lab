------------------------------------------------------------
--                        Enums                           --
------------------------------------------------------------
--CREATE TYPE IF NOT EXISTS user_roles AS ENUM ('admin', 'athlete');
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_roles') THEN
        CREATE TYPE user_roles AS ENUM ('admin', 'athlete');
    END IF;
END
$$;

-- CREATE TYPE user_status AS ENUM ('active', 'inactive');
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status') THEN
        CREATE TYPE user_status AS ENUM ('active', 'inactive');
    END IF;
END
$$;

------------------------------------------------------------
--                   Users Table                          --
------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
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

INSERT INTO users (email, firstname, lastname, date_of_birth, password_hash, status, role, created_at, updated_at)
VALUES
  ('jan@ullrich.de', 'Jan', 'Ullrich', '1973-12-02', 'hash456', 'inactive', 'athlete', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('mathieu@van-der-pole.be', 'Mathieu', 'van der Poel', '1995-01-19', 'hash456', 'active', 'athlete', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('admin@cycling-coach-lab.de', 'Admin', 'Admin', '1990-01-01', '$2a$10$zxfv3JU3K/YvMrt3AvyCG.fM2PHKMTytAD2Oa4SqBdYZpD9FlDLqO', 'active', 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);


------------------------------------------------------------
--                   Sessions Table                       --
------------------------------------------------------------
CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  user_id INT REFERENCES users(id) ON DELETE CASCADE 
);


------------------------------------------------------------
--               Global Settings Table                    --
------------------------------------------------------------
CREATE TABLE IF NOT EXISTS globalSettings
(
    SectionName VARCHAR(50),
    SettingName VARCHAR(50),
    SettingValue VARCHAR(1000),
    SettingType SMALLINT DEFAULT 1,

    PRIMARY KEY (SectionName, SettingName)
);

-- Create triggers to enforce lowercase on insert and update
CREATE OR REPLACE FUNCTION lowercase_setting_name()
RETURNS TRIGGER AS $$
BEGIN
    NEW.SectionName := LOWER(NEW.SectionName);
    NEW.SettingName := LOWER(NEW.SettingName);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS lowercase_setting_name_trigger ON globalSettings;
CREATE TRIGGER lowercase_setting_name_trigger
BEFORE INSERT OR UPDATE ON globalSettings
FOR EACH ROW
EXECUTE FUNCTION lowercase_setting_name();

INSERT INTO globalSettings (SectionName, SettingName, SettingValue, SettingType) VALUES ('APP', 'initialized', 'false', 2);
