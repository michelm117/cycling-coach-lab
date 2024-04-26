-- Migration for creating sessions and settings table
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS globalSettings;


CREATE TABLE sessions (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  user_id INT REFERENCES users(id) ON DELETE CASCADE -- Assuming each session belongs to a user
);

CREATE TABLE globalSettings
(
    SectionName VARCHAR(50),
    SettingName VARCHAR(50),
    SettingValue VARCHAR(1000),
    SettingType SMALLINT DEFAULT 0
);

INSERT INTO globalSettings (SectionName, SettingName, SettingValue, SettingType) VALUES ('APP', 'initialized', 'false', 2);
