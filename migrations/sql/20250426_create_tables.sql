-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied.

CREATE TABLE courses (
                         id SERIAL PRIMARY KEY,
                         name VARCHAR(255) NOT NULL UNIQUE,
                         description TEXT,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE chapters (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL UNIQUE,
                          description TEXT,
                          "order" INT,
                          course_id INT REFERENCES courses(id) ON DELETE CASCADE,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE lessons (
                         id SERIAL PRIMARY KEY,
                         name VARCHAR(255) NOT NULL UNIQUE,
                         description TEXT,
                         content TEXT,
                         "order" INT,
                         chapter_id INT REFERENCES chapters(id) ON DELETE CASCADE,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back.

DROP TABLE IF EXISTS lessons;
DROP TABLE IF EXISTS chapters;
DROP TABLE IF EXISTS courses;
