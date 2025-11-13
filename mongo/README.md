# MongoDB Setup

This directory contains Docker Compose configuration for running MongoDB locally during development.

## Purpose

The orchestrator service requires MongoDB for persistent storage of tasks and conversation history. This setup provides a quick way to run MongoDB in a container without manual installation.

## Usage

### Start MongoDB

```bash
cd mongo
docker-compose up -d
# or with Podman
podman-compose up -d
```

### Stop MongoDB

```bash
docker-compose down
# or
podman-compose down
```

## Configuration

- **Port**: `27020` (mapped to container port `27017`)
- **Username**: `admin`
- **Password**: `password`
- **Data Persistence**: Stored in Docker volume `mongodb_data`
