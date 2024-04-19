DROP TABLE users

DROP TABLE trains

DROP TABLE bookings

ALTER TABLE bookings DROP CONSTRAINT bookings_user_id_fkey;

ALTER TABLE bookings DROP CONSTRAINT bookings_train_id_fkey;