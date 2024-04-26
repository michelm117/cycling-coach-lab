-- Reset sessions table
DELETE FROM sessions;

-- Reset users table
DELETE FROM users;
-- INSERT INTO users (email, firstname, lastname, date_of_birth, password_hash, status, role, created_at, updated_at)
-- VALUES
--   ('admin@example.com', 'admin', 'admin', '1990-01-01', 'hash123', 'active', 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
--   ('jan@ullrich.de', 'Jan', 'Ullrich', '1973-12-02', 'hash456', 'inactive', 'athlete', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
--   ('mathieu@van-der-pole.be', 'Mathieu', 'van der Poel', '1995-01-19', 'hash456', 'active', 'athlete', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
--   ('john@doe.com', 'John', 'Doe', '1990-01-01', '$2a$10$NsbCtrtJiJ/ynNEYLuiIZO/7wPKWt4mCpkfcJtjLd2k7/wObMZm9m', 'active', 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Reset globalSettings table
DELETE FROM globalSettings;
-- INSERT INTO globalSettings (SectionName, SettingName, SettingValue, SettingType) VALUES ('APP', 'initialized', 'true', 2);

