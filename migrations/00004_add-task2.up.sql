-- Drop the table if it exists
DROP TABLE IF EXISTS tasks;

-- Create tasks table
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    content VARCHAR(500) NOT NULL
);

-- Insert initial data
INSERT INTO tasks (title, content) VALUES ('title1', 'content1');
INSERT INTO tasks (title, content) VALUES ('title2', 'content2');
INSERT INTO users (username, email) VALUES ('user3', 'user3@example.com');
