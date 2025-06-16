-- +goose Up
CREATE TABLE attachments (
                             id SERIAL PRIMARY KEY,
                             name_original VARCHAR(255) NOT NULL,
                             name_hashed VARCHAR(255) NOT NULL,
                             url VARCHAR(255) NOT NULL UNIQUE,
                             lesson_id INT REFERENCES lessons(id) ON DELETE CASCADE,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS attachments;
