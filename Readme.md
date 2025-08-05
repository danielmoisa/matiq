# Matiq App

A modern flow automation platform that allows users to design, build, and manage automated flows through a visual interface. Similar to n8n, this platform enables seamless integration between different services and APIs with robust authentication and authorization.

## Features

- Visual flow designer with drag-and-drop interface
- Keycloak-based authentication and group-based permissions
- RESTful API with Bearer token security
- Support for various node types (databases, APIs, triggers, actions)
- Real-time flow execution and monitoring

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

1. `make docker-compose`
2. Import keycloak.json to create a new realm at `http://localhost:8888`.
3. Configure environment variables for database and Keycloak.
4. Create .env and add the variables from env.local.
5. Run `make seed`.
6. Run the Go backend: `make run`
7. Start the Next.js frontend: `make run-web`
8. Access the application at `http://localhost:3000`
9. Access api swagger at `http://localhost:8080/swagger/index.html`
