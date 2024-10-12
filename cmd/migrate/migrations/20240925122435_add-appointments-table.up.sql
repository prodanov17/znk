-- Create ENUM type for appointment status
CREATE TYPE appointment_status_enum AS ENUM ('pending', 'approved', 'rejected', 'completed', 'canceled');

-- Create table
CREATE TABLE IF NOT EXISTS appointments (
    id SERIAL PRIMARY KEY,  -- Use SERIAL for auto-incrementing IDs
    barber_id SERIAL NOT NULL,
    user_id SERIAL DEFAULT NULL,  -- Nullable for guest users
    name VARCHAR(255) NOT NULL,  -- Name of the client
    date DATE NOT NULL,
    time TIME NOT NULL,
    status appointment_status_enum NOT NULL DEFAULT 'pending',  -- Use the custom ENUM type
    service_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key references
    FOREIGN KEY (barber_id) REFERENCES barbers(user_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (service_id) REFERENCES services(id)
);

-- Add a trigger to update the updated_at field on row updates
CREATE OR REPLACE FUNCTION update_timestamp_appointments()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger for the appointments table
CREATE TRIGGER update_appointments_updated_at
BEFORE UPDATE ON appointments
FOR EACH ROW
EXECUTE FUNCTION update_timestamp_appointments();
