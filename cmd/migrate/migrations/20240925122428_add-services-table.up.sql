-- Create table
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,  -- Use SERIAL for auto-incrementing IDs
    barber_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NULL,
    price DECIMAL(10, 2) NOT NULL,
    duration INT NOT NULL,  -- Duration in minutes
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key reference
    FOREIGN KEY (barber_id) REFERENCES barbers(user_id)
);

-- Add a comment for the duration column
COMMENT ON COLUMN services.duration IS 'Duration in minutes';

-- Add a trigger to update the updated_at field on row updates
CREATE OR REPLACE FUNCTION update_timestamp_services()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger for the services table
CREATE TRIGGER update_services_updated_at
BEFORE UPDATE ON services
FOR EACH ROW
EXECUTE FUNCTION update_timestamp_services();
