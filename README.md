# CBC Backend

A RESTful API backend service for managing CBC (Competency-Based Curriculum) educational resources and related functionalities.

## Overview

This backend service provides APIs for managing educational resources, user accounts, and job postings. It's built using the Beego framework and PostgreSQL database, focusing on serving CBC-related content to educational institutions.

## Features

### Resource Management
- Upload educational resources (PDF, DOC, DOCX, TXT, RTF)
- Search resources by title, level, and description
- Paginated resource listing (20 items per page)
- File type validation
- Automatic file storage management

### User Management
- User registration and authentication
- Profile management
- Role-based access control

### Job Portal
- Post teaching and educational jobs
- Search and filter job listings
- Job application management

## Technical Stack

- **Framework**: Beego v2.1.0
- **Database**: PostgreSQL
- **Language**: Go 1.23
- **Authentication**: JWT (JSON Web Tokens)
- **Documentation**:

## API Endpoints

### Resources
GET /v1/resources
Query Parameters:
level: Filter by educational level
title: Search in resource titles
description: Search in resource descriptions
page: Page number for pagination (default: 1)
Response includes pagination metadata

POST /v1/resources
Multipart form data:
title: Resource title
description: Resource description
level: Educational level
file: Resource file (PDF/DOC/DOCX/TXT/RTF only)

### Users
POST /v1/user/signup
POST /v1/user/login
GET /v1/user/logout
PUT /v1/user/:uid
DELETE /v1/user/:uid

### Jobs
GET /v1/jobs
POST /v1/jobs


## Setup and Installation

1. Prerequisites:
   ```bash
   - Go 1.23 or higher
   - PostgreSQL
   ```

2. Database Setup:
   ```sql
   CREATE DATABASE cbcexams;
   ```

3. Configuration:
   ```bash
   # Update conf/app.conf with your database credentials
   sqlconn = "user=postgres password=your_password dbname=cbcexams sslmode=disable"
   ```

4. Install Dependencies:
   ```bash
   go mod tidy
   ```

5. Run the Application:
   ```bash
   go run main.go
   ```

## Project Structure

├── conf/
│ └── app.conf
├── controllers/
│ ├── resource.go
│ ├── user.go
│ └── job.go
├── models/
│ ├── db.go
│ ├── resource.go
│ ├── user.go
│ └── job.go
├── routers/
│ └── router.go
├── uploads/ # Resource files storage
├── main.go
└── README.md

## Development

- Run tests:
  ```bash
  go test ./tests
  ```

- Generate Swagger documentation:
  ```bash
  bee generate docs
  ```

## API Documentation

To be added


