-- Create ENUM type for days of the week
CREATE TYPE day_of_week_enum AS ENUM ('SUNDAY', 'MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY', 'SATURDAY');

-- Create table
CREATE TABLE IF NOT EXISTS barber_availabilities (
    barber_id INT PRIMARY KEY,  -- barber_id as primary key
    day_of_week day_of_week_enum NOT NULL,  -- Use the custom ENUM type for days
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key reference
    FOREIGN KEY (barber_id) REFERENCES barbers(user_id)
);

-- Add a trigger to update the updated_at field on row updates
CREATE OR REPLACE FUNCTION update_timestamp_barber_availabilities()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger for the barber_availabilities table
CREATE TRIGGER update_barber_availabilities_updated_at
BEFORE UPDATE ON barber_availabilities
FOR EACH ROW
EXECUTE FUNCTION update_timestamp_barber_availabilities();
