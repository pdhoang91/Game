# Oden - Idle RPG Game MVP

An Idle RPG mobile game inspired by "Ode To Heroes" developed for iOS using Unity for the client, Golang for the backend, and AWS for server deployment.

## Project Structure

```
Oden/
├── client/              # Unity client project
│   ├── Assets/          # Unity assets
│   │   ├── Scripts/     # C# code for game logic
│   │   ├── Prefabs/     # Reusable game objects
│   │   ├── Scenes/      # Game scenes
│   │   └── ...
├── server/              # Golang backend
│   ├── cmd/             # Application entrypoints
│   ├── internal/        # Private application code
│   │   ├── api/         # API handlers
│   │   ├── auth/        # Authentication logic
│   │   ├── config/      # Configuration
│   │   ├── db/          # Database interactions
│   │   ├── game/        # Game logic
│   │   └── model/       # Data models
│   ├── pkg/             # Public libraries
│   └── ...
└── docs/                # Documentation
    ├── architecture.md  # System architecture
    ├── api.md           # API documentation
    └── deployment.md    # Deployment guide
```

## System Architecture

![Architecture Diagram](docs/images/architecture_diagram.png)

### Client-Server Flow
1. Player launches the game on iOS device
2. Authentication via Unity client to Golang backend
3. Game data (heroes, resources) fetched from the backend
4. Battles calculated on the server, results sent to client
5. Client displays battle animations and rewards
6. Idle rewards calculated on server based on offline time

### AWS Services Used
- EC2: Hosting the Golang API server
- RDS (MySQL): Game database storage
- S3: Assets storage (hero images, etc.)
- CloudFront: Content delivery for static assets
- Route 53: DNS management
- CloudWatch: Monitoring and logging

## Docker Local Development Setup

### Prerequisites

- Docker and Docker Compose
- Git

### Setup Instructions

1. Clone the repository:
   ```bash
   git clone https://your-repository-url/oden.git
   cd oden
   ```

2. Setup configuration:
   ```bash
   cd server/internal/config
   cp config.example.json config.json
   cd ../../..
   ```

3. Start the services:
   ```bash
   docker-compose up -d
   ```

4. Access the services:
   - API Server: http://localhost:8080
   - MinIO Console (for S3-compatible storage): http://localhost:9001
     - Username: minioadmin
     - Password: minioadmin

5. Check if the services are running:
   ```bash
   docker-compose ps
   ```

6. Stop the services:
   ```bash
   docker-compose down
   ```

### Important Notes

- The MySQL database is exposed on port 3306
- Initial database migrations are automatically applied from `server/internal/db/migrations`
- The API server is rebuilt when you change the source code
- MinIO provides S3-compatible storage for game assets
- Configuration can be adjusted in `server/internal/config/config.json`

## Getting Started

See the following documentation:
- [Client Setup Guide](docs/client_setup.md)
- [Server Setup Guide](docs/server_setup.md)
- [Deployment Guide](docs/deployment.md)
- [Testing Guide](docs/testing.md) 