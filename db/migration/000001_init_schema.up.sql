CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "email" varchar,
  "name" varchar,
  "password" varchar,
  "created_at" timestamp
);

CREATE TABLE "trains" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar,
  "departure_time" timestamp,
  "arrival_time" timestamp,
  "from" varchar,
  "to" varchar,
  "total_seats" integer,
  "fare" integer,
  "created_at" timestamp
);

CREATE TABLE "bookings" (
  "id" SERIAL PRIMARY KEY,
  "user_id" integer,
  "train_id" integer,
  "seats" integer,
  "status" varchar,
  "created_at" timestamp
);

ALTER TABLE "bookings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "bookings" ADD FOREIGN KEY ("train_id") REFERENCES "trains" ("id");