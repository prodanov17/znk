-- Create table
CREATE TABLE IF NOT EXISTS barbers (
    user_id SERIAL NOT NULL,
    profile_picture VARCHAR(255),
    establishment_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key references
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (establishment_id) REFERENCES establishments(id)
);

-- Add a trigger to update the updated_at field on row updates
CREATE OR REPLACE FUNCTION update_timestamp_barbers()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger for the barbers table
CREATE TRIGGER update_barbers_updated_at
BEFORE UPDATE ON barbers
FOR EACH ROW
EXECUTE FUNCTION update_timestamp_barbers();
