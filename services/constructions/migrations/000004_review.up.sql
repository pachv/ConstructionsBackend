CREATE TABLE reviews (
    id VARCHAR(255),
    name VARCHAR(255),
    position VARCHAR(255),
    text VARCHAR(2000),
    rating INT,
    image_path VARCHAR(1024),
    consent BOOLEAN,
    can_publish BOOLEAN,
    created_at TIMESTAMP
);
