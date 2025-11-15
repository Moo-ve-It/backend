# MooveIt Backend

A RESTful API backend for monitoring and managing a smart farm system. The backend provides endpoints to track the state of cows, a robo-dog, and a drone, all equipped with various sensors that return real-time data.

## 🚀 Overview

MooveIt Backend is a Go-based API server designed to monitor and manage farm operations. It tracks:

- **Cows**: Individual cow health, location, and sensor data
- **Robo-Dog**: Autonomous monitoring robot with environmental sensors
- **Drone**: Aerial surveillance with environmental and positioning sensors

The system provides real-time monitoring capabilities through a RESTful API, enabling web applications to display current farm state, animal health, and equipment status.

## ✨ Features

- **Farm State Monitoring**: Get overall farm statistics including total cows, health status, and equipment states
- **Cow Tracking**: Monitor individual cows with detailed health metrics, location tracking, and sensor data
- **Robo-Dog Monitoring**: Track robo-dog status, location, and environmental sensor readings
- **Drone Surveillance**: Monitor drone status, altitude, location, and environmental conditions
- **Health Check Endpoint**: Server health and status monitoring
- **Metrics Endpoint**: Application metrics and debugging information
- **Structured JSON Logging**: Comprehensive logging with structured JSON output
- **Error Handling**: Robust error handling with proper HTTP status codes
- **Panic Recovery**: Automatic panic recovery middleware
- **Request Logging**: All requests are logged with method and URL

## 🏗️ Architecture

The backend follows a clean, modular architecture:

- **HTTP Router**: Uses `httprouter` for efficient routing
- **Middleware Chain**: Request logging and panic recovery
- **JSON Responses**: Consistent JSON response format with envelope pattern
- **Mock Data**: Currently uses in-memory mock data (ready for database integration)
- **Structured Logging**: JSON-formatted logs with severity levels
- **Version Control**: Automatic version tracking from VCS

## 📡 API Endpoints

### Farm Monitoring

#### Get Farm State
```http
GET /api/farm/state
```

Returns the overall state of the farm including:
- Total number of cows
- Healthy vs sick cow counts
- Robo-dog status
- Drone status
- Last update timestamp

**Response:**
```json
{
  "farm_state": {
    "total_cows": 5,
    "healthy_cows": 4,
    "sick_cows": 1,
    "robodog_status": "active",
    "drone_status": "flying",
    "last_updated": "2024-01-15T10:30:00Z"
  }
}
```

#### List All Cows
```http
GET /api/cows
```

Returns a list of all cows with their complete sensor data.

**Response:**
```json
{
  "cows": [
    {
      "id": 1,
      "name": "Bessie",
      "tag": "COW-001",
      "location": {
        "latitude": 40.7128,
        "longitude": -74.0060,
        "zone": "Pasture A"
      },
      "health": {
        "status": "healthy",
        "temperature": 38.5,
        "heart_rate": 65,
        "activity": "grazing"
      },
      "sensors": {
        "temperature": 38.5,
        "heart_rate": 65,
        "activity": "grazing",
        "battery_level": 85
      },
      "last_updated": "2024-01-15T10:30:00Z"
    }
  ],
  "total": 5
}
```

#### Get Specific Cow
```http
GET /api/cows/:id
```

Returns detailed information for a specific cow by ID.

**Response:**
```json
{
  "cow": {
    "id": 1,
    "name": "Bessie",
    "tag": "COW-001",
    ...
  }
}
```

#### Get Robo-Dog Status
```http
GET /api/robodog
```

Returns the current state and sensor data of the robo-dog.

**Response:**
```json
{
  "robodog": {
    "id": 1,
    "name": "Rex",
    "status": "active",
    "location": {
      "latitude": 40.7129,
      "longitude": -74.0061,
      "zone": "Central Area"
    },
    "sensors": {
      "temperature": 22.5,
      "humidity": 65.0,
      "motion_detected": true,
      "camera_status": "active",
      "audio_level": 45.2
    },
    "battery_level": 72,
    "last_updated": "2024-01-15T10:30:00Z"
  }
}
```

#### Get Drone Status
```http
GET /api/drone
```

Returns the current state and sensor data of the drone.

**Response:**
```json
{
  "drone": {
    "id": 1,
    "name": "SkyEye",
    "status": "flying",
    "location": {
      "latitude": 40.7132,
      "longitude": -74.0059,
      "zone": "Airspace"
    },
    "altitude": 150.0,
    "sensors": {
      "temperature": 18.3,
      "humidity": 58.0,
      "wind_speed": 12.5,
      "camera_status": "active",
      "gps_accuracy": 2.5,
      "air_quality": 45.0
    },
    "battery_level": 68,
    "last_updated": "2024-01-15T10:30:00Z"
  }
}
```

### System Endpoints

#### Health Check
```http
GET /api/healthcheck
```

Returns server health status and system information.

**Response:**
```json
{
  "status": "available",
  "system_info": {
    "environment": "development",
    "version": "2024-01-15T10:00:00Z-abc123"
  }
}
```

#### Metrics
```http
GET /api/debug/vars
```

Returns application metrics including:
- Version information
- Active goroutines count
- Current timestamp

## 🛠️ Technology Stack

- **Language**: Go 1.21.6
- **HTTP Router**: [httprouter](https://github.com/julienschmidt/httprouter)
- **Logging**: Custom JSON logger
- **Deployment**: Railway (configured)

## 📦 Project Structure

```
mooveit-backend/
├── cmd/
│   └── api/
│       ├── main.go              # Application entry point
│       ├── routes.go            # Route definitions and middleware
│       ├── helpers.go           # HTTP helper functions
│       ├── healthcheck.go       # Health check handler
│       └── farm_handlers.go     # Farm monitoring handlers
├── internal/
│   ├── jsonlog/                 # Structured JSON logging
│   │   └── log.go
│   ├── validator/               # Input validation utilities
│   │   └── validator.go
│   └── vcs/                     # Version control system utilities
│       └── vcs.go
├── migrations/                  # Database migrations (future)
├── bin/                         # Compiled binaries
├── go.mod                       # Go module dependencies
├── go.sum                       # Go module checksums
├── Makefile                     # Build automation
├── railway.json                 # Railway deployment configuration
└── README.md                    # This file
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21.6 or later
- Git

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd mooveit-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o bin/api ./cmd/api
```

Or use the Makefile:
```bash
make build
```

### Running the Server

#### Development Mode

Run directly with Go:
```bash
go run ./cmd/api
```

Or run the compiled binary:
```bash
./bin/api
```

The server will start on port `4000` by default (or the port specified in the `PORT` environment variable).

#### With Custom Port

```bash
PORT=8080 go run ./cmd/api
```

Or:
```bash
./bin/api -port=8080
```

### Configuration

The application supports configuration through command-line flags and environment variables:

- **Port**: `-port` flag or `PORT` environment variable (default: 4000)
- **Environment**: `-env` flag or `ENV` environment variable (default: development)
- **Version**: Display version with `-version` flag

**Environment Variables:**
- `PORT`: Server port number
- `ENV`: Environment (development|staging|production)
- `RAILWAY_PUBLIC_DOMAIN`: Railway public domain (auto-set by Railway)
- `RAILWAY_STATIC_URL`: Railway service URL (auto-set by Railway)
- `PUBLIC_DOMAIN`: Custom public domain

## 🔧 Development

### Building

```bash
go build -o bin/api ./cmd/api
```

### Testing Endpoints

Once the server is running, you can test the endpoints:

```bash
# Health check
curl http://localhost:4000/api/healthcheck

# Get farm state
curl http://localhost:4000/api/farm/state

# List all cows
curl http://localhost:4000/api/cows

# Get specific cow
curl http://localhost:4000/api/cows/1

# Get robo-dog status
curl http://localhost:4000/api/robodog

# Get drone status
curl http://localhost:4000/api/drone

# Get metrics
curl http://localhost:4000/api/debug/vars
```

### Code Structure

- **Handlers**: Located in `cmd/api/farm_handlers.go` - contain business logic for farm monitoring
- **Routes**: Defined in `cmd/api/routes.go` - maps URLs to handlers
- **Helpers**: Utility functions in `cmd/api/helpers.go` - JSON responses, error handling
- **Logging**: Custom JSON logger in `internal/jsonlog/` - structured logging with severity levels
- **Validation**: Input validation utilities in `internal/validator/`

## 📊 Data Models

### Cow
- ID, Name, Tag
- Location (GPS coordinates, zone)
- Health status (healthy/sick/injured)
- Health metrics (temperature, heart rate, activity)
- Sensor data (temperature, heart rate, activity, battery level)

### Robo-Dog
- ID, Name, Status
- Location
- Sensors (temperature, humidity, motion detection, camera, audio)
- Battery level

### Drone
- ID, Name, Status
- Location, Altitude
- Sensors (temperature, humidity, wind speed, camera, GPS accuracy, air quality)
- Battery level

## 🚢 Deployment

### Railway Deployment

The project is configured for Railway deployment with `railway.json`:

```json
{
  "build": {
    "builder": "NIXPACKS",
    "buildCommand": "go build -o bin/api ./cmd/api"
  },
  "deploy": {
    "startCommand": "./bin/api",
    "restartPolicyType": "ON_FAILURE",
    "restartPolicyMaxRetries": 10
  }
}
```

Railway automatically:
- Detects the Go project
- Builds the application
- Sets the `PORT` environment variable
- Provides `RAILWAY_PUBLIC_DOMAIN` for public access

## 📝 Logging

The application uses structured JSON logging with the following levels:

- **INFO**: General information messages
- **ERROR**: Error messages with properties
- **FATAL**: Fatal errors that terminate the application

Logs include:
- Timestamp (PST timezone)
- Severity level
- Message
- Properties (key-value pairs)
- Stack trace (for ERROR and FATAL levels)

Example log entry:
```json
{
  "level": "INFO",
  "time": "15-Jan-24 10:30:00.123 PST",
  "message": "💭 Server starting",
  "properties": {
    "port": "4000",
    "environment": "development"
  }
}
```

## 🔒 Error Handling

The API returns consistent error responses:

- **404 Not Found**: Resource not found
- **500 Internal Server Error**: Server errors

Error response format:
```json
{
  "error": "Error message here"
}
```

## 🧪 Future Enhancements

- Database integration (PostgreSQL/MySQL)
- Authentication and authorization
- Real-time sensor data updates via WebSocket
- Historical data tracking
- Alert system for health issues
- Database migrations
- Unit and integration tests
- API documentation with Swagger/OpenAPI

## 📄 License

[Add your license here]

## 👥 Contributing

[Add contribution guidelines here]

## 📞 Support

[Add support information here]
