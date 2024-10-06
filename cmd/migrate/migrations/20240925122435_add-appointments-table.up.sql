CREATE TABLE IF NOT EXISTS appointments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    barber_id INT NOT NULL,
    user_id INT DEFAULT NULL, -- Nullable for guest users
    name VARCHAR(255) NOT NULL, -- Name of the client
    date DATE NOT NULL,
    time TIME NOT NULL,
    status ENUM('pending', 'approved', 'rejected', 'completed', 'canceled') NOT NULL DEFAULT 'pending',
    service_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (barber_id) REFERENCES barbers(user_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (service_id) REFERENCES services(id)
);
