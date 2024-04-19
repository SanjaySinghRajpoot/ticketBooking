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

### Features

1. **User Registration**: Users can register with a unique username and password.

2. **User Login**: Registered users can log in to their accounts securely.

3. **Admin Privileges**: Admin users have special privileges like adding new trains, updating total seats, etc.

4. **Add a New Train**: Admin users can add new trains with their source, destination, and total available seats.

5. **Get Seat Availability**: Users can check seat availability between any two stations.

6. **Book a Seat**: Logged-in users can book seats on available trains.

7. **Real-time Updates**: Seat availability updates in real-time, ensuring accurate information for users.

8. **Role-Based Access Control**: Differentiate between admin and regular users to provide appropriate access.

9. **Concurrency Handling**: Prevent race conditions while booking seats to ensure data integrity and consistency.

10. **Security**: Authentication and authorization mechanisms in place to secure user data and transactions.

11. **Error Handling**: Comprehensive error handling to provide meaningful feedback to users in case of failures.

12. **Scalability**: Optimized codebase and database design to handle large traffic without compromising performance.


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
