## Train Booking System

### Introduction

This project is a train booking system developed using Golang, Gin framework, and PostgreSQL database. It allows users to register, login, add new trains, check seat availability between stations, and book seats on available trains.

### Setup

1. **Clone Repository**: Clone this repository to your local machine using the following command:

   ```
   git clone <repository_url>
   ```

2. **Install Dependencies**: Navigate to the project directory and install dependencies using:

   ```
   go mod download
   ```

3. **Database Setup**:

   - run `docker-compose up` cmd to start the postgres docker image, make sure you have docker installed on your computer.

4. **Environment Variables**:

   - Create a `.env` file in the root directory.
   - Set the `ADMIN_KEY` var with the help of which Admin APIs will work

5. **Run Migrations**: Apply database migrations to create tables. Make sure you have `golang-migrate` installed on your machine. Run the following command:

   ```
   migrate -path db/migration -database "postgres://postgres:postgres@localhost/ticketbooking" -verbose up
   ```

6. **Start Server**: Start the Gin server by running:
   ```
   go run main.go
   ```

### Test Data to create a Train

```
{
  "name": "Test Train 4",
  "departure_time": "2024-04-20T08:00:00Z",
  "arrival_time": "2024-04-20T20:00:00Z",
  "from": "indore",
  "to": "bhopal",
  "total_seats": 10,
  "fare": 500,
  "admin_key":"123456"
}
```
