# API-Golang-Books

A RESTful API for managing books, users, and related data (such as villagers, image uploads), built with **Go (Golang)**, **Gin Web Framework**, **GORM**, and **Swagger** for API documentation. This project demonstrates a clean architecture with JWT authentication, file uploads (MinIO), and database migrations.

## Features

- **User Authentication** – Register, login, and JWT-based authorization.
- **Book Management** – CRUD operations for books (admin/user role-based).
- **Image Upload** – Upload book covers or user avatars via **MinIO** (or local fallback).
- **Villager Management** – Additional entity demonstration (e.g., for game or reference data).
- **API Documentation** – Auto-generated Swagger UI (via `swaggo`).
- **CORS Support** – Configured for cross-origin requests.
- **Environment Configuration** – Uses `.env` for secrets and settings.
- **CI/CD Ready** – GitHub Actions workflow included.

## Tech Stack

| Category      | Technology                                                    |
| ------------- | ------------------------------------------------------------- |
| Language      | Go 1.19+                                                      |
| Web Framework | Gin                                                           |
| ORM           | GORM                                                          |
| Database      | PostgreSQL (or SQLite for local dev – adjust `dsn` in `.env`) |
| Auth          | JWT (golang-jwt)                                              |
| File Storage  | MinIO (S3-compatible) – fallback to local disk                |
| Docs          | Swagger (swaggo/swag)                                         |
| Logging       | Logrus                                                        |
| Other Tools   | godotenv, bcrypt, cors, validator                             |

## Project Structure

API-Golang-Books/
├── controllers/ # All handler functions
├── models/ # Database models (User, Book/Post, Villager, Loan, etc.)
├── routers/ # Route definitions
├── middleware/ # JWT auth middleware
├── initializer/ # DB connection, env loading, migration
├── utils/ # Helper functions
├── migrate/ # Migration scripts (if any)
├── docs/ # Swagger documentation
├── .env.example
├── main.go
└── go.mod

## Quick Start

### 1. Clone the repositorybash

git clone https://github.com/ferytell/API-Golang-Books.git
cd API-Golang-Books

### 2. Copy environment filebash

cp .env.example .env

Edit .env with your credentials:env

DB_URL=postgres://user:password@host:port/dbname?sslmode=disable
MINIO_ENDPOINT="your-minio-endpoint"
MINIO_ACCESS_KEY="your-access-key"
MINIO_SECRET_KEY="your-secret-key"
MINIO_BUCKET="your-bucket-name"

### 3. Install dependenciesbash

go mod tidy

### 4. Run the serverbash

go run main.go

Or build and run:bash

go build -o api-books
./api-books

The server will start on http://localhost:8080 (or the PORT env variable).
