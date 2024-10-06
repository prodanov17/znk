CREATE TABLE IF NOT EXISTS barber_unavailabilities (
    barber_id INT PRIMARY KEY,
    FOREIGN KEY (barber_id) REFERENCES barbers(user_id),
    date DATE NOT NULL,
    start_time TIMESTAMP COMMENT 'Null for the entire day',
    end_time TIMESTAMP COMMENT 'Null for the entire day',
    reason VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
