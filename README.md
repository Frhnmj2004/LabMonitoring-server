# Lab Resource Monitoring System

A real-time computer lab resource monitoring system built with Go and Fiber framework. This system monitors CPU, memory, and network usage of lab computers and provides a real-time admin dashboard with WebSocket support.

## Features

- Real-time resource monitoring (CPU, Memory, Network)
- WebSocket-based live updates
- JWT Authentication with role-based access control
- Offline data buffering with sync mechanism
- Alert system for high resource usage
- Historical data viewing
- PostgreSQL database with GORM ORM

## Prerequisites

- Go 1.19 or higher
- PostgreSQL 12 or higher
- Environment variables set in `.env` file

## Setup

1. Clone the repository:
```bash
git clone https://github.com/Frhnmj2004/LabMonitoring-server.git
cd LabMonitoring-server
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up your PostgreSQL database and update the `.env` file with your configuration:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=lab_monitor
JWT_SECRET=your-super-secret-jwt-key-change-in-production
PORT=8080
```

4. Run the server:
```bash
go run main.go
```

## API Endpoints

### Authentication
- `POST /api/v1/login`: User login
- `POST /api/v1/signup`: Create new user (Admin only)

### Resource Monitoring
- `POST /api/v1/resource`: Submit resource data
- `GET /api/v1/history`: Get resource history (Admin only)
- `GET /api/v1/alerts`: Get resource alerts (Admin only)

### WebSocket
- `WS /ws/resources`: WebSocket endpoint for real-time updates

## Project Structure

```
lab-monitor/
├── main.go                # Entry point
├── config/
│   └── db.go             # Database configuration
├── controllers/
│   ├── authController.go # Authentication handlers
│   └── resourceController.go # Resource monitoring handlers
├── middleware/
│   └── authMiddleware.go # JWT authentication middleware
├── models/
│   ├── user.go          # User model
│   ├── computer.go      # Computer model
│   └── resourceLog.go   # Resource log model
├── routes/
│   └── routes.go        # API routes
├── storage/
│   └── buffer.go        # Offline data buffer
├── utils/
│   ├── jwt.go           # JWT utilities
│   └── logger.go        # Logging utility
└── websocket/
    └── handlers.go      # WebSocket handlers
```

## Security

- JWT-based authentication
- Role-based access control
- Password hashing with bcrypt
- Input validation
- CORS protection
- Panic recovery middleware

## Error Handling

The system includes comprehensive error handling:
- Database connection errors
- Authentication failures
- Invalid input validation
- WebSocket connection issues
- Resource monitoring errors

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

MIT License
