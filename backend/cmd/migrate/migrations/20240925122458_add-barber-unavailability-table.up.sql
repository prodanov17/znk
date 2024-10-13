-- Create table
CREATE TABLE IF NOT EXISTS barber_unavailabilities (
    barber_id INT PRIMARY KEY,  -- barber_id as primary key
    date DATE NOT NULL,
    start_time TIMESTAMP,  -- Null for the entire day
    end_time TIMESTAMP,    -- Null for the entire day
    reason VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key reference
    FOREIGN KEY (barber_id) REFERENCES barbers(user_id)
);

-- Add a comment for the start_time and end_time columns
COMMENT ON COLUMN barber_unavailabilities.start_time IS 'Null for the entire day';
COMMENT ON COLUMN barber_unavailabilities.end_time IS 'Null for the entire day';

-- Add a trigger to update the updated_at field on row updates
CREATE OR REPLACE FUNCTION update_timestamp_barber_unavailabilities()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger for the barber_unavailabilities table
CREATE TRIGGER update_barber_unavailabilities_updated_at
BEFORE UPDATE ON barber_unavailabilities
FOR EACH ROW
EXECUTE FUNCTION update_timestamp_barber_unavailabilities();
