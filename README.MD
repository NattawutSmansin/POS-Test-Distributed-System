A Go-based microservice architecture with **Core-Service** and **Branch-Service**, powered by MySQL and RabbitMQ.

---

## Environment Setup

Before running the application, rename the environment file:

```bash
mv env_dev .env
```

> Make sure `.env` is configured with your correct database and RabbitMQ credentials before proceeding.

---

## Getting Started

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Generate Swagger Docs

```bash
swag init
```

### 3. Run Services (MySQL + RabbitMQ)

```bash
docker-compose up -d
```

### 4. Check Container Status

```bash
docker ps
```

### 5. Run Application

```bash
go run .
```

---

## Swagger UI

After the application is running, access the API documentation via the following URLs:

| Service          | Swagger URL                                                  |
|------------------|--------------------------------------------------------------|
| **Core-Service**  | `http://localhost:3344/api/core/docs/swagger/index.html`     |
| **Branch-Service** | `http://localhost:5566/api/branch/docs/swagger/index.html` |

---

## Architecture Overview

| Component        | Description                              |
|------------------|------------------------------------------|
| **Core-Service**  | Core service running on port `3344`      |
| **Branch-Service** | Branch service running on port `5566`  |
| **MySQL**        | Relational database (via Docker)         |
| **RabbitMQ**     | Message broker (via Docker)              |

---

## Install `swag` if not installed
```bash
- go install github.com/swaggo/swag/cmd/swag@latest