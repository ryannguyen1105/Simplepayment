# Simplepayment

A backend payment system built with Golang. The project focuses on processing transactions securely, authenticating users, and managing database state.

## Tech Stack

* Core Backend: Golang, Gin Framework
* Database: PostgreSQL, SQLC (Type-safe SQL generation)
* Security & Auth: PASETO, JWT, Bcrypt (Password hashing)
* Infrastructure & Cloud: Docker, Docker Compose, Kubernetes, Render, Neon DB

## Features Implemented

### 1. Atomic Transaction Management
* Developed payment processing logic using PostgreSQL transactions to ensure ACID compliance and prevent data inconsistency during balance transfers.

### 2. Type-Safe Database Layer
* Implemented database queries using SQLC to ensure type safety and reduce runtime errors.

### 3. Secure Authentication & Authorization
* Implemented token-based authentication using JWT and PASETO, including token creation and verification.
* Encrypted user passwords using Bcrypt hashing before saving to the database.

### 4. Multi-Protocol API (REST & gRPC)
* Designed RESTful APIs using Gin with request validation, and implemented inter-service communication using gRPC and Protocol Buffers.

### 5. Containerized Infrastructure
* Containerized the application using Docker and Docker Compose for consistent development and deployment environments.

## Cloud Deployment

The application is deployed and hosted on cloud platforms:

* Deployment Endpoint: Render (Public Domain)
* Database Provider: Neon DB (Managed PostgreSQL)
