# Matiq Frontend

A React/Next.js frontend for building and managing automated workflows with a visual drag-and-drop interface.

## Features

- **Visual Matiq**: Drag-and-drop interface for creating workflows
- **Node Types**: Support for triggers, databases, APIs, and action nodes
- **Real-time Connections**: Connect nodes with arrows to define workflow flow
- **Properties Panel**: Configure each node with specific settings
- **API Integration**: Connects to a separate REST backend for persistence

## Setup

1. **Install Dependencies**
   ```bash
   cd web
   npm install
   ```

2. **Environment Configuration**
   ```bash
   cp .env.example .env.local
   ```
   
   Update `.env.local` with your backend API URL:
   ```
   NEXT_PUBLIC_API_URL=http://localhost:8080/api
   ```

3. **Run Development Server**
   ```bash
   npm run dev
   ```

## Backend Integration

This frontend is designed to work with a separate REST backend. See `../API-SPECIFICATION.md` for the complete API specification that your backend should implement.

### Required Backend Endpoints

- `GET /api/workflows` - List all workflows
- `POST /api/workflows` - Create new workflow
- `GET /api/workflows/{id}` - Get specific workflow
- `PUT /api/workflows/{id}` - Update workflow
- `DELETE /api/workflows/{id}` - Delete workflow
- `POST /api/workflows/{id}/execute` - Execute workflow

## Project Structure

```
web/
├── src/
│   ├── app/                 # Next.js app directory
│   ├── components/          # React components
│   │   ├── Canvas.tsx       # Main workflow canvas
│   │   ├── NodeComponent.tsx # Individual node rendering
│   │   ├── Sidebar.tsx      # Node palette
│   │   ├── PropertiesPanel.tsx # Node configuration
│   │   └── WorkflowBuilder.tsx # Main container
│   ├── hooks/               # Custom React hooks
│   │   └── useWorkflow.ts   # Workflow state management
│   ├── services/            # API services
│   │   ├── api-client.ts    # HTTP client
│   │   └── workflow-service.ts # Workflow API calls
│   └── types/               # TypeScript types
│       └── workflow.ts      # Workflow type definitions
└── README.md               # This file
```

## Node Types

### Triggers
- **Schedule**: Time-based triggers (cron expressions)
- **Webhook**: HTTP webhook triggers

### Databases
- PostgreSQL, MySQL, MariaDB, TiDB, Neon
- MongoDB, Snowflake, Supabase
- ClickHouse, Hydra

### APIs
- REST API calls
- GraphQL queries

### Actions
- AI Agent: Custom AI processing
- Transformer: JavaScript data transformation
- Condition: Conditional logic
- Loop: Iterative processing
- Response: HTTP responses
- Error Handler: Error management

## Usage

1. **Create Workflow**: Start by adding trigger nodes from the sidebar
2. **Add Nodes**: Drag additional nodes to build your workflow
3. **Connect Nodes**: Click output points and connect to input points
4. **Configure**: Select nodes to configure their properties
5. **Save**: Changes are automatically saved to the backend

## Development

- **Tech Stack**: React 19, Next.js 15, TypeScript, Tailwind CSS
- **Drag & Drop**: @dnd-kit for smooth node interactions
- **State Management**: Custom hooks with API integration

## Backend Requirements

Your backend should implement the API specification in `../API-SPECIFICATION.md`. The frontend expects:

- RESTful API with JSON responses
- CORS enabled for frontend domain
- Proper error handling with standard HTTP status codes
- Workflow execution engine (optional for basic functionality)
