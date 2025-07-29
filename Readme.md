# Matiq App

A modern workflow automation platform that allows users to design, build, and manage automated workflows through a visual interface. Similar to n8n, this platform enables seamless integration between different services and APIs with robust authentication and authorization.

## Features

- Visual workflow designer with drag-and-drop interface
- Keycloak-based authentication and group-based permissions
- RESTful API with Bearer token security
- Support for various node types (databases, APIs, triggers, actions)
- Real-time workflow execution and monitoring

## Tech Stack

**Backend:**
- Go (Gin framework)
- PostgreSQL
- Keycloak (Authentication & Authorization)
- UUID-based entity management

**Frontend:**
- Next.js 14 (React)
- TypeScript
- NextAuth.js
- Tailwind CSS

## Quick Start

1. Set up Keycloak and configure groups (`workflow-viewers`, `workflow-managers`)
2. Configure environment variables for database and Keycloak
3. Run the Go backend: `go run cmd/workflow-builder-backend/main.go`
4. Start the Next.js frontend: `npm run dev`
5. Access the application at `http://localhost:3000`