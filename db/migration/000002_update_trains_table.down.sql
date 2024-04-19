-- Add deleted_at column to users table
ALTER TABLE users
ADD COLUMN deleted_at TIMESTAMP;

-- Add deleted_at column to trains table
ALTER TABLE trains
ADD COLUMN deleted_at TIMESTAMP;

-- Add deleted_at column to bookings table
ALTER TABLE bookings
ADD COLUMN deleted_at TIMESTAMP;

-- Add status column to trains table
ALTER TABLE trains
ADD COLUMN status VARCHAR;
