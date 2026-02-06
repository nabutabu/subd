# subd

A daemon that runs on all hosts provisioned by Crane-OSS to maintain system state and communicate with the Dominator service.

## Overview

subd is a lightweight agent that continuously monitors host state and synchronizes with a central Dominator service. It operates on a simple observe-report-reconcile loop to ensure hosts maintain their desired configuration.

## Architecture

- **Collector**: Gathers current system state (running services, filesystem, processes)
- **Client**: Communicates with the Dominator API via HTTP
- **Agent**: Orchestrates the observation and reconciliation loop

## Configuration

Create `appsettings.json` with:

```json
{
  "DominatorURL": "https://your-dominator.example.com",
  "Token": "your-auth-token"
}
```

## Usage

```bash
go run cmd/main.go
```

The agent will:
1. Collect current system state every 2 seconds
2. Send heartbeat with current state to Dominator
3. Receive desired state response
4. Compute diff and reconcile if needed

## Development

Build:
```bash
go build -o subd cmd/main.go
```
