# Taxi Booking Backend - Golang

A modern, scalable taxi booking backend built with Go, Gin, GORM, and PostgreSQL. Supports real-time ride tracking via WebSocket, driver allocation, and trip management.

## Architecture Overview

The backend follows a clean, layered architecture:

```
├── config/          # Database and environment configuration
├── models/          # Data models (User, Driver, Ride, Vehicle, Location)
├── repositories/    # Data access layer (abstraction over GORM)
├── services/        # Business logic (ride management, driver allocation)
├── controllers/     # HTTP handlers (Gin endpoints)
├── routes/          # API route definitions
├── realtime/        # WebSocket hub for live updates
├── main.go          # Application entry point
└── docker-compose.yml  # Local development environment
```

## Core Features

### 1. **Ride Management**
- Request new rides with pickup/dropoff locations
- Real-time ride status tracking (requested → assigned → started → completed)
- Automatic fare estimation based on distance
- Ride history for passengers and drivers
- Cancel rides with reasons

### 2. **Driver Management**
- Driver registration and verification
- Online/offline status management
- Location tracking
- Driver availability and rating system
- Ride history and statistics

### 3. **Intelligent Driver Allocation**
- Find nearest available drivers using Haversine formula
- Allocate drivers based on proximity, ratings, and availability
- Status transitions with validation
- Real-time status updates

### 4. **Real-Time Updates**
- WebSocket connection for live ride updates
- Broadcast ride status changes to dashboard clients
- Scalable hub-based architecture with goroutines and channels

### 5. **Data Models**
- **User** - Passengers with profile and ride history
- **Driver** - Drivers with verification, location, and statistics
- **Ride** - Trip records with status, pricing, and locations
- **Vehicle** - Vehicle information linked to drivers
- **Location** - Saved addresses for quick booking

## API Endpoints

### Health Check
- `GET /api/v1/health` - Server health status

### Rides
- `POST /api/v1/rides` - Request a new ride
- `GET /api/v1/rides/:id` - Get ride details
- `GET /api/v1/rides/passenger/:id` - Get passenger's ride history
- `GET /api/v1/rides/pending` - Get all pending rides
- `PUT /api/v1/rides/:id/cancel` - Cancel a ride
- `PUT /api/v1/rides/:id/start` - Start a ride
- `PUT /api/v1/rides/:id/complete` - Complete a ride
- `POST /api/v1/rides/:id/allocate` - Allocate driver to ride
- `POST /api/v1/rides/drivers/nearest` - Find nearest drivers

### Drivers
- `GET /api/v1/drivers/:id` - Get driver details
- `GET /api/v1/drivers/available` - Get all available drivers
- `GET /api/v1/drivers/online` - Get all online drivers
- `GET /api/v1/drivers/:id/rides` - Get driver's ride history
- `PUT /api/v1/drivers/:id/location` - Update driver location
- `PUT /api/v1/drivers/:id/status` - Update driver status

### WebSocket
- `GET /api/v1/ws/ride-updates` - WebSocket connection for real-time updates

## Prerequisites

- **Go** 1.25.0 or higher
- **PostgreSQL** 12 or higher
- **Docker** (optional, for containerized setup)
- **Docker Compose** (optional, for orchestration)

## Local Development Setup

### 1. Clone Repository
```bash
git clone <repository-url>
cd backend
```

### 2. Create `.env` File
Copy `.env.example` to `.env` and configure:
```bash
cp .env.example .env
```

Edit `.env`:
```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/taxi_booking?sslmode=disable
PORT=8080
GIN_MODE=debug
```

### 3. Start PostgreSQL Database

**Option A: Using Docker Compose (Recommended)**
```bash
docker-compose up -d
```

This starts:
- PostgreSQL on `localhost:5432`
- Backend on `localhost:8080`

**Option B: Local PostgreSQL**
```bash
# Create database
createdb taxi_booking

# Update DATABASE_URL in .env with your local credentials
DATABASE_URL=postgres://your_user:your_password@localhost:5432/taxi_booking?sslmode=disable
```

### 4. Install Dependencies
```bash
go mod download
go mod tidy
```

### 5. Run Server
```bash
# Hot reload with air (install first: go install github.com/cosmtrek/air@latest)
air

# Or standard go run
go run main.go
```

Server starts on `http://localhost:8080`

## Development Workflow

### Adding New Features

#### 1. Define Models
Create new model in `models/your_model.go`:
```go
type YourModel struct {
    ID    uint   `gorm:"primaryKey" json:"id"`
    Name  string `json:"name"`
}
```

#### 2. Create Repository
Add methods in `repositories/your_repository.go`:
```go
type YourRepository interface {
    Create(item *YourModel) error
    GetByID(id uint) (*YourModel, error)
}
```

#### 3. Implement Service Logic
Add business logic in `services/your_service.go`:
```go
func (s *yourService) DoSomething() error {
    // Business logic here
}
```

#### 4. Create Controller Handler
Add HTTP handlers in `controllers/your_controller.go`:
```go
func (c *YourController) HandleRequest(ctx *gin.Context) {
    // HTTP handling
}
```

#### 5. Register Routes
Update `routes/routes.go`:
```go
group := api.Group("/your-resource")
{
    group.GET("/:id", controller.Handler)
}
```

### Database Migrations

Migrations run automatically on startup via `config.DB.AutoMigrate()` in `main.go`.

To add new migrations, add to the `migrateDatabase()` function:
```go
config.DB.AutoMigrate(&models.NewModel{})
```

### WebSocket Integration

Broadcast ride updates from services:
```go
// In ride_service.go or controller
update := &realtime.RideUpdate{
    RideID: ride.ID,
    Status: string(ride.Status),
}
wsHub.BroadcastUpdate(update)
```

## Testing

### Manual API Testing

Using `curl`:
```bash
# Request a ride
curl -X POST http://localhost:8080/api/v1/rides \
  -H "Content-Type: application/json" \
  -d '{
    "passenger_id": 1,
    "pickup_latitude": 40.7128,
    "pickup_longitude": -74.0060,
    "dropoff_latitude": 40.7580,
    "dropoff_longitude": -73.9855,
    "dropoff_address": "Times Square, NY"
  }'

# Get ride details
curl http://localhost:8080/api/v1/rides/1

# Update driver location
curl -X PUT http://localhost:8080/api/v1/drivers/1/location \
  -H "Content-Type: application/json" \
  -d '{"latitude": 40.7128, "longitude": -74.0060}'
```

Using Postman:
1. Import API endpoints from comments
2. Set `Content-Type: application/json`
3. Test request/response payloads

### WebSocket Testing

Using `websocat` or similar:
```bash
# Install websocat: cargo install websocat
websocat ws://localhost:8080/api/v1/ws/ride-updates
```

## Performance Considerations

1. **Database Indexes** - GORM creates indexes for primaryKey, unique, and explicit `index` tags
2. **Connection Pooling** - PostgreSQL driver handles connection pooling automatically
3. **WebSocket Buffering** - Set appropriate buffer sizes in `realtime/handler.go`
4. **Goroutines** - Each WebSocket client runs 2 goroutines (read/write)

## Deployment

### Build Docker Image
```bash
docker build -t taxi-backend .
```

### Run with Docker Compose
```bash
docker-compose up --build
```

### Production Considerations

Before deploying to production:

1. **Update `.env` with production values**
   ```env
   DATABASE_URL=postgres://prod_user:prod_password@prod_host:5432/taxi_booking?sslmode=require
   GIN_MODE=release
   PORT=8080
   ```

2. **Enable HTTPS** - Use Gin's `RunTLS()` or reverse proxy (Nginx)

3. **Add Authentication** - Implement JWT middleware (future)

4. **Security**
   - Set proper CORS in Gin middleware
   - Validate all inputs
   - Use environment variables for secrets
   - Implement rate limiting

5. **Monitoring**
   - Add logging middleware
   - Export metrics (Prometheus)
   - Monitor database performance

6. **Scaling**
   - Use connection pooling
   - Implement caching (Redis)
   - Consider message queues for async operations
   - Load balance WebSocket connections

## Project Structure Details

### Config
- `database.go` - PostgreSQL connection and initialization

### Models
- `user.go` - Passenger/customer model
- `driver.go` - Driver model with verification
- `vehicle.go` - Vehicle information
- `ride.go` - Trip/ride records
- `location.go` - Saved user locations

### Repositories
- `ride_repository.go` - Ride data access
- `driver_repository.go` - Driver data access
- `user_repository.go` - User data access

### Services
- `ride_service.go` - Ride business logic
- `driver_service.go` - Driver business logic

### Controllers
- `ride_controller.go` - Ride HTTP handlers
- `driver_controller.go` - Driver HTTP handlers
- `health_controller.go` - Health check handler

### Realtime
- `hub.go` - WebSocket hub and client management
- `handler.go` - WebSocket HTTP upgrade handler

## Common Issues

### Database Connection Failed
```
Failed to connect to database: dial tcp 127.0.0.1:5432: connection refused
```
**Solution**: Ensure PostgreSQL is running and DATABASE_URL is correct

### Port Already in Use
```
listen tcp :8080: bind: address already in use
```
**Solution**: Change PORT in `.env` or kill process using port 8080

### Missing Dependencies
```
no required module provides package github.com/gorilla/websocket
```
**Solution**: Run `go mod tidy` and `go mod download`

## Future Enhancements

- [ ] User authentication (JWT)
- [ ] Payment integration
- [ ] Email/SMS notifications
- [ ] Driver rating and review system
- [ ] Surge pricing algorithm
- [ ] Ride scheduling
- [ ] Multi-language support
- [ ] Analytics dashboard
- [ ] Admin panel
- [ ] Mobile app backend optimizations

## Contributing

1. Create a feature branch
2. Make changes following the architecture patterns
3. Test thoroughly
4. Submit pull request

## License

MIT License

## Support

For issues or questions:
1. Check existing issues
2. Review API documentation
3. Test with provided examples
4. Create new issue with details

---

**Built with Go, Gin, GORM, PostgreSQL, and WebSockets**
