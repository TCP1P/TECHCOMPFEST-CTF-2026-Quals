# Go-backend for super note

A  API for Notes management built with Go. This project serves as a backend of the application, demonstrating industry best practices for structuring a Go backend application.


## Features

- **RESTful API Design**: Clean and consistent API endpoints following REST principles
- **PostgreSQL Database**: Robust data persistence with GORM ORM
- **JWT Authentication**: Secure user authentication and authorization
- **Structured Logging**: Comprehensive logging with Zap logger
- **API Documentation**: Auto-generated Swagger documentation
- **Environment Configuration**: Flexible configuration via environment variables
- **Graceful Shutdown**: Proper handling of server shutdown
- **Middleware Support**: Extensible middleware architecture
- **DTO Pattern**: Clean separation of data transfer objects
- **Repository Pattern**: Separation of data access logic
- **Service Layer**: Business logic encapsulation

## Project Structure

```text
├── api/                  # API layer
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   └── routes/           # Route definitions
├── config/               # Configuration management
├── database/             # Database connection and migrations
├── docs/                 # Swagger documentation
├── internal/             # Internal application code
│   ├── dtos/             # Data Transfer Objects
│   ├── logger/           # Logger configuration
│   ├── models/           # Database models
│   ├── repositories/     # Data access layer
│   ├── services/         # Business logic layer
│   └── utils/            # Utility functions
├── .env                  # Environment variables
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── main.go               # Application entry point
└── dev.sh               # Development script
```

## Prerequisites

- Go 1.18 or higher
- PostgreSQL 12 or higher
- Git

## License

MIT

