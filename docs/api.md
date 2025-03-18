# Lab Monitoring Server API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
The API uses JWT (JSON Web Token) for authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## Endpoints

### Authentication

#### 1. Login
```http
POST /login
```
**Request Body:**
```json
{
    "email": "string",
    "password": "string"
}
```
**Response:**
```json
{
    "token": "string",
    "user": {
        "id": "uuid",
        "email": "string",
        "role": "string"
    }
}
```

#### 2. Signup (Admin only)
```http
POST /signup
```
**Request Body:**
```json
{
    "email": "string",
    "password": "string",
    "role": "string"  // "admin" or "user"
}
```
**Response:**
```json
{
    "id": "uuid",
    "email": "string",
    "role": "string"
}
```

### Resource Monitoring

#### 1. Submit Resource Data
```http
POST /resource
```
**Request Body:**
```json
{
    "computer_id": "uuid",
    "cpu": "float",        // CPU usage percentage (0-100)
    "memory": "float",     // Memory usage percentage (0-100)
    "network_in": "float", // Network incoming traffic (bytes/sec)
    "network_out": "float" // Network outgoing traffic (bytes/sec)
}
```
**Response:**
```json
{
    "id": "uuid",
    "computer_id": "uuid",
    "cpu": "float",
    "memory": "float",
    "network_in": "float",
    "network_out": "float",
    "timestamp": "string"
}
```

#### 2. Get Resource History (Admin only)
```http
GET /resources/history
```
**Query Parameters:**
- `computer_id` (optional): Filter by computer ID
- `start_time` (optional): Start time for filtering (ISO 8601)
- `end_time` (optional): End time for filtering (ISO 8601)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 50)

**Response:**
```json
{
    "data": [
        {
            "id": "uuid",
            "computer_id": "uuid",
            "cpu": "float",
            "memory": "float",
            "network_in": "float",
            "network_out": "float",
            "timestamp": "string"
        }
    ],
    "pagination": {
        "current_page": "integer",
        "total_pages": "integer",
        "total_items": "integer",
        "per_page": "integer"
    }
}
```

### Alerts

#### 1. Get All Alerts (Admin only)
```http
GET /alerts
```
**Query Parameters:**
- `computer_id` (optional): Filter by computer ID
- `type` (optional): Filter by alert type ("HIGH_CPU" or "HIGH_MEMORY")
- `resolved` (optional): Filter by resolution status (true/false)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 50)

**Response:**
```json
{
    "data": [
        {
            "id": "uuid",
            "computer_id": "uuid",
            "type": "string",
            "message": "string",
            "timestamp": "string",
            "resolved": "boolean"
        }
    ],
    "pagination": {
        "current_page": "integer",
        "total_pages": "integer",
        "total_items": "integer",
        "per_page": "integer"
    }
}
```

#### 2. Get Active Alerts (Admin only)
```http
GET /alerts/active
```
**Query Parameters:**
- `computer_id` (optional): Filter by computer ID

**Response:**
```json
{
    "data": [
        {
            "id": "uuid",
            "computer_id": "uuid",
            "type": "string",
            "message": "string",
            "timestamp": "string",
            "resolved": false
        }
    ]
}
```

#### 3. Get Alert History (Admin only)
```http
GET /alerts/history
```
**Query Parameters:**
- `computer_id` (optional): Filter by computer ID
- `type` (optional): Filter by alert type
- `start_time` (optional): Start time for filtering (ISO 8601)
- `end_time` (optional): End time for filtering (ISO 8601)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 50)

**Response:**
```json
{
    "data": [
        {
            "id": "uuid",
            "computer_id": "uuid",
            "type": "string",
            "message": "string",
            "timestamp": "string",
            "resolved": "boolean"
        }
    ],
    "pagination": {
        "current_page": "integer",
        "total_pages": "integer",
        "total_items": "integer",
        "per_page": "integer"
    }
}
```

#### 4. Get Alert Statistics (Admin only)
```http
GET /alerts/stats
```
**Response:**
```json
{
    "total_stats": {
        "total_alerts": "integer",
        "active_alerts": "integer",
        "resolved_alerts": "integer",
        "high_cpu_alerts": "integer",
        "high_memory_alerts": "integer"
    },
    "last_24h": {
        "total_alerts": "integer"
    }
}
```

#### 5. Resolve Alert (Admin only)
```http
PUT /alerts/:id/resolve
```
**Response:**
```json
{
    "message": "Alert resolved successfully",
    "data": {
        "id": "uuid",
        "computer_id": "uuid",
        "type": "string",
        "message": "string",
        "timestamp": "string",
        "resolved": true
    }
}
```

## WebSocket Connection

### Resource Updates WebSocket
```
ws://localhost:8080/ws/resources
```

The WebSocket connection provides real-time updates for:
1. Resource data updates
2. New alerts
3. Alert resolutions

**Message Types:**

1. Resource Update:
```json
{
    "type": "resource_update",
    "data": {
        "computer_id": "uuid",
        "cpu": "float",
        "memory": "float",
        "network_in": "float",
        "network_out": "float",
        "timestamp": "string"
    }
}
```

2. New Alert:
```json
{
    "type": "alert",
    "data": {
        "id": "uuid",
        "computer_id": "uuid",
        "type": "string",
        "message": "string",
        "timestamp": "string",
        "resolved": false
    }
}
```

3. Alert Resolution:
```json
{
    "type": "alert_resolved",
    "data": {
        "id": "uuid",
        "computer_id": "uuid",
        "type": "string",
        "message": "string",
        "timestamp": "string",
        "resolved": true
    }
}
```

## Error Responses

The API uses standard HTTP status codes and returns error messages in the following format:

```json
{
    "error": "string"
}
```

Common status codes:
- 200: Success
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error
- 503: Service Unavailable (when database is down and data is buffered)

## Rate Limiting
Currently, there are no rate limits implemented, but it's recommended to implement client-side throttling for resource submissions (e.g., once every 5-10 seconds per computer).

## Notes for Flutter Developers

1. **WebSocket Management**:
   - Implement reconnection logic with exponential backoff
   - Handle connection drops gracefully
   - Consider using packages like `web_socket_channel`

2. **Authentication**:
   - Store JWT token securely
   - Implement token refresh mechanism
   - Add token to all HTTP requests via interceptor

3. **State Management**:
   - Consider using state management solutions (e.g., Provider, Bloc, Riverpod) for handling real-time updates
   - Implement local caching for offline support

4. **Error Handling**:
   - Implement global error handling
   - Show appropriate error messages to users
   - Handle network connectivity issues gracefully

5. **Data Models**:
   - Create Dart models for all API responses
   - Use JSON serialization for type safety
   - Implement proper null safety
