# CBC Backend

A RESTful API backend service for managing CBC (Competency-Based Curriculum) educational resources and related functionalities.

## Overview

This backend service provides APIs for managing educational resources, user accounts, and job postings. It's built using the Beego framework and PostgreSQL database, focusing on serving CBC-related content to educational institutions.

---

## Features

### Resource Management
- Upload educational resources (PDF, DOC, DOCX, TXT, RTF).
- Search resources by title, level, and description.
- Paginated resource listing (20 items per page).
- File type validation.
- Automatic file storage management.

### User Management
- User registration and authentication.
- Profile management.
- Role-based access control.

### Job Portal
- Post teaching and educational jobs.
- Search and filter job listings.
- Manage job applications.

---

## Technical Stack

- **Framework**: Beego v2.1.0
- **Database**: PostgreSQL
- **Language**: Go 1.23
- **Authentication**: JWT (JSON Web Tokens)
- **Documentation**: Swagger (to be added)

---

## API Endpoints

### Resources

#### List Resources
**GET** `/v1/resources`

**Query Parameters:**
- `level`: Filter by educational level.
- `title`: Search in resource titles.
- `description`: Search in resource descriptions.
- `page`: Page number for pagination (default: 1).

**Response:**
- Includes pagination metadata.

#### Upload Resource
**POST** `/v1/resources`

**Multipart Form Data:**
- `title`: Resource title.
- `description`: Resource description.
- `level`: Educational level.
- `file`: Resource file (PDF/DOC/DOCX/TXT/RTF only).

### Users

#### Sign Up
**POST** `/v1/user/signup`

#### Log In
**POST** `/v1/user/login`

#### Log Out
**GET** `/v1/user/logout`

#### Update User
**PUT** `/v1/user/:uid`

#### Delete User
**DELETE** `/v1/user/:uid`

### Jobs

#### List Jobs
**GET** `/v1/jobs`

#### Post Job
**POST** `/v1/jobs`

---

## Setup and Installation

### Prerequisites

- Go 1.23 or higher.
- PostgreSQL.

### Database Setup

Run the following SQL command to create the database:
```sql
CREATE DATABASE cbcexams;
```

### Configuration

Update the `conf/app.conf` file with your database credentials:
```bash
sqlconn = "user=postgres password=your_password dbname=cbcexams sslmode=disable"
```

### Install Dependencies

Run the following command to install dependencies:
```bash
go mod tidy
```

### Run the Application

Start the server with:
```bash
go run main.go
```

---

## Project Structure

```
├── conf/
│   └── app.conf
├── controllers/
│   ├── resource.go
│   ├── user.go
│   └── job.go
├── models/
│   ├── db.go
│   ├── resource.go
│   ├── user.go
│   └── job.go
├── routers/
│   └── router.go
├── uploads/  # Resource files storage
├── main.go
└── README.md
```

---

## Development

### Run Tests

Run tests using:
```bash
go test ./tests
```

### Generate Swagger Documentation

Generate API documentation with:
```bash
bee generate docs
```

---

## API Documentation

To be added.

