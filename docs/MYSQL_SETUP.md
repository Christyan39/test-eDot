# MySQL Setup Guide

## Prerequisites

1. **Install MySQL Server** (if not already installed):
   ```bash
   # macOS using Homebrew
   brew install mysql
   
   # Or download from: https://dev.mysql.com/downloads/mysql/
   ```

2. **Start MySQL Service**:
   ```bash
   # macOS
   brew services start mysql
   
   # Or manually
   sudo mysql.server start
   ```

## Database Setup

1. **Connect to MySQL**:
   ```bash
   mysql -u root -p
   ```

2. **Create Database and User**:
   ```sql
   CREATE DATABASE edot_user CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   
   -- Optional: Create dedicated user
   CREATE USER 'edot_user'@'localhost' IDENTIFIED BY 'password';
   GRANT ALL PRIVILEGES ON edot_user.* TO 'edot_user'@'localhost';
   FLUSH PRIVILEGES;
   ```

3. **Create Users Table**:
   ```sql
   USE edot_user;
   
   CREATE TABLE IF NOT EXISTS users (
       id INT AUTO_INCREMENT PRIMARY KEY,
       name VARCHAR(255) NOT NULL,
       email VARCHAR(255) UNIQUE NOT NULL,
       phone VARCHAR(20),
       password VARCHAR(255) NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
   );
   ```

## Environment Configuration

1. **Copy environment file**:
   ```bash
   cp .env.example .env
   ```

2. **Edit .env file** with your MySQL credentials:
   ```env
   # Server Configuration
   PORT=8080

   # MySQL Database Configuration
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=password
   DB_NAME=edot_user

   # JWT Configuration
   JWT_SECRET=your-super-secret-jwt-key-change-in-production

   # Environment
   ENV=development
   ```

## Running the Application

1. **Start the server**:
   ```bash
   go run main.go
   ```

2. **Test the connection**:
   The application will automatically try to connect to MySQL and show connection status in the console.

## API Testing

1. **Health Check**:
   ```bash
   curl http://localhost:8080/health
   ```

2. **Create User**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/users \
     -H "Content-Type: application/json" \
     -d '{"name":"John Doe","email":"john@example.com"}'
   ```

3. **Get Users**:
   ```bash
   curl http://localhost:8080/api/v1/users
   ```

## Troubleshooting

- **Connection Refused**: Check if MySQL is running
- **Access Denied**: Verify username/password in .env
- **Database Not Found**: Create the database first
- **Port Already in Use**: Change PORT in .env file