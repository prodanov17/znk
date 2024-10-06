CREATE TABLE IF NOT EXISTS barbers (
    user_id INT NOT NULL,
    profile_picture VARCHAR(255),
    establishment_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (establishment_id) REFERENCES establishments(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
