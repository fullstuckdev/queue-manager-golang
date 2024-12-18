# Queue Manager Service

A multi-tenant queue management system built with Go, featuring RabbitMQ integration and a CLI tool.

## Features

- ✨ Create and manage tenants with dedicated RabbitMQ queues
- 🔄 Automatic consumer setup for each tenant
- 🗑️ Soft delete tenants with proper resource cleanup
- 📤 Process JSON payloads for specific tenants
- 🔁 Auto-reconnect for RabbitMQ connections
- 🛡️ Panic recovery in consumers
- 🛠️ CLI tool for tenant management

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or higher (for local development)

## Quick Start

1. Clone the repository:
```bash
git clone <repository-url>
cd queue-manager
```

2. Start the services using Docker:
```bash
docker-compose up --build
```

The services will be available at:
- API Server: http://localhost:8080
- RabbitMQ Management: http://localhost:15672 (guest/guest)
- PostgreSQL: localhost:5432

## API Endpoints

### Create Tenant
```bash
POST /tenants

Request:
{
    "client_id": "tenant123",
    "name": "Example Tenant"
}

Response: (201 Created)
{
    "id": 1,
    "client_id": "tenant123",
    "name": "Example Tenant",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
}
```

### Delete Tenant
```bash
DELETE /tenants/{clientID}

Response: 204 No Content
```

### Process Payload
```bash
POST /tenants/process

Request:
{
    "client_id": "tenant123",
    "payload": {
        "key": "value",
        "message": "Hello, World!"
    }
}

Response: (200 OK)
{
    "status": "Message published successfully"
}
```
